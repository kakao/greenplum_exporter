package impl

import (
	"greenplum-exporter/collector/desc"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
)

type BgWriterScraper struct {
	Gpdb *sqlx.DB
}

type bgWriter struct {
	CheckpointsTimed    int64     `db:"checkpoints_timed"`
	CheckpointsReq      int64     `db:"checkpoints_req"`
	CheckpointWriteTime float64   `db:"checkpoint_write_time"`
	CheckpointSyncTime  float64   `db:"checkpoint_sync_time"`
	BuffersCheckpoint   int64     `db:"buffers_checkpoint"`
	BuffersClean        int64     `db:"buffers_clean"`
	MaxwrittenClean     int64     `db:"maxwritten_clean"`
	BuffersBackend      int64     `db:"buffers_backend"`
	BuffersBackendFsync int64     `db:"buffers_backend_fsync"`
	BuffersAlloc        int64     `db:"buffers_alloc"`
	StatsReset          time.Time `db:"stats_reset"`
}

func (b *BgWriterScraper) Scrape(ch chan<- prometheus.Metric) {
	if b.Gpdb == nil {
		return
	}
	db := b.Gpdb
	rows := []bgWriter{}
	if err := db.Select(&rows, desc.BgwriterQ); err != nil {
		log.Errorf("fail to scrap: query=%s error=%v", desc.BgwriterQ, err)
	}

	for _, row := range rows {
		ch <- prometheus.MustNewConstMetric(desc.CheckpointsTimedDesc, prometheus.CounterValue, float64(row.CheckpointsTimed))
		ch <- prometheus.MustNewConstMetric(desc.CheckpointsReqDesc, prometheus.CounterValue, float64(row.CheckpointsReq))
		ch <- prometheus.MustNewConstMetric(desc.CheckpointWriteTimeDesc, prometheus.CounterValue, row.CheckpointWriteTime/1000)
		ch <- prometheus.MustNewConstMetric(desc.CheckpointSyncTimeDesc, prometheus.CounterValue, row.CheckpointSyncTime/1000)
		ch <- prometheus.MustNewConstMetric(desc.BuffersCheckpointDesc, prometheus.CounterValue, float64(row.BuffersCheckpoint))
		ch <- prometheus.MustNewConstMetric(desc.BuffersCleanDesc, prometheus.CounterValue, float64(row.BuffersClean))
		ch <- prometheus.MustNewConstMetric(desc.MaxwrittenCleanDesc, prometheus.CounterValue, float64(row.MaxwrittenClean))
		ch <- prometheus.MustNewConstMetric(desc.BuffersBackendDesc, prometheus.CounterValue, float64(row.BuffersBackend))
		ch <- prometheus.MustNewConstMetric(desc.BuffersBackendFsyncDesc, prometheus.CounterValue, float64(row.BuffersBackendFsync))
		ch <- prometheus.MustNewConstMetric(desc.BuffersAllocDesc, prometheus.CounterValue, float64(row.BuffersAlloc))
		ch <- prometheus.MustNewConstMetric(desc.StatsResetDesc, prometheus.GaugeValue, float64(row.StatsReset.UTC().Unix()))
	}
}
