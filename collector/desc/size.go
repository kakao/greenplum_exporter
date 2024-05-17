package desc

import (
	"github.com/prometheus/client_golang/prometheus"
)

const (
	DbSizeQ          = `SELECT sodddatname, sodddatsize as db_size from gp_toolkit.gp_size_of_database;`
	SegDiskFreeSizeQ = `SELECT dfhostname,sum(dfspace)/count(dfspace) as dfspace_free from gp_toolkit.gp_disk_free GROUP BY dfhostname;`

	nameSize = "size"
)

var (
	DbSizeDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, nameSize, "db"),
		"Total File System Size (in MB) of Each Database", []string{"dbname"}, nil,
	)

	SegDiskFreeSizeDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, nameSize, "segment_disk_free"),
		"The total amount of free space available on the disk in the file system for each segment node",
		[]string{"hostname"}, nil,
	)
)
