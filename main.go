package main

import (
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/alecthomas/kingpin/v2"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"

	// "flag"
	"greenplum-exporter/collector"
	"net/http"
)

var (
	listenAddress = kingpin.Flag("web.listen-address", "web endpoint").Default("0.0.0.0:9101").String()
	excludeScraps = kingpin.Flag("excludeScraps", "excluding metric names: user,segment,conn,lock,size,size_detail,bgwriter").Default("").String()
	gpVersion     = kingpin.Flag("greenplumVersion", "greenplum Version, greenplum 6.x, greenplum 7.x: 6, 7").Default("6").String()
	loglevel      = kingpin.Flag("loglevel", "Set the level of log: debug, info, warn, error").Default("info").String()
)

var (
	gathers prometheus.Gatherers
	db      *sqlx.DB
)

func establishDBConnection() error {
	// create a database object
	var err error
	pgConnStr := os.Getenv("GPDB_CONNECTION")
	db, err = sqlx.Connect("postgres", pgConnStr)
	if err != nil {
		return err
	}
	if err := db.Ping(); err != nil {
		log.Errorf("Check Connection Information")
		_ = db.Close()
		return err
	}

	db.SetMaxIdleConns(1)
	db.SetMaxOpenConns(1)
	log.Infof("Establish Connection")

	return nil
}

func initHandler(disableDefaultMetrics string, scrapers []collector.Scraper) http.HandlerFunc {
	registry := prometheus.NewRegistry()
	greenplumCollector := collector.New(db, scrapers)

	registry.MustRegister(greenplumCollector)

	gathers = prometheus.Gatherers{registry, prometheus.DefaultGatherer}
	handler := promhttp.HandlerFor(gathers, promhttp.HandlerOpts{
		ErrorHandling: promhttp.ContinueOnError,
	})

	return handler.ServeHTTP
}

func main() {
	kingpin.Version("2.0.0")
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()
	level, err := logrus.ParseLevel(*loglevel)
	if err != nil {
		logrus.WithError(err).Errorf("Invalid log level: %s", *loglevel)
		return
	}
	logrus.SetLevel(level)

	establishDBConnection()

	parsedScraps := collector.NewScrapMetrics(*excludeScraps, db)
	metricsHandleFunc := initHandler(*excludeScraps, parsedScraps)

	mux := http.NewServeMux()
	mux.HandleFunc("/metrics", metricsHandleFunc)

	log.Warnf("exporter's address : %s", *listenAddress)
	log.Error(http.ListenAndServe(*listenAddress, mux).Error())
	defer db.Close()
}
