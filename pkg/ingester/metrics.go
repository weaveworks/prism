package ingester

import (
	"sync"

	"github.com/cortexproject/cortex/pkg/util"
	"github.com/go-kit/kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	memSeriesCreatedTotalName = "cortex_ingester_memory_series_created_total"
	memSeriesCreatedTotalHelp = "The total number of series that were created per user."

	memSeriesRemovedTotalName = "cortex_ingester_memory_series_removed_total"
	memSeriesRemovedTotalHelp = "The total number of series that were removed per user."
)

type ingesterMetrics struct {
	flushQueueLength      prometheus.Gauge
	ingestedSamples       prometheus.Counter
	ingestedSamplesFail   prometheus.Counter
	queries               prometheus.Counter
	queriedSamples        prometheus.Histogram
	queriedSeries         prometheus.Histogram
	queriedChunks         prometheus.Histogram
	memSeries             prometheus.Gauge
	memUsers              prometheus.Gauge
	memSeriesCreatedTotal *prometheus.CounterVec
	memSeriesRemovedTotal *prometheus.CounterVec
}

func newIngesterMetrics(r prometheus.Registerer, registerMetricsConflictingWithTSDB bool) *ingesterMetrics {
	m := &ingesterMetrics{
		flushQueueLength: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "cortex_ingester_flush_queue_length",
			Help: "The total number of series pending in the flush queue.",
		}),
		ingestedSamples: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "cortex_ingester_ingested_samples_total",
			Help: "The total number of samples ingested.",
		}),
		ingestedSamplesFail: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "cortex_ingester_ingested_samples_failures_total",
			Help: "The total number of samples that errored on ingestion.",
		}),
		queries: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "cortex_ingester_queries_total",
			Help: "The total number of queries the ingester has handled.",
		}),
		queriedSamples: prometheus.NewHistogram(prometheus.HistogramOpts{
			Name: "cortex_ingester_queried_samples",
			Help: "The total number of samples returned from queries.",
			// Could easily return 10m samples per query - 10*(8^(8-1)) = 20.9m.
			Buckets: prometheus.ExponentialBuckets(10, 8, 8),
		}),
		queriedSeries: prometheus.NewHistogram(prometheus.HistogramOpts{
			Name: "cortex_ingester_queried_series",
			Help: "The total number of series returned from queries.",
			// A reasonable upper bound is around 100k - 10*(8^(6-1)) = 327k.
			Buckets: prometheus.ExponentialBuckets(10, 8, 6),
		}),
		queriedChunks: prometheus.NewHistogram(prometheus.HistogramOpts{
			Name: "cortex_ingester_queried_chunks",
			Help: "The total number of chunks returned from queries.",
			// A small number of chunks per series - 10*(8^(7-1)) = 2.6m.
			Buckets: prometheus.ExponentialBuckets(10, 8, 7),
		}),
		memSeries: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "cortex_ingester_memory_series",
			Help: "The current number of series in memory.",
		}),
		memUsers: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "cortex_ingester_memory_users",
			Help: "The current number of users in memory.",
		}),
		memSeriesCreatedTotal: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: memSeriesCreatedTotalName,
			Help: memSeriesCreatedTotalHelp,
		}, []string{"user"}),
		memSeriesRemovedTotal: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: memSeriesRemovedTotalName,
			Help: memSeriesRemovedTotalHelp,
		}, []string{"user"}),
	}

	if r != nil {
		r.MustRegister(
			m.flushQueueLength,
			m.ingestedSamples,
			m.ingestedSamplesFail,
			m.queries,
			m.queriedSamples,
			m.queriedSeries,
			m.queriedChunks,
			m.memSeries,
			m.memUsers,
		)

		if registerMetricsConflictingWithTSDB {
			r.MustRegister(
				m.memSeriesCreatedTotal,
				m.memSeriesRemovedTotal,
			)
		}
	}

	return m
}

// TSDB metrics. Each tenant has its own registry, that TSDB code uses.
type tsdbMetrics struct {
	// We aggregate metrics from individual TSDB registries into
	// a single set of counters, which are exposed as Cortex metrics.
	dirSyncs        *prometheus.Desc // sum(thanos_shipper_dir_syncs_total)
	dirSyncFailures *prometheus.Desc // sum(thanos_shipper_dir_sync_failures_total)
	uploads         *prometheus.Desc // sum(thanos_shipper_uploads_total)
	uploadFailures  *prometheus.Desc // sum(thanos_shipper_upload_failures_total)

	// These two metrics replace metrics in ingesterMetrics, as we count them differently
	memSeriesCreatedTotal *prometheus.Desc
	memSeriesRemovedTotal *prometheus.Desc

	// These maps drive the collection output. Key = original metric name to group.
	sumCountersGlobally map[string]*prometheus.Desc
	sumCountersPerUser  map[string]*prometheus.Desc

	regsMu sync.RWMutex                    // custom mutex for shipper registry, to avoid blocking main user state mutex on collection
	regs   map[string]*prometheus.Registry // One prometheus registry per tenant
}

func newTSDBMetrics(r prometheus.Registerer) *tsdbMetrics {
	m := &tsdbMetrics{
		regs: make(map[string]*prometheus.Registry),

		dirSyncs: prometheus.NewDesc(
			"cortex_ingester_shipper_dir_syncs_total",
			"TSDB: Total dir sync attempts",
			nil, nil),
		dirSyncFailures: prometheus.NewDesc(
			"cortex_ingester_shipper_dir_sync_failures_total",
			"TSDB: Total number of failed dir syncs",
			nil, nil),
		uploads: prometheus.NewDesc(
			"cortex_ingester_shipper_uploads_total",
			"TSDB: Total object upload attempts",
			nil, nil),
		uploadFailures: prometheus.NewDesc(
			"cortex_ingester_shipper_upload_failures_total",
			"TSDB: Total number of failed object uploads",
			nil, nil),

		memSeriesCreatedTotal: prometheus.NewDesc(memSeriesCreatedTotalName, memSeriesCreatedTotalHelp, []string{"user"}, nil),
		memSeriesRemovedTotal: prometheus.NewDesc(memSeriesRemovedTotalName, memSeriesRemovedTotalHelp, []string{"user"}, nil),
	}

	m.sumCountersGlobally = map[string]*prometheus.Desc{
		"thanos_shipper_dir_syncs_total":         m.dirSyncs,
		"thanos_shipper_dir_sync_failures_total": m.dirSyncFailures,
		"thanos_shipper_uploads_total":           m.uploads,
		"thanos_shipper_upload_failures_total":   m.uploadFailures,
	}

	m.sumCountersPerUser = map[string]*prometheus.Desc{
		"prometheus_tsdb_head_series_created_total": m.memSeriesCreatedTotal,
		"prometheus_tsdb_head_series_removed_total": m.memSeriesRemovedTotal,
	}

	if r != nil {
		r.MustRegister(m)
	}
	return m
}

func (sm *tsdbMetrics) Describe(out chan<- *prometheus.Desc) {
	out <- sm.dirSyncs
	out <- sm.dirSyncFailures
	out <- sm.uploads
	out <- sm.uploadFailures
	out <- sm.memSeriesCreatedTotal
	out <- sm.memSeriesRemovedTotal
}

func (sm *tsdbMetrics) Collect(out chan<- prometheus.Metric) {
	regs := sm.registries()
	data := util.NewMetricFamiliersPerUser()

	for userID, r := range regs {
		m, err := r.Gather()
		if err == nil {
			err = data.AddGatheredDataForUser(userID, m)
		}
		if err != nil {
			level.Warn(util.Logger).Log("msg", "failed to gather metrics from TSDB shipper", "user", userID, "err", err)
			continue
		}
	}

	// OK, we have it all. Let's build results.
	for metric, desc := range sm.sumCountersGlobally {
		out <- prometheus.MustNewConstMetric(desc, prometheus.CounterValue, data.SumCountersAcrossAllUsers(metric))
	}

	for metric, desc := range sm.sumCountersPerUser {
		userValues := data.SumCountersPerUser(metric)
		for user, val := range userValues {
			out <- prometheus.MustNewConstMetric(desc, prometheus.CounterValue, val, user)
		}
	}
}

// make a copy of the map, so that metrics can be gathered while the new registry is being added.
func (sm *tsdbMetrics) registries() map[string]*prometheus.Registry {
	sm.regsMu.RLock()
	defer sm.regsMu.RUnlock()

	regs := make(map[string]*prometheus.Registry, len(sm.regs))
	for u, r := range sm.regs {
		regs[u] = r
	}
	return regs
}

func (sm *tsdbMetrics) setRegistryForUser(userID string, registry *prometheus.Registry) {
	sm.regsMu.Lock()
	sm.regs[userID] = registry
	sm.regsMu.Unlock()
}
