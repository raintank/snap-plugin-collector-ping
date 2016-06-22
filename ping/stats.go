package ping

import (
	"fmt"
	"math"
	"net"
	"sort"
	"time"

	"github.com/raintank/go-pinger"
)

// global co-oridinator shared between all go-routines.
var GlobalPinger *pinger.Pinger

func init() {
	GlobalPinger = pinger.NewPinger()
}

// results. we use pointers so that missing data will be
// encoded as 'null' in the json response.
type PingResult struct {
	Loss   *float64
	Min    *float64
	Max    *float64
	Avg    *float64
	Median *float64
	Mdev   *float64
}

// Our check definition.
type RaintankProbePing struct {
	Hostname string
	Count    int
	Timeout  time.Duration
}

// parse the json request body to build our check definition.
func NewRaintankPingProbe(hostname string, count int, timeout float64) (*RaintankProbePing, error) {
	p := RaintankProbePing{
		Hostname: hostname,
		Count:    count,
	}
	if p.Hostname == "" {
		return nil, fmt.Errorf("no host passed.")
	}
	if p.Count <= 0 {
		return nil, fmt.Errorf("invalid count, but be greater then 0.")
	}

	if timeout <= 0.0 {
		return nil, fmt.Errorf("invalid value for timeout, must be greater then 0.")
	}
	p.Timeout = time.Duration(time.Millisecond * time.Duration(int(1000.0*timeout)))

	return &p, nil
}

func (p *RaintankProbePing) Run() (*PingResult, error) {
	deadline := time.Now().Add(p.Timeout)
	result := &PingResult{}

	var ipAddr string

	// get IP from hostname.
	addrs, err := net.LookupHost(p.Hostname)
	if err != nil || len(addrs) < 1 {
		loss := 100.0
		result.Loss = &loss
		return result, nil
	}
	if time.Now().After(deadline) {
		//timeout resolving IP address of hostname
		loss := 100.0
		result.Loss = &loss
		return result, nil
	}
	ipAddr = addrs[0]

	resultsChan, err := GlobalPinger.Ping(ipAddr, p.Count, deadline)
	if err != nil {
		return nil, err
	}

	results := <-resultsChan

	// derive stats from results.
	successCount := results.Received
	failCount := results.Sent - results.Received

	measurements := make([]float64, len(results.Latency))
	for i, m := range results.Latency {
		measurements[i] = m.Seconds() * 1000
	}

	tsum := 0.0
	tsum2 := 0.0
	min := math.Inf(1)
	max := 0.0
	for _, r := range measurements {
		if r > max {
			max = r
		}
		if r < min {
			min = r
		}
		tsum += r
		tsum2 += (r * r)
	}

	if successCount > 0 {
		avg := tsum / float64(successCount)
		result.Avg = &avg
		root := math.Sqrt((tsum2 / float64(successCount)) - ((tsum / float64(successCount)) * (tsum / float64(successCount))))
		result.Mdev = &root
		sort.Float64s(measurements)
		median := measurements[successCount/2]
		result.Median = &median
		result.Min = &min
		result.Max = &max
	}
	if failCount == 0 {
		loss := 0.0
		result.Loss = &loss
	} else {
		loss := 100.0 * (float64(failCount) / float64(results.Sent))
		result.Loss = &loss
	}

	return result, nil
}
