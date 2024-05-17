package desc

import (
	"github.com/prometheus/client_golang/prometheus"
)

const (
	LockQ = `SELECT l.pid, d.datname, a.usename, l.locktype, 
				l.mode, a.application_name, a.state, a.query, l.granted,
				LEAST(a.query_start, a.xact_start) AS start,
				COUNT(*) as cnt
			FROM pg_locks as l
				JOIN pg_database d ON l.database = d.oid
				JOIN pg_stat_activity a ON l.pid = a.pid
			WHERE l.pid <> pg_backend_pid() AND a.application_name <> 'pg_statsinfod'
			GROUP BY
				l.pid, d.datname, a.usename, l.locktype,     
				l.mode, a.application_name, a.state,
				l.granted, a.query, start`

	LockStateQ = `SELECT mode, count(*) as cnt FROM pg_locks GROUP BY mode;`
	nameLock   = "lock"
)

var (
	LockDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, nameLock, "detail"),
		"Table locks detail",
		[]string{"pid", "datname", "usename", "locktype", "mode", "application_name", "state", "query", "lock_status", "count"},
		nil,
	)

	LockStateDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, nameLock, "state"),
		"Lock state",
		[]string{"mode"}, nil,
	)
)
