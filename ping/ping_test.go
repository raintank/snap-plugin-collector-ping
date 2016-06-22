package ping

import (
	"testing"

	"github.com/intelsdi-x/snap/control/plugin"
	"github.com/intelsdi-x/snap/control/plugin/cpolicy"
	"github.com/intelsdi-x/snap/core"
	"github.com/intelsdi-x/snap/core/cdata"
	"github.com/intelsdi-x/snap/core/ctypes"
	. "github.com/smartystreets/goconvey/convey"
)

func TestPingPlugin(t *testing.T) {
	Convey("Meta should return metadata for the plugin", t, func() {
		meta := Meta()
		So(meta.Name, ShouldResemble, Name)
		So(meta.Version, ShouldResemble, Version)
		So(meta.Type, ShouldResemble, plugin.CollectorPluginType)
	})

	Convey("Create Ping Collector", t, func() {
		collector := New()
		Convey("So Ping collector should not be nil", func() {
			So(collector, ShouldNotBeNil)
		})
		Convey("So ping collector should be of Ping type", func() {
			So(collector, ShouldHaveSameTypeAs, &Ping{})
		})
		Convey("collector.GetConfigPolicy() should return a config policy", func() {
			configPolicy, _ := collector.GetConfigPolicy()
			Convey("So config policy should not be nil", func() {
				So(configPolicy, ShouldNotBeNil)
				t.Log(configPolicy)
			})
			Convey("So config policy should be a cpolicy.ConfigPolicy", func() {
				So(configPolicy, ShouldHaveSameTypeAs, &cpolicy.ConfigPolicy{})
			})
			Convey("So config policy namespace should be /raintank/ping", func() {
				conf := configPolicy.Get([]string{"raintank", "ping"})
				So(conf, ShouldNotBeNil)
				So(conf.HasRules(), ShouldBeTrue)
				tables := conf.RulesAsTable()
				So(len(tables), ShouldEqual, 3)
				for _, rule := range tables {
					So(rule.Name, ShouldBeIn, "hostname", "timeout", "count")
					switch rule.Name {
					case "hostname":
						So(rule.Required, ShouldBeTrue)
						So(rule.Type, ShouldEqual, "string")
					case "timeout":
						So(rule.Required, ShouldBeFalse)
						So(rule.Type, ShouldEqual, "float")
					case "count":
						So(rule.Required, ShouldBeFalse)
						So(rule.Type, ShouldEqual, "integer")
					}
				}
			})
		})
	})
}

func TestPingollectMetrics(t *testing.T) {
	cfg := setupCfg("127.0.0.1")

	Convey("Ping collector", t, func() {
		p := New()
		mt, err := p.GetMetricTypes(cfg)
		if err != nil {
			t.Fatal(err)
		}
		So(len(mt), ShouldEqual, 6)

		Convey("collect metrics", func() {
			mts := []plugin.MetricType{
				plugin.MetricType{
					Namespace_: core.NewNamespace(
						"raintank", "ping", "avg"),
					Config_: cfg.ConfigDataNode,
				},
				plugin.MetricType{
					Namespace_: core.NewNamespace(
						"raintank", "ping", "min"),
					Config_: cfg.ConfigDataNode,
				},
				plugin.MetricType{
					Namespace_: core.NewNamespace(
						"raintank", "ping", "max"),
					Config_: cfg.ConfigDataNode,
				},
				plugin.MetricType{
					Namespace_: core.NewNamespace(
						"raintank", "ping", "median"),
					Config_: cfg.ConfigDataNode,
				},
				plugin.MetricType{
					Namespace_: core.NewNamespace(
						"raintank", "ping", "mdev"),
					Config_: cfg.ConfigDataNode,
				},
				plugin.MetricType{
					Namespace_: core.NewNamespace(
						"raintank", "ping", "loss"),
					Config_: cfg.ConfigDataNode,
				},
			}
			metrics, err := p.CollectMetrics(mts)
			So(err, ShouldBeNil)
			So(metrics, ShouldNotBeNil)
			So(len(metrics), ShouldEqual, 6)
			So(metrics[0].Namespace()[0].Value, ShouldEqual, "raintank")
			So(metrics[0].Namespace()[1].Value, ShouldEqual, "ping")
			for _, m := range metrics {
				So(m.Namespace()[2].Value, ShouldBeIn, "avg", "min", "max", "median", "mdev", "loss")
				t.Log(m.Namespace()[2].Value, m.Data())
			}
		})
	})
}

func setupCfg(host string) plugin.ConfigType {
	node := cdata.NewNode()
	node.AddItem("hostname", ctypes.ConfigValueStr{Value: host})
	return plugin.ConfigType{ConfigDataNode: node}
}
