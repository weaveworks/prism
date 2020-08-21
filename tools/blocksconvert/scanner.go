package blocksconvert

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/prometheus/tsdb/errors"
	"github.com/thanos-io/thanos/pkg/objstore"

	"github.com/cortexproject/cortex/pkg/chunk"
	"github.com/cortexproject/cortex/pkg/storage/tsdb"
	"github.com/cortexproject/cortex/pkg/util/services"
)

type ScannerConfig struct {
	BigtableProject  string
	BigtableInstance string

	TableName    string
	SchemaConfig chunk.SchemaConfig

	OutputDirectory string
	Concurrency     int

	UploadFiles  bool
	Bucket       tsdb.BucketConfig
	BucketPrefix string
}

func (cfg *ScannerConfig) RegisterFlags(f *flag.FlagSet) {
	cfg.SchemaConfig.RegisterFlags(flag.CommandLine)
	cfg.Bucket.RegisterFlags(flag.CommandLine)

	f.StringVar(&cfg.BigtableProject, "bigtable.project", "", "The Google Cloud Platform project ID. Required.")
	f.StringVar(&cfg.BigtableInstance, "bigtable.instance", "", "The Google Cloud Bigtable instance ID. Required.")
	f.StringVar(&cfg.TableName, "table", "", "Table to generate plan files from. If not used, tables are discovered via schema.")
	f.StringVar(&cfg.OutputDirectory, "scanner.local-dir", "", "Local directory used for storing temporary plan files (will be deleted and recreated!).")
	f.IntVar(&cfg.Concurrency, "scanner.concurrency", 16, "Number of concurrent index processors.")
	f.BoolVar(&cfg.UploadFiles, "scanner.upload", true, "Upload plan files.")
	f.StringVar(&cfg.BucketPrefix, "workspace.prefix", "migration", "Prefix in the bucket for storing plan files.")
}

type Scanner struct {
	services.Service

	cfg         ScannerConfig
	indexReader IndexReader

	series    prometheus.Counter
	openFiles prometheus.Gauge
	logger    log.Logger

	tablePeriod time.Duration

	table       string
	tablePrefix string
	bucket      objstore.Bucket
}

func NewScanner(cfg ScannerConfig, l log.Logger, reg prometheus.Registerer) (*Scanner, error) {
	if cfg.BigtableProject == "" || cfg.BigtableInstance == "" {
		return nil, fmt.Errorf("missing BigTable configuration")
	}

	tablePrefix := ""
	tablePeriod := time.Duration(0)
	if cfg.TableName == "" {
		err := cfg.SchemaConfig.Load()
		if err != nil {
			return nil, fmt.Errorf("no table name provided, and schema failed to load: %w", err)
		}

		for _, c := range cfg.SchemaConfig.Configs {
			if c.IndexTables.Period%(24*time.Hour) != 0 {
				return nil, fmt.Errorf("invalid index table period: %v", c.IndexTables.Period)
			}

			if c.Schema != "v9" && c.Schema != "v10" && c.Schema != "v11" {
				return nil, fmt.Errorf("unsupported schema version: %v", c.Schema)
			}

			if tablePrefix == "" {
				tablePrefix = c.IndexTables.Prefix
				tablePeriod = c.IndexTables.Period
			} else if tablePrefix != c.IndexTables.Prefix {
				return nil, fmt.Errorf("multiple index table prefixes found in schema: %v, %v", tablePrefix, c.IndexTables.Prefix)
			} else if tablePeriod != c.IndexTables.Period {
				return nil, fmt.Errorf("multiple index table periods found in schema: %v, %v", tablePeriod, c.IndexTables.Period)
			}
		}
	}

	if cfg.OutputDirectory == "" {
		return nil, fmt.Errorf("no output directory")
	}

	var bucketClient objstore.Bucket
	if cfg.UploadFiles {
		if err := cfg.Bucket.Validate(); err != nil {
			return nil, fmt.Errorf("invalid bucket config: %w", err)
		}

		bucket, err := tsdb.NewBucketClient(context.Background(), cfg.Bucket, "bucket", l, reg)
		if err != nil {
			return nil, fmt.Errorf("failed to create bucket: %w", err)
		}

		bucketClient = bucket
	}

	err := os.MkdirAll(cfg.OutputDirectory, os.FileMode(0700))
	if err != nil {
		return nil, fmt.Errorf("failed to create new output directory %s: %w", cfg.OutputDirectory, err)
	}

	s := &Scanner{
		cfg:         cfg,
		indexReader: NewBigtableIndexReader(cfg.BigtableProject, cfg.BigtableInstance, l, reg),
		table:       cfg.TableName,
		tablePrefix: tablePrefix,
		tablePeriod: tablePeriod,
		logger:      l,
		bucket:      bucketClient,

		series: promauto.With(reg).NewCounter(prometheus.CounterOpts{
			Name: "scanner_series_written_total",
			Help: "Number of series written to the plan files",
		}),

		openFiles: promauto.With(reg).NewGauge(prometheus.GaugeOpts{
			Name: "scanner_open_files",
			Help: "Number of series written to the plan files",
		}),
	}

	s.Service = services.NewBasicService(nil, s.running, nil)
	return s, nil
}

func (s *Scanner) running(ctx context.Context) error {
	tables := []string{}
	if s.table == "" {
		// Use table prefix to discover tables to scan.
		// TODO: use min/max day
		tableNames, err := s.indexReader.IndexTableNames(ctx)
		if err != nil {
			return err
		}

		tables = findTables(s.logger, tableNames, s.tablePrefix, s.tablePeriod)
		level.Info(s.logger).Log("msg", fmt.Sprintf("found %d tables to scan", len(tables)), "prefix", s.tablePrefix, "period", s.tablePeriod)
	} else {
		tables = []string{s.table}
	}

	for _, t := range tables {
		dir := filepath.Join(s.cfg.OutputDirectory, t)
		s.logger.Log("msg", "scanning table", "table", t, "output", dir)

		err := scanSingleTable(ctx, s.indexReader, t, dir, s.cfg.Concurrency, s.openFiles, s.series)
		if err != nil {
			return fmt.Errorf("failed to process table %s: %w", t, err)
		}

		if s.bucket != nil {
			s.logger.Log("msg", "uploading generated plan files", "source", dir)

			err := objstore.UploadDir(ctx, s.logger, s.bucket, dir, s.cfg.BucketPrefix)
			if err != nil {
				return fmt.Errorf("failed to upload %s to bucket: %w", err)
			}
		}
	}

	return nil
}

func findTables(logger log.Logger, tableNames []string, prefix string, period time.Duration) []string {
	type table struct {
		name        string
		periodIndex int64
	}

	var tables []table

	for _, t := range tableNames {
		if !strings.HasPrefix(t, prefix) {
			continue
		}

		if period == 0 {
			tables = append(tables, table{
				name:        t,
				periodIndex: 0,
			})
			continue
		}

		p, err := strconv.ParseInt(t[len(prefix):], 10, 64)
		if err != nil {
			level.Warn(logger).Log("msg", "failed to parse period index of table", "table", t)
			continue
		}

		tables = append(tables, table{
			name:        t,
			periodIndex: p,
		})
	}

	sort.Slice(tables, func(i, j int) bool {
		return tables[i].periodIndex < tables[j].periodIndex
	})

	var out []string
	for _, t := range tables {
		out = append(out, t.name)
	}

	return out
}

func scanSingleTable(ctx context.Context, indexReader IndexReader, tableName string, outDir string, concurrency int, openFiles prometheus.Gauge, series prometheus.Counter) error {
	err := os.RemoveAll(outDir)
	if err != nil {
		return fmt.Errorf("failed to delete directory %s: %w", outDir, err)
	}

	err = os.MkdirAll(outDir, os.FileMode(0700))
	if err != nil {
		return fmt.Errorf("failed to prepare directory %s: %w", outDir, err)
	}

	files := newOpenFiles(128*1024, openFiles)

	var ps []IndexEntryProcessor

	for i := 0; i < concurrency; i++ {
		ps = append(ps, newProcessor(outDir, files, series))
	}

	err = indexReader.ReadIndexEntries(ctx, tableName, ps)
	if err != nil {
		return err
	}

	errs := files.closeAllFiles(func() interface{} {
		return PlanFooter{Complete: true}
	})
	if len(errs) > 0 {
		return errors.MultiError(errs)
	}

	return nil
}
