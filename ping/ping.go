package ping

import (
	"fmt"
	"time"

	"github.com/intelsdi-x/snap/control/plugin"
	"github.com/intelsdi-x/snap/control/plugin/cpolicy"
	"github.com/intelsdi-x/snap/core"
	"github.com/intelsdi-x/snap/core/ctypes"
)

const (
	// Name of plugin
	Name = "ping"
	// Version of plugin
	Version = 1
	// Type of plugin
	Type = plugin.CollectorPluginType
)

var (
	// make sure that we actually satisify requierd interface
	_ plugin.CollectorPlugin = (*Ping)(nil)

	metricNames = []string{
		"avg",
		"min",
		"max",
		"median",
		"mdev",
		"loss",
	}
)

type Ping struct {
}

func New() *Ping {
	return &Ping{}
}

// CollectMetrics collects metrics for testing
func (p *Ping) CollectMetrics(mts []plugin.MetricType) ([]plugin.MetricType, error) {
	var err error

	conf := mts[0].Config().Table()
	hostname, ok := conf["hostname"]
	if !ok || hostname.(ctypes.ConfigValueStr).Value == "" {
		return nil, fmt.Errorf("hostname missing from config, %v", conf)
	}
	var timeout float64
	timeoutConf, ok := conf["timeout"]
	if !ok || timeoutConf.(ctypes.ConfigValueFloat).Value == 0 {
		timeout = 10.0
	} else {
		timeout = timeoutConf.(ctypes.ConfigValueFloat).Value
	}
	var count int
	countConf, ok := conf["count"]
	if !ok || countConf.(ctypes.ConfigValueInt).Value == 0 {
		count = 5
	} else {
		count = countConf.(ctypes.ConfigValueInt).Value
	}

	metrics, err := ping(hostname.(ctypes.ConfigValueStr).Value, count, timeout, mts)
	if err != nil {
		return nil, err
	}

	return metrics, nil
}

func ping(host string, count int, timeout float64, mts []plugin.MetricType) ([]plugin.MetricType, error) {
	check, err := NewRaintankPingProbe(host, count, timeout)
	if err != nil {
		return nil, err
	}
	runTime := time.Now()
	result, err := check.Run()
	if err != nil {
		return nil, err
	}
	stats := make(map[string]float64)
	if result.Avg != nil {
		stats["avg"] = *result.Avg
	}
	if result.Min != nil {
		stats["min"] = *result.Min
	}
	if result.Max != nil {
		stats["max"] = *result.Max
	}
	if result.Median != nil {
		stats["median"] = *result.Median
	}
	if result.Mdev != nil {
		stats["mdev"] = *result.Mdev
	}
	if result.Loss != nil {
		stats["loss"] = *result.Loss
	}

	metrics := make([]plugin.MetricType, 0, len(stats))
	for _, m := range mts {
		stat := m.Namespace()[2].Value
		if value, ok := stats[stat]; ok {
			mt := plugin.MetricType{
				Data_:      value,
				Namespace_: core.NewNamespace("raintank", "ping", stat),
				Timestamp_: runTime,
				Version_:   m.Version(),
			}
			metrics = append(metrics, mt)
		}
	}

	return metrics, nil
}

//GetMetricTypes returns metric types for testing
func (p *Ping) GetMetricTypes(cfg plugin.ConfigType) ([]plugin.MetricType, error) {
	mts := []plugin.MetricType{}
	for _, metricName := range metricNames {
		mts = append(mts, plugin.MetricType{
			Namespace_: core.NewNamespace("raintank", "ping", metricName),
		})
	}
	return mts, nil
}

//GetConfigPolicy returns a ConfigPolicyTree for testing
func (p *Ping) GetConfigPolicy() (*cpolicy.ConfigPolicy, error) {
	c := cpolicy.New()
	rule0, _ := cpolicy.NewStringRule("hostname", true)
	rule1, _ := cpolicy.NewFloatRule("timeout", false, 10.0)
	rule2, _ := cpolicy.NewIntegerRule("count", false, 5)
	cp := cpolicy.NewPolicyNode()
	cp.Add(rule0)
	cp.Add(rule1)
	cp.Add(rule2)
	c.Add([]string{"raintank", "ping"}, cp)
	return c, nil
}

//Meta returns meta data for testing
func Meta() *plugin.PluginMeta {
	return plugin.NewPluginMeta(
		Name,
		Version,
		Type,
		[]string{plugin.SnapGOBContentType},
		[]string{plugin.SnapGOBContentType},
		plugin.Unsecure(true),
		plugin.ConcurrencyCount(5000),
	)
}
