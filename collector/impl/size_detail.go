package impl

import (
	"greenplum-exporter/collector/desc"

	"github.com/jmoiron/sqlx"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
)

type DBSizeDetailScraper struct {
	Gpdb *sqlx.DB
}

func (d *DBSizeDetailScraper) Scrape(ch chan<- prometheus.Metric) {
	d.scrapeTblSize(ch)
	d.scrapeIdxSize(ch)
}

type scrapeTblSize struct {
	Schema    string `db:"schema"`
	TblName   string `db:"relname"`
	DataSize  int64  `db:"dsize"`
	ToastSize int64  `db:"toastsize"`
	OtherSize int64  `db:"othersize"`
}

func (d *DBSizeDetailScraper) scrapeTblSize(ch chan<- prometheus.Metric) {
	if d.Gpdb == nil {
		return
	}
	db := d.Gpdb

	rows := []scrapeTblSize{}
	if err := db.Select(&rows, desc.TblSizeQ); err != nil {
		log.Errorf("fail to scrap: query=%s error=%v", desc.TblSizeQ, err)
	}
	for _, row := range rows {
		ch <- prometheus.MustNewConstMetric(desc.TblDataSizeDesc, prometheus.GaugeValue, float64(row.DataSize), row.Schema, row.TblName)
		ch <- prometheus.MustNewConstMetric(desc.TblToastSizeDesc, prometheus.GaugeValue, float64(row.ToastSize), row.Schema, row.TblName)
		ch <- prometheus.MustNewConstMetric(desc.TblOtherSizeDesc, prometheus.GaugeValue, float64(row.OtherSize), row.Schema, row.TblName)
	}
}

type scrapeIdxSize struct {
	Schema  string `db:"schema"`
	TblName string `db:"relname"`
	IdxName string `db:"iname"`
	IdxSize int64  `db:"isize"`
}

func (d *DBSizeDetailScraper) scrapeIdxSize(ch chan<- prometheus.Metric) {
	if d.Gpdb == nil {
		return
	}
	db := d.Gpdb

	rows := []scrapeIdxSize{}
	if err := db.Select(&rows, desc.IdxSizeQ); err != nil {
		log.Errorf("fail to scrap: query=%s error=%v", desc.IdxSizeQ, err)
	}
	for _, row := range rows {
		ch <- prometheus.MustNewConstMetric(desc.IdxSizeDesc, prometheus.GaugeValue, float64(row.IdxSize), row.Schema, row.TblName, row.IdxName)
	}
}
