// Copyright 2015 CoreOS, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package store

import (
	"github.com/coreos/etcd/Godeps/_workspace/src/github.com/prometheus/client_golang/prometheus"
	"time"
)

// Set of raw Prometheus metrics.
// Labels
// * type = declared in event.go
// * outcome = Outcome
// Do not increment directly, use Report* methods.
var (
	latencyBucketInSeconds = prometheus.ExponentialBuckets(0.001, 2, 13)  // 0.001s to 8.192 sec

	readCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "etcd",
			Subsystem: "store",
			Name:      "reads",
			Help:      "Counter of reads type by (get/getRecursive), outcome (success/failure).",
		}, []string{"type", "outcome"})

	readHandlingTime = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "etcd",
			Subsystem: "store",
			Name:      "read_time_s",
			Help:      "Bucketed histogram of read times (s) by type (get/getRecursive), outcome (success/failure).",
			Buckets:   latencyBucketInSeconds,
		}, []string{"type", "outcome"})

	writeCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "etcd",
			Subsystem: "store",
			Name:      "writes",
			Help:      "Counter of writes by type (set/delete/update/create/compareAndSwap/compareAndDelete/expire) " +
			"outcome(success/failure).",
		}, []string{"type", "outcome"})

	writeHandlingTime = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "etcd",
			Subsystem: "store",
			Name:      "write_time_s",
			Help:      "Bucketed histogram of write times (s) by type " +
			"(set/delete/update/create/compareAndSwap/compareAndDelete/expire) outcome (success/failure).",
			Buckets:   latencyBucketInSeconds,
		}, []string{"type", "outcome"})

	expireCounter = prometheus.NewCounter(
		prometheus.CounterOpts {
			Namespace: "etcd",
			Subsystem: "store",
			Name:      "expires",
			Help:      "Counter of number of key expirations.",
		})

	watchRequests = prometheus.NewCounter(
		prometheus.CounterOpts {
			Namespace: "etcd",
			Subsystem: "store",
			Name:      "watch_requests",
			Help:      "Counter of watch requests incoming into the system.",
		})

	watcherCount = prometheus.NewGauge(
		prometheus.GaugeOpts {
			Namespace: "etcd",
			Subsystem: "store",
			Name:      "watchers",
			Help:      "Number of active watchers.",
		})
)

type Outcome string

const (
	Success Outcome = "success"
	Failure Outcome = "failure"
)

const (
	GetRecursive = "getRecursive"
)

func init() {
	prometheus.MustRegister(readCounter)
	prometheus.MustRegister(writeCounter)
	prometheus.MustRegister(readHandlingTime)
	prometheus.MustRegister(writeHandlingTime)
	prometheus.MustRegister(expireCounter)
	prometheus.MustRegister(watchRequests)
	prometheus.MustRegister(watcherCount)
}

func ReportReadRequest(read_type string, outcome Outcome, start_time time.Time) {
	readCounter.WithLabelValues(read_type, string(outcome)).Inc()
	readHandlingTime.WithLabelValues(read_type, string(outcome)).Observe(time.Since(start_time).Seconds())
}

func ReportWriteRequest(write_type string, outcome Outcome, start_time time.Time) {
	writeCounter.WithLabelValues(write_type, string(outcome)).Inc()
	writeHandlingTime.WithLabelValues(write_type, string(outcome)).Observe(time.Since(start_time).Seconds())
}

func ReportExpiredKey() {
	expireCounter.Inc()
}

func ReportWatchRequest() {
	watchRequests.Inc()
}

func ReportWatcherAdded() {
	watcherCount.Inc()
}

func ReportWatcherRemoved() {
	watcherCount.Dec()
}