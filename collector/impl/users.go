package impl

import (
	"database/sql"
	"greenplum-exporter/collector/desc"

	"github.com/jmoiron/sqlx"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
)

type UserScraper struct {
	Gpdb *sqlx.DB
}

type scrapeUsername struct {
	Usename string `db:"usename"`
}

func (u *UserScraper) Scrape(ch chan<- prometheus.Metric) {
	u.scrapeUsername(ch)
	u.scrapeUserConn(ch)
}

func (u *UserScraper) scrapeUsername(ch chan<- prometheus.Metric) {
	if u.Gpdb == nil {
		return
	}
	db := u.Gpdb

	rows := []scrapeUsername{}
	if err := db.Select(&rows, desc.UserQ); err != nil {
		log.Errorf("fail to scrap: query=%s error=%v", desc.UserQ, err)
	}
	for _, row := range rows {
		ch <- prometheus.MustNewConstMetric(desc.UserNameDesc, prometheus.GaugeValue, 1, row.Usename)
	}
	ch <- prometheus.MustNewConstMetric(desc.UserCountDesc, prometheus.GaugeValue, float64(len(rows)))
}

type scrapeUserConn struct {
	Usename sql.NullString `db:"usename"`
	Idle    float64        `db:"user_idle"`
	Active  float64        `db:"user_active"`
	Wait    float64        `db:"user_wait"`
	Abort   float64        `db:"user_abort"`
}

func (u *UserScraper) scrapeUserConn(ch chan<- prometheus.Metric) {
	if u.Gpdb == nil {
		return
	}
	db := u.Gpdb

	rows := []scrapeUserConn{}
	if err := db.Select(&rows, desc.ConnUserQ); err != nil {
		log.Errorf("fail to scrap: query=%s error=%v", desc.ConnUserQ, err)
	}

	totalOnlineUserCount := 0
	for _, row := range rows {
		if row.Usename.String == "" {
			continue
		}
		ch <- prometheus.MustNewConstMetric(desc.IdlePerUserDesc, prometheus.GaugeValue, row.Idle, row.Usename.String)
		ch <- prometheus.MustNewConstMetric(desc.ActivePerUserDesc, prometheus.GaugeValue, row.Active, row.Usename.String)
		ch <- prometheus.MustNewConstMetric(desc.WaitPerUserDesc, prometheus.GaugeValue, row.Wait, row.Usename.String)
		ch <- prometheus.MustNewConstMetric(desc.AbortPerUserDesc, prometheus.GaugeValue, row.Abort, row.Usename.String)
		ch <- prometheus.MustNewConstMetric(desc.TotalPerUserDesc, prometheus.GaugeValue, row.Idle+row.Active+row.Wait+row.Abort, row.Usename.String)

		totalOnlineUserCount++
	}

	ch <- prometheus.MustNewConstMetric(desc.TotalCountOnlineUsersDesc, prometheus.GaugeValue, float64(totalOnlineUserCount))
}
