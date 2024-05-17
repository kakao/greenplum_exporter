package desc

import (
	"github.com/prometheus/client_golang/prometheus"
)

const (
	UserQ     = `SELECT usename from pg_catalog.pg_user;`
	ConnUserQ = `SELECT usename, 
					COUNT(*) FILTER(WHERE state = 'idle') as user_idle,
					COUNT(*) FILTER(WHERE state = 'active') as user_active,
					COUNT(*) FILTER(WHERE state = 'idle in transaction') as user_wait,
					COUNT(*) FILTER(WHERE state = 'idle in transaction (aborted)') as user_abort
				FROM pg_stat_activity group by usename;`
	nameUser = "user"
)

var (
	UserCountDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, nameUser, "count"),
		"The number of users", nil, nil,
	)

	UserNameDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, nameUser, "name"),
		"User names", []string{"username"}, nil,
	)

	TotalPerUserDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, nameUser, "total_connection_per"),
		"Total Connections Allocated to a Specific Database User", []string{"usename"}, nil,
	)

	IdlePerUserDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, nameUser, "idle_connection_per"),
		"Idle Connections Assigned to a Specific Database User", []string{"usename"}, nil,
	)

	ActivePerUserDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, nameUser, "active_connection_per"),
		"Active Connections Associated with a Specified Database User", []string{"usename"}, nil,
	)

	WaitPerUserDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, nameUser, "wait_connection_per"),
		"Waited Connections Associated with a Specified Database User", []string{"usename"}, nil,
	)
	AbortPerUserDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, nameUser, "abort_connection_per"),
		"Aborted Connections Associated with a Specified Database User", []string{"usename"}, nil,
	)

	TotalCountOnlineUsersDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, nameUser, "total_online_count"),
		"Total Online User Count", nil, nil,
	)
)
