package collector

import (
	"github.com/prometheus/client_golang/prometheus"
)

// interface of Scraper
type Scraper interface {
	Scrape(ch chan<- prometheus.Metric) //scraper
}
