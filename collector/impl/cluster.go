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

const (
	primary float64 = 1
	mirror  float64 = 2

	sync      float64 = 1
	resync    float64 = 2
	changeLog float64 = 3
	nosync    float64 = 4

	upStatus   float64 = 1
	downStatus float64 = 0
)

func role2f(role string) float64 {
	if strings.EqualFold(role, "p") {
		return primary
	}
	return mirror
}

func mode2f(mode string) float64 {
	if strings.EqualFold(mode, "s") {
		return sync
	} else if strings.EqualFold(mode, "r") {
		return resync
	} else if strings.EqualFold(mode, "c") {
		return changeLog
	} else {
		return nosync
	}
}

func status2f(status string) float64 {
	if strings.EqualFold(status, "u") {
		return upStatus
	}
	return downStatus
}

type ClusterScraper struct {
	Gpdb *sqlx.DB
}

func (s *ClusterScraper) Scrape(ch chan<- prometheus.Metric) {
	s.scrapeSegConf(ch)
	s.scrapeVersion(ch)
	s.scrapeUpTime(ch)
	s.scrapeSync(ch)
	s.scrapeRate(ch)
}

type scrapeSegConf struct {
	DbID          int64          `db:"dbid"`
	Content       int64          `db:"content"`
	Role          string         `db:"role"`
	PreferredRole string         `db:"preferred_role"`
	Mode          string         `db:"mode"`
	Status        string         `db:"status"`
	Port          int64          `db:"port"`
	Hostname      string         `db:"hostname"`
	Address       string         `db:"address"`
	Datadir       sql.NullString `db:"datadir"`
}

func (s *ClusterScraper) scrapeSegConf(ch chan<- prometheus.Metric) {
	if s.Gpdb == nil {
		return
	}
	db := s.Gpdb
	rows := []scrapeSegConf{}
	if err := db.Select(&rows, desc.SegConfQ); err != nil {
		log.Errorf("fail to scrap: query=%s error=%v", desc.SegConfQ, err)
	}

	totalPSeg := 0
	totalSeg := 0
	upSeg := 0
	for _, row := range rows {
		ch <- prometheus.MustNewConstMetric(desc.SegStatusDesc, prometheus.GaugeValue, status2f(row.Status), strconv.FormatInt(row.DbID, 10), strconv.FormatInt(row.Content, 10), row.PreferredRole, strconv.FormatInt(row.Port, 10), row.Hostname, row.Datadir.String)
		ch <- prometheus.MustNewConstMetric(desc.SegRoleDesc, prometheus.GaugeValue, role2f(row.Role), strconv.FormatInt(row.DbID, 10), strconv.FormatInt(row.Content, 10), row.PreferredRole, strconv.FormatInt(row.Port, 10), row.Hostname, row.Datadir.String)
		ch <- prometheus.MustNewConstMetric(desc.SegModeDesc, prometheus.GaugeValue, mode2f(row.Mode), strconv.FormatInt(row.DbID, 10), strconv.FormatInt(row.Content, 10), row.PreferredRole, strconv.FormatInt(row.Port, 10), row.Hostname, row.Datadir.String)
		if row.Content == -1 && row.Role == "p" {
			ch <- prometheus.MustNewConstMetric(desc.CoordinatorDesc, prometheus.GaugeValue, 1, row.Hostname)
		} else if row.Content == -1 && row.Role == "m" {
			ch <- prometheus.MustNewConstMetric(desc.StandbyDesc, prometheus.GaugeValue, 1, row.Hostname)
		}

		if row.Content != -1 {
			totalSeg++
			if row.Status == "u" {
				upSeg++
			}
			if row.Role == "p" {
				totalPSeg++
			}
		}
	}
	ch <- prometheus.MustNewConstMetric(desc.TotalSegDesc, prometheus.GaugeValue, float64(totalSeg))
	ch <- prometheus.MustNewConstMetric(desc.UpSegDesc, prometheus.GaugeValue, float64(upSeg))
	ch <- prometheus.MustNewConstMetric(desc.TotalPSegDesc, prometheus.GaugeValue, float64(totalPSeg))
}

type scrapeUpTime struct {
	Uptime float64 `db:"date_part"`
}

func (s *ClusterScraper) scrapeUpTime(ch chan<- prometheus.Metric) {
	if s.Gpdb == nil {
		return
	}
	db := s.Gpdb

	rows := []scrapeUpTime{}
	if err := db.Select(&rows, desc.UpTimeQ); err != nil {
		log.Errorf("fail to scrap: query=%s error=%v", desc.UpTimeQ, err)
	}

	for _, row := range rows {
		ch <- prometheus.MustNewConstMetric(desc.UpTimeDesc, prometheus.GaugeValue, row.Uptime)
	}
}

type scrapeVersion struct {
	Version string `db:"version"`
}

func (s *ClusterScraper) scrapeVersion(ch chan<- prometheus.Metric) {
	if s.Gpdb == nil {
		return
	}
	db := s.Gpdb

	rows := []scrapeVersion{}
	if err := db.Select(&rows, desc.VersionQ); err != nil {
		log.Errorf("fail to scrap: query=%s error=%v", desc.VersionQ, err)
	}

	for _, row := range rows {
		lastIndex := strings.LastIndex(row.Version, ".")
		version := row.Version[:lastIndex]

		f, err := strconv.ParseFloat(version, 64)
		if err != nil {
			log.Errorf("Parse error: %v", err)
		}
		ch <- prometheus.MustNewConstMetric(desc.VersionDesc, prometheus.GaugeValue, f)
	}
}

type scrapeSync struct {
	ClientAddr sql.NullString `db:"client_addr"`
	SyncState  string         `db:"sync_state"`
	SyncCount  int64          `db:"count"`
}

func (s *ClusterScraper) scrapeSync(ch chan<- prometheus.Metric) {
	if s.Gpdb == nil {
		return
	}
	db := s.Gpdb

	rows := []scrapeSync{}
	if err := db.Select(&rows, desc.SyncQ); err != nil {
		log.Errorf("fail to scrap: query=%s error=%v", desc.SyncQ, err)
	}

	if len(rows) == 0 {
		ch <- prometheus.MustNewConstMetric(desc.SyncDesc, prometheus.GaugeValue, 0, "", "")
	}

	for _, row := range rows {
		if row.ClientAddr.String == "" {
			continue
		}
		ch <- prometheus.MustNewConstMetric(desc.SyncDesc, prometheus.GaugeValue, float64(row.SyncCount), row.ClientAddr.String, row.SyncState)
	}
}

type scrapeRate struct {
	HitRate  float64 `db:"hit_rate"`
	TxCommit float64 `db:"tx_commit"`
}

func (d *ClusterScraper) scrapeRate(ch chan<- prometheus.Metric) {
	if d.Gpdb == nil {
		return
	}
	db := d.Gpdb

	rows := []scrapeRate{}
	if err := db.Select(&rows, desc.RateQ); err != nil {
		log.Errorf("fail to scrap: query=%s error=%v", desc.RateQ, err)
	}
	for _, row := range rows {
		ch <- prometheus.MustNewConstMetric(desc.HitRateDesc, prometheus.GaugeValue, row.HitRate)
		ch <- prometheus.MustNewConstMetric(desc.TxCommitRateDesc, prometheus.GaugeValue, row.TxCommit)
	}
}
