package collector

import (
	"strings"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/prometheus/client_golang/prometheus"

	"greenplum-exporter/collector/impl"
)

type GreenplumCollector struct {
	db       *sqlx.DB
	scrapers []Scraper
}

func New(db *sqlx.DB, scrapers []Scraper) *GreenplumCollector {
	return &GreenplumCollector{
		db:       db,
		scrapers: scrapers,
	}
}

func (g *GreenplumCollector) Collect(ch chan<- prometheus.Metric) {
	for _, scraper := range g.scrapers {
		scraper.Scrape(ch)
	}
}

func (c *GreenplumCollector) Describe(ch chan<- *prometheus.Desc) {
}

func NewScrapMetrics(excludeScrapsFlag string, db *sqlx.DB) []Scraper {
	var globalScraps = make([]Scraper, 0, 8)

	globalScraps = append(globalScraps, &impl.UpScrap{Gpdb: db})

	// v, _ := strconv.Atoi(version)
	scrapers := map[string]Scraper{
		"user":        &impl.UserScraper{Gpdb: db},
		"segment":     &impl.ClusterScraper{Gpdb: db},
		"conn":        &impl.ConnectionsScraper{Gpdb: db},
		"lock":        &impl.LockScraper{Gpdb: db},
		"size":        &impl.DBSizeScraper{Gpdb: db},
		"size_detail": &impl.DBSizeDetailScraper{Gpdb: db},
		"bgwriter":    &impl.BgWriterScraper{Gpdb: db},
	}

	for _, name := range strings.Split(excludeScrapsFlag, ",") {
		name = strings.TrimSpace(name)
		delete(scrapers, name)
	}

	for _, metric := range scrapers {
		globalScraps = append(globalScraps, metric)
	}

	return globalScraps
}
