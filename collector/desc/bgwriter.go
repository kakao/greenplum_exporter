package desc

import (
	"github.com/prometheus/client_golang/prometheus"
)

const (
	BgwriterQ    = `SELECT * FROM pg_stat_bgwriter`
	nameBgwriter = "bgwriter"
)

var (
	CheckpointsTimedDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, nameBgwriter, "checkpoints_timed"),
		"The total amount of planned checkpoints that have been completed", nil, nil,
	)

	CheckpointsReqDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, nameBgwriter, "checkpoints_req"),
		"The count of requested checkpoints that have been carried out", nil, nil,
	)

	CheckpointWriteTimeDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, nameBgwriter, "checkpoint_write_time"),
		"The total duration spent while writing files to the disk during a checkpoint", nil, nil,
	)

	CheckpointSyncTimeDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, nameBgwriter, "checkpoint_sync_time"),
		"The total duration of time dedicated to synchronizing files to the disk during a checkpoint", nil, nil,
	)

	BuffersCheckpointDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, nameBgwriter, "buffers_checkpoint"),
		"The amount of buffered data during a checkpoint", nil, nil,
	)

	BuffersCleanDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, nameBgwriter, "buffers_clean"),
		"The amount of buffered data which have been cleaned during a checkpoint", nil, nil,
	)

	MaxwrittenCleanDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, nameBgwriter, "maxwritten_clean"),
		"The number of data where the background writer halted a cleaning scan due to excessive buffer writing", nil, nil,
	)

	BuffersBackendDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, nameBgwriter, "buffers_backend"),
		"The amount of buffers directly written by a backend", nil, nil,
	)

	BuffersBackendFsyncDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, nameBgwriter, "buffers_backend_fsync"),
		"The frequency of fsync call", nil, nil,
	)

	BuffersAllocDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, nameBgwriter, "buffers_alloc"),
		"The amount of memory buffers that have been assigned or allocated", nil, nil,
	)

	StatsResetDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, nameBgwriter, "stats_reset"),
		"The most recent point in time when these statistics were reset", nil, nil,
	)
)
