package impl

import (
	"database/sql"
	"greenplum-exporter/collector/desc"
	"strconv"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
)

type LockScraper struct {
	Gpdb *sqlx.DB
}

func (l *LockScraper) Scrape(ch chan<- prometheus.Metric) {
	l.scrapeLock(ch)
	l.scrapeLockState(ch)
}

type scrapeLock struct {
	Pid             string       `db:"pid"`
	Datname         string       `db:"datname"`
	Usename         string       `db:"usename"`
	Locktype        string       `db:"locktype"`
	Mode            string       `db:"mode"`
	ApplicationName string       `db:"application_name"`
	Query           string       `db:"query"`
	State           string       `db:"state"`
	Granted         string       `db:"granted"`
	Start           sql.NullTime `db:"start"`
	Count           int64        `db:"cnt"`
}

func granted2status(granted string) string {
	if strings.EqualFold(granted, "t") {
		return "wait_lock"
	}
	return "get_lock"
}

func (l *LockScraper) scrapeLock(ch chan<- prometheus.Metric) {
	if l.Gpdb == nil {
		return
	}
	db := l.Gpdb

	rows := []scrapeLock{}
	if err := db.Select(&rows, desc.LockQ); err != nil {
		log.Errorf("fail to scrap: query=%s error=%v", desc.LockQ, err)
	}
	for _, row := range rows {
		ch <- prometheus.MustNewConstMetric(desc.LockDesc, prometheus.GaugeValue, float64(row.Start.Time.UTC().Unix()),
			row.Pid, row.Datname, row.Usename, row.Locktype, row.Mode, row.ApplicationName, row.State, granted2status(row.Granted),
			row.Query, strconv.FormatInt(row.Count, 10))
	}
}

type scrapeLockState struct {
	Mode  string `db:"mode"`
	Count int64  `db:"cnt"`
}

func (l *LockScraper) scrapeLockState(ch chan<- prometheus.Metric) {
	if l.Gpdb == nil {
		return
	}
	db := l.Gpdb

	rows := []scrapeLockState{}
	if err := db.Select(&rows, desc.LockStateQ); err != nil {
		log.Errorf("fail to scrap: query=%s error=%v", desc.LockStateQ, err)
	}
	for _, row := range rows {
		ch <- prometheus.MustNewConstMetric(desc.LockStateDesc, prometheus.GaugeValue, float64(row.Count), row.Mode)
	}
}
