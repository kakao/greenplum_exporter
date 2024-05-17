package desc

import (
	"github.com/prometheus/client_golang/prometheus"
)

const (
	SegConfQ = `select * from gp_segment_configuration;`
	VersionQ = `SELECT (regexp_matches(version(), 'Greenplum Database (\d+\.\d+\.\d+)'))[1] as version;`
	UpTimeQ  = `select extract(epoch from now() - pg_postmaster_start_time())`
	SyncQ    = `SELECT client_addr, sync_state, count(*) 
				FROM pg_stat_replication where state='streaming' 
				GROUP By client_addr, sync_state`
	RateQ = `SELECT 
				sum(blks_hit)/(sum(blks_read)+sum(blks_hit))*100 as hit_rate,
				sum(xact_commit)/(sum(xact_commit)+sum(xact_rollback))*100 as tx_commit 
			FROM pg_stat_database;
	`
	nameCluster = "cluster"
)

var (
	TotalPSegDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, nameCluster, "total_segments"),
		"The number of primary segments ", nil, nil,
	)

	VersionDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, nameCluster, "version"),
		"Greenplum version", nil, nil,
	)

	CoordinatorDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, nameCluster, "coordinator"),
		"Coordinator name", []string{"coordinator"}, nil,
	)

	StandbyDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, nameCluster, "standby"),
		"Standby name", []string{"standby"}, nil,
	)

	UpTimeDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, nameCluster, "uptime"),
		"The duration in seconds since GP was last started", nil, nil,
	)

	SyncDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, nameCluster, "coordinator_sync_with_standby"),
		"The status of whether the coordinator is currently synchronizing to the standby", []string{"client", "state"}, nil,
	)

	TotalSegDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, nameCluster, "the_number_of_segments"),
		"The number of segments within a cluster", nil, nil,
	)

	UpSegDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, nameCluster, "segment_up"),
		"The number of active segments", nil, nil,
	)
	SegStatusDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, nameCluster, "segment_status"),
		"UP(1) indicates that the segment is on, otherwise off",
		[]string{"dbid", "content", "preferred_role", "port", "hostname", "datadir"}, nil,
	)

	SegRoleDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, nameCluster, "segment_role"),
		"The current role of a segment, which can be either primary or mirror",
		[]string{"dbid", "content", "preferred_role", "port", "hostname", "datadir"}, nil,
	)

	SegModeDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, nameCluster, "segment_mode"),
		"The status of replication for each segment",
		[]string{"dbid", "content", "preferred_role", "port", "hostname", "datadir"}, nil,
	)
	HitRateDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, nameCluster, "hit_rate"),
		"Cache hit rate for databases", nil, nil,
	)

	TxCommitRateDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, nameCluster, "transition_commit_rate"),
		"Transaction commit rate for databases", nil, nil,
	)
)
