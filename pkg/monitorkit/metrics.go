package monitorkit

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var NoahHeartbeatRequestCount = promauto.NewCounterVec(
	prometheus.CounterOpts{
		Name: "noah_heartbeat_requests_total",
		Help: "Total number of heartbeat requests by skey.",
	},
	[]string{"mkey"},
)

var NoahUserLoginRequestCount = promauto.NewCounterVec(
	prometheus.CounterOpts{
		Name: "user_login_requests_total",
		Help: "Total number of user login requests by skey.",
	},
	[]string{"mkey"},
)
