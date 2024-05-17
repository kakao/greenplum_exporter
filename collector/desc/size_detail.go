package desc

import (
	"github.com/prometheus/client_golang/prometheus"
)

const (
	TblSizeQ = `SELECT sotdschemaname as schema, relname, sotdsize as dsize, sotdtoastsize as toastsize, sotdadditionalsize as othersize
				FROM gp_toolkit.gp_size_of_table_disk as tbl, pg_class
				WHERE tbl.sotdoid=pg_class.oid;`

	IdxSizeQ = `SELECT soiindexschemaname as schema, soitablename as relname, soiindexname as iname, soisize as isize 
				FROM pg_class, gp_toolkit.gp_size_of_index as idx 
				WHERE pg_class.oid = idx.soioid and pg_class.relkind='i';`

	nameDSize = "size_detail"
)

var (
	TblDataSizeDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, nameDSize, "table_data"),
		"Table size (data)", []string{"schema", "table"}, nil,
	)

	TblToastSizeDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, nameDSize, "table_toast"),
		"Table size (data)", []string{"schema", "table"}, nil,
	)
	TblOtherSizeDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, nameDSize, "table_other"),
		"Table size (other)", []string{"schema", "table"}, nil,
	)
	IdxSizeDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, nameDSize, "index"),
		"Index size", []string{"schema", "table", "index"}, nil,
	)
)
