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

package proxy

import (
	"github.com/coreos/etcd/Godeps/_workspace/src/github.com/prometheus/client_golang/prometheus"
	"net/http"
	"time"
)

var (
	requestsIncoming = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "etcd",
			Subsystem: "proxy",
			Name:      "requests",
			Help:      "Counter requests incoming by method.",
		}, []string{"method"})

	requestsHandled = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "etcd",
			Subsystem: "proxy",
			Name:      "handled",
			Help:      "Counter of requests fully handled (by authoratitave servers)",
		}, []string{"method", "code"})

	requestsDropped = prometheus.NewCounterVec(
		prometheus.CounterOpts {
			Namespace: "etcd",
			Subsystem: "proxy",
			Name:      "dropped",
			Help:      "Counter of requests dropped on the proxy.",
		},[]string{"method", "proxying_error"})

	requestsHandlingTime = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "etcd",
			Subsystem: "http",
			Name:      "handling_time_s",
			Help:      "Bucketed histogram of handling time of successful events (non-watches), by method " +
			"(GET/PUT etc.).",
			Buckets:   prometheus.ExponentialBuckets(0.0005, 2, 13),
		}, []string{"method"})
)

type ProxyingError string

const (
	ZeroEndpoints ProxyingError = "zero_endpoints"
	FailedSendingRequest ProxyingError = "failed_sending_request"
	FailedGettingResponse ProxyingError = "failed_getting_response"
)

func init() {
	prometheus.MustRegister(requestsIncoming)
	prometheus.MustRegister(requestsHandled)
	prometheus.MustRegister(requestsDropped)
	prometheus.MustRegister(requestsHandlingTime)
}

func ReportIncomingRequest(request *http.Request) {
	requestsIncoming.WithLabelValues(request.Method).Inc()
}

func ReportRequestHandled(request *http.Request, response *http.Response, startTime time.Time) {
	method := request.Method
	requestsHandled.WithLabelValues(method).Inc()
	requestsHandlingTime.WithLabelValues(method).Observe(time.Since(startTime).Seconds())
}

func ReportRequestDropped(request *http.Request, err ProxyingError) {
	requestsDropped.WithLabelValues(request.Method, string(err))
}