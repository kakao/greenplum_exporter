package impl

import (
	"greenplum-exporter/collector/desc"

	"github.com/jmoiron/sqlx"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
)

type UpScrap struct {
	Gpdb *sqlx.DB
}

type up struct {
	Ok string `db:"ok"`
}

func (u *UpScrap) Scrape(ch chan<- prometheus.Metric) {
	if u.Gpdb == nil {
		return
	}
	db := u.Gpdb

	rows := []up{}
	if err := db.Select(&rows, desc.UpQ); err != nil {
		log.Errorf("fail to scrap: query=%s error=%v", desc.UpQ, err)
	}
	ch <- prometheus.MustNewConstMetric(desc.UpDesc, prometheus.GaugeValue, float64(len(rows)))
}
