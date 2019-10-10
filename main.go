package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/push"
)

var (
	completionTime = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: "test",
		Name:      "last_completion_timestamp_seconds",
		Help:      "The timestamp of the last completion, successful or not.",
	})
	successTime = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: "test",
		Name:      "last_success_timestamp_seconds",
		Help:      "The timestamp of the last successful completion.",
	})
	duration = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: "test",
		Name:      "duration_seconds",
		Help:      "The duration of the last in seconds.",
	})
)

func main() {
	rand.Seed(time.Now().Unix())
	pusher := push.New("http://pushgateway.kube-pushgateway-test.svc.cluster.local.:9091", "test-job")

	startTime := time.Now()

	// Do something
	time.Sleep(time.Second * 2)

	completionTime.SetToCurrentTime()
	if rand.Intn(2) == 0 {
		successTime.SetToCurrentTime()
		pusher.Collector(successTime)
	}

	duration.Set(time.Since(startTime).Seconds())
	if err := pusher.
		Collector(completionTime).
		Collector(duration).
		Add(); err != nil {
		fmt.Println("Could not push completion time to Pushgateway:", err)
	}
}
