package desc

import (
	"github.com/prometheus/client_golang/prometheus"
)

const UpQ = `select 'ALIVE' as ok`

var (
	UpDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "up"),
		"Greenplum Cluster Accessibility Status", nil, nil,
	)
)
