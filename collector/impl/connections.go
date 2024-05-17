package impl

import (
	"database/sql"
	"greenplum-exporter/collector/desc"

	"github.com/jmoiron/sqlx"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
)

type ConnectionsScraper struct {
	Gpdb *sqlx.DB
}

func (c *ConnectionsScraper) Scrape(ch chan<- prometheus.Metric) {
	c.scrapeClient(ch)
	c.scrapeMaxConn(ch)
}

type scrapeClient struct {
	ClientAddr   sql.NullString `db:"client_addr"`
	ClientIdle   float64        `db:"client_idle"`
	ClientActive float64        `db:"client_active"`
	ClientWait   float64        `db:"client_wait"`
	ClientAbort  float64        `db:"client_abort"`
}

func (c *ConnectionsScraper) scrapeClient(ch chan<- prometheus.Metric) {
	if c.Gpdb == nil {
		return
	}
	db := c.Gpdb
	rows := []scrapeClient{}
	if err := db.Select(&rows, desc.ConnStateQ); err != nil {
		log.Errorf("fail to scrap: query=%s error=%v", desc.ConnStateQ, err)
	}

	var clientActive float64 = 0

	var idle float64 = 0   // sum of client idle
	var active float64 = 0 // sum of client active
	var wait float64 = 0   // sum of client wait
	var abort float64 = 0  // sum of client abort

	totalClientCount := 0
	for _, row := range rows {
		if row.ClientAddr.String == "" {
			continue
		}
		ch <- prometheus.MustNewConstMetric(desc.IdlePerClientDesc, prometheus.GaugeValue, row.ClientIdle, row.ClientAddr.String)
		ch <- prometheus.MustNewConstMetric(desc.ActivePerClientDesc, prometheus.GaugeValue, clientActive, row.ClientAddr.String)
		ch <- prometheus.MustNewConstMetric(desc.WaitPerClientDesc, prometheus.GaugeValue, clientActive, row.ClientAddr.String)
		ch <- prometheus.MustNewConstMetric(desc.AbortPerClientDesc, prometheus.GaugeValue, clientActive, row.ClientAddr.String)
		ch <- prometheus.MustNewConstMetric(desc.TotalPerClientDesc, prometheus.GaugeValue, row.ClientIdle+clientActive, row.ClientAddr.String)

		idle += row.ClientIdle
		active += row.ClientActive
		wait += row.ClientWait
		abort += row.ClientAbort
		//totalClientCount is the sum of client addresses (not null)
		totalClientCount++
	}

	ch <- prometheus.MustNewConstMetric(desc.TotalCountClientDesc, prometheus.GaugeValue, float64(totalClientCount))

	ch <- prometheus.MustNewConstMetric(desc.IdleConnDesc, prometheus.GaugeValue, idle)
	ch <- prometheus.MustNewConstMetric(desc.ActiveConnDesc, prometheus.GaugeValue, active)
	ch <- prometheus.MustNewConstMetric(desc.WaitConnDesc, prometheus.GaugeValue, wait)
	ch <- prometheus.MustNewConstMetric(desc.AbortConnDesc, prometheus.GaugeValue, abort)
	ch <- prometheus.MustNewConstMetric(desc.TotalConnDesc, prometheus.GaugeValue, idle+active+wait+abort)
}

type scrapeMaxConn struct {
	Name    string  `db:"name"`
	MaxConn float64 `db:"setting"`
}

func (c *ConnectionsScraper) scrapeMaxConn(ch chan<- prometheus.Metric) {
	if c.Gpdb == nil {
		return
	}
	db := c.Gpdb
	rows := []scrapeMaxConn{}
	if err := db.Select(&rows, desc.MaxConnQ); err != nil {
		log.Errorf("fail to scrap: query=%s error=%v", desc.MaxConnQ, err)
	}

	var max_connections float64
	var superuser_reserved_connections float64
	for _, row := range rows {
		if row.Name == "max_connections" {
			max_connections = float64(row.MaxConn)
		} else if row.Name == "superuser_reserved_connections" {
			superuser_reserved_connections = float64(row.MaxConn)
		}
	}
	ch <- prometheus.MustNewConstMetric(desc.MaxConnDesc, prometheus.GaugeValue, max_connections-superuser_reserved_connections)
}
