package impl

import (
	"greenplum-exporter/collector/desc"

	"github.com/jmoiron/sqlx"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
)

type DBSizeScraper struct {
	Gpdb *sqlx.DB
}

func (d *DBSizeScraper) Scrape(ch chan<- prometheus.Metric) {
	d.scrapeDbSize(ch)
	d.scrapeSegDiskFree(ch)
}

type scrapeDbSize struct {
	DbName string  `db:"sodddatname"`
	DbSize float64 `db:"db_size"`
}

func (d *DBSizeScraper) scrapeDbSize(ch chan<- prometheus.Metric) {
	if d.Gpdb == nil {
		return
	}
	db := d.Gpdb

	rows := []scrapeDbSize{}
	if err := db.Select(&rows, desc.DbSizeQ); err != nil {
		log.Errorf("fail to scrap: query=%s error=%v", desc.DbSizeQ, err)
	}
	for _, row := range rows {
		ch <- prometheus.MustNewConstMetric(desc.DbSizeDesc, prometheus.GaugeValue, row.DbSize, row.DbName)
	}
}

type scrapeSegDiskFree struct {
	HostName string  `db:"dfhostname"`
	MbSize   float64 `db:"dfspace_free"`
}

func (d *DBSizeScraper) scrapeSegDiskFree(ch chan<- prometheus.Metric) {
	if d.Gpdb == nil {
		return
	}
	db := d.Gpdb

	rows := []scrapeSegDiskFree{}
	if err := db.Select(&rows, desc.SegDiskFreeSizeQ); err != nil {
		log.Errorf("fail to scrap: query=%s error=%v", desc.SegDiskFreeSizeQ, err)
	}
	for _, row := range rows {
		ch <- prometheus.MustNewConstMetric(desc.SegDiskFreeSizeDesc, prometheus.GaugeValue, row.MbSize, row.HostName)
	}
}
