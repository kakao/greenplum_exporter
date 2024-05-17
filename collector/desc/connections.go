package desc

import (
	"github.com/prometheus/client_golang/prometheus"
)

const (
	ConnStateQ = `SELECT client_addr,
				COUNT(*) FILTER(WHERE state = 'idle') as client_idle,
				COUNT(*) FILTER(WHERE state = 'active') as client_active,
				COUNT(*) FILTER(WHERE state = 'idle in transaction') as client_wait,
				COUNT(*) FILTER(WHERE state = 'idle in transaction (aborted)') as client_abort
			FROM pg_stat_activity
			WHERE pid <> pg_backend_pid() 
			GROUP BY client_addr;`

	MaxConnQ = `SELECT name, setting 
				FROM pg_settings 
				WHERE name In('max_connections', 'superuser_reserved_connections');`

	nameConnection = "connection"
)

var (
	TotalConnDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, nameConnection, "total"),
		"The number of current connections", nil, nil,
	)

	IdleConnDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, nameConnection, "idle"),
		"The number of idle connections", nil, nil,
	)

	ActiveConnDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, nameConnection, "active"),
		"The number of active connections", nil, nil,
	)

	WaitConnDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, nameConnection, "wait"),
		"The number of waited queries", nil, nil,
	)

	AbortConnDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, nameConnection, "abort"),
		"The number of aborted queries", nil, nil,
	)

	MaxConnDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, nameConnection, "max"),
		"The maximum number of connections", nil, nil,
	)
	TotalPerClientDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, nameConnection, "total_per_client"),
		"The number of connections for each client", []string{"client"}, nil,
	)
	IdlePerClientDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, nameConnection, "idle_per_client"),
		"The number of idle connections for each client", []string{"client"}, nil,
	)
	ActivePerClientDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, nameConnection, "active_per_client"),
		"The number of active connections for each client", []string{"client"}, nil,
	)
	WaitPerClientDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, nameConnection, "wait_per_client"),
		"The number of waited connections for each client", []string{"client"}, nil,
	)
	AbortPerClientDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, nameConnection, "abort_per_client"),
		"The number of aborted connections for each client", []string{"client"}, nil,
	)
	TotalCountClientDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, nameConnection, "total_client_count"),
		"The number of client connections", nil, nil,
	)
)
