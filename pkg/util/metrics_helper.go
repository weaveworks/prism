package util

import (
	"bytes"
	"errors"
	"fmt"
	"sync"

	"github.com/go-kit/kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
)

type SingleValueWithLabels struct {
	Value       float64
	LabelValues []string
}

type SingleValueWithLabelsMap map[string]SingleValueWithLabels

func (m SingleValueWithLabelsMap) aggregateFn(labelsKey string, labelValues []string, value float64) {
	r := m[labelsKey]
	if r.LabelValues == nil {
		r.LabelValues = labelValues
	}

	r.Value += value
	m[labelsKey] = r
}

func (m SingleValueWithLabelsMap) WriteToMetricChannel(out chan<- prometheus.Metric, desc *prometheus.Desc, valueType prometheus.ValueType) {
	for _, cr := range m {
		out <- prometheus.MustNewConstMetric(desc, valueType, cr.Value, cr.LabelValues...)
	}
}

// MetricFamilyMap is a map of metric names to their family (metrics with same name, but different labels)
// Keeping map of metric name to its family makes it easier to do searches later.
type MetricFamilyMap map[string]*dto.MetricFamily

// NewMetricFamilyMap sorts output from Gatherer.Gather method into a map.
// Gatherer.Gather specifies that there metric families are uniquely named, and we use that fact here.
// If they are not, this method returns error.
func NewMetricFamilyMap(metrics []*dto.MetricFamily) (MetricFamilyMap, error) {
	perMetricName := MetricFamilyMap{}

	for _, m := range metrics {
		name := m.GetName()
		// these errors should never happen when passing Gatherer.Gather() output.
		if name == "" {
			return nil, errors.New("empty name for metric family")
		}
		if perMetricName[name] != nil {
			return nil, fmt.Errorf("non-unique name for metric family: %q", name)
		}

		perMetricName[name] = m
	}

	return perMetricName, nil
}

func (mfm MetricFamilyMap) SumCounters(name string) float64 {
	return sum(mfm[name], counterValue)
}

func (mfm MetricFamilyMap) SumCountersWithLabels(name string, labelNames ...string) SingleValueWithLabelsMap {
	result := SingleValueWithLabelsMap{}
	mfm.sumOfSingleValuesWithLabels(name, labelNames, counterValue, result.aggregateFn)
	return result
}

func (mfm MetricFamilyMap) SumGauges(name string) float64 {
	return sum(mfm[name], gaugeValue)
}

func (mfm MetricFamilyMap) SumGaugesWithLabels(name string, labelNames ...string) map[string]SingleValueWithLabels {
	result := SingleValueWithLabelsMap{}
	mfm.sumOfSingleValuesWithLabels(name, labelNames, gaugeValue, result.aggregateFn)
	return result
}

func (mfm MetricFamilyMap) sumOfSingleValuesWithLabels(metric string, labelNames []string, extractFn func(*dto.Metric) float64, aggregateFn func(labelsKey string, labelValues []string, value float64)) {
	metricsPerLabelValue := getMetricsWithLabelNames(mfm[metric], labelNames)

	for key, mlv := range metricsPerLabelValue {
		for _, m := range mlv.metrics {
			val := extractFn(m)
			aggregateFn(key, mlv.labelValues, val)
		}
	}
}

func (mfm MetricFamilyMap) SumHistograms(s string) HistogramData {
	hd := HistogramData{}
	mfm.SumHistogramsTo(s, &hd)
	return hd
}

func (mfm MetricFamilyMap) SumHistogramsTo(name string, hd *HistogramData) {
	for _, m := range mfm[name].GetMetric() {
		hd.AddHistogram(m.GetHistogram())
	}
}

func (mfm MetricFamilyMap) SumSummaries(s string) SummaryData {
	sd := SummaryData{}
	mfm.SumSummariesTo(s, &sd)
	return sd
}

func (mfm MetricFamilyMap) SumSummariesTo(name string, hd *SummaryData) {
	for _, m := range mfm[name].GetMetric() {
		hd.AddSummary(m.GetSummary())
	}
}

// MetricFamiliesPerUser is a collection of metrics gathered via calling Gatherer.Gather() method on different
// gatherers, one per user.
type MetricFamiliesPerUser map[string]MetricFamilyMap

func BuildMetricFamiliesPerUserFromUserRegistries(regs map[string]*prometheus.Registry) MetricFamiliesPerUser {
	data := MetricFamiliesPerUser{}
	for userID, r := range regs {
		m, err := r.Gather()
		if err == nil {
			var mfm MetricFamilyMap = nil
			mfm, err = NewMetricFamilyMap(m)
			if err == nil {
				data[userID] = mfm
			}
		}

		if err != nil {
			level.Warn(Logger).Log("msg", "failed to gather metrics from registry", "user", userID, "err", err)
			continue
		}
	}
	return data
}

func (d MetricFamiliesPerUser) SendSumOfCounters(out chan<- prometheus.Metric, desc *prometheus.Desc, counter string) {
	result := float64(0)
	for _, perUser := range d {
		result += perUser.SumCounters(counter)
	}
	out <- prometheus.MustNewConstMetric(desc, prometheus.CounterValue, result)
}

func (d MetricFamiliesPerUser) SendSumOfCountersWithLabels(out chan<- prometheus.Metric, desc *prometheus.Desc, counter string, labelNames ...string) {
	d.sumOfSingleValuesWithLabels(counter, counterValue, labelNames).WriteToMetricChannel(out, desc, prometheus.CounterValue)
}

func (d MetricFamiliesPerUser) SendSumOfCountersPerUser(out chan<- prometheus.Metric, desc *prometheus.Desc, counter string) {
	for user, perMetric := range d {
		v := perMetric.SumCounters(counter)

		out <- prometheus.MustNewConstMetric(desc, prometheus.CounterValue, v, user)
	}
}

func (d MetricFamiliesPerUser) SendSumOfGauges(out chan<- prometheus.Metric, desc *prometheus.Desc, gauge string) {
	result := float64(0)
	for _, perMetric := range d {
		result += perMetric.SumGauges(gauge)
	}
	out <- prometheus.MustNewConstMetric(desc, prometheus.GaugeValue, result)
}

func (d MetricFamiliesPerUser) SendSumOfGaugesWithLabels(out chan<- prometheus.Metric, desc *prometheus.Desc, gauge string, labelNames ...string) {
	d.sumOfSingleValuesWithLabels(gauge, gaugeValue, labelNames).WriteToMetricChannel(out, desc, prometheus.GaugeValue)
}

func (d MetricFamiliesPerUser) sumOfSingleValuesWithLabels(metric string, fn func(*dto.Metric) float64, labelNames []string) SingleValueWithLabelsMap {
	result := SingleValueWithLabelsMap{}
	for _, userMetrics := range d {
		userMetrics.sumOfSingleValuesWithLabels(metric, labelNames, fn, result.aggregateFn)
	}
	return result
}

func (d MetricFamiliesPerUser) SendSumOfSummaries(out chan<- prometheus.Metric, desc *prometheus.Desc, summaryName string) {
	summaryData := SummaryData{}
	for _, userMetrics := range d {
		userMetrics.SumSummariesTo(summaryName, &summaryData)
	}
	out <- summaryData.Metric(desc)
}

func (d MetricFamiliesPerUser) SendSumOfSummariesWithLabels(out chan<- prometheus.Metric, desc *prometheus.Desc, summaryName string, labelNames ...string) {
	type summaryResult struct {
		data        SummaryData
		labelValues []string
	}

	result := map[string]summaryResult{}

	for _, userMetrics := range d {
		metricsPerLabelValue := getMetricsWithLabelNames(userMetrics[summaryName], labelNames)

		for key, mwl := range metricsPerLabelValue {
			for _, m := range mwl.metrics {
				r := result[key]
				if r.labelValues == nil {
					r.labelValues = mwl.labelValues
				}

				r.data.AddSummary(m.GetSummary())
				result[key] = r
			}
		}
	}

	for _, sr := range result {
		out <- sr.data.Metric(desc, sr.labelValues...)
	}
}

func (d MetricFamiliesPerUser) SendSumOfHistograms(out chan<- prometheus.Metric, desc *prometheus.Desc, histogramName string) {
	hd := HistogramData{}
	for _, userMetrics := range d {
		userMetrics.SumHistogramsTo(histogramName, &hd)
	}
	out <- hd.Metric(desc)
}

type metricsWithLabels struct {
	labelValues []string
	metrics     []*dto.Metric
}

func getMetricsWithLabelNames(mf *dto.MetricFamily, labelNames []string) map[string]metricsWithLabels {
	result := map[string]metricsWithLabels{}

	for _, m := range mf.GetMetric() {
		lbls, include := getLabelValues(m, labelNames)
		if !include {
			continue
		}

		key := getLabelsString(lbls)
		r := result[key]
		if r.labelValues == nil {
			r.labelValues = lbls
		}
		r.metrics = append(r.metrics, m)
		result[key] = r
	}
	return result
}

func getLabelValues(m *dto.Metric, labelNames []string) ([]string, bool) {
	all := map[string]string{}
	for _, lp := range m.GetLabel() {
		all[lp.GetName()] = lp.GetValue()
	}

	result := make([]string, 0, len(labelNames))
	for _, ln := range labelNames {
		lv, ok := all[ln]
		if !ok {
			// required labels not found
			return nil, false
		}
		result = append(result, lv)
	}
	return result, true
}

func getLabelsString(labelValues []string) string {
	buf := bytes.Buffer{}
	for _, v := range labelValues {
		buf.WriteString(v)
		buf.WriteByte(0) // separator, not used in prometheus labels
	}
	return buf.String()
}

// sum returns sum of values from all metrics from same metric family (= series with the same metric name, but different labels)
// Supplied function extracts value.
func sum(mf *dto.MetricFamily, fn func(*dto.Metric) float64) float64 {
	result := float64(0)
	for _, m := range mf.GetMetric() {
		result += fn(m)
	}
	return result
}

// This works even if m is nil, m.Counter is nil or m.Counter.Value is nil (it returns 0 in those cases)
func counterValue(m *dto.Metric) float64 { return m.GetCounter().GetValue() }
func gaugeValue(m *dto.Metric) float64   { return m.GetGauge().GetValue() }

// SummaryData keeps all data needed to create summary metric
type SummaryData struct {
	sampleCount uint64
	sampleSum   float64
	quantiles   map[float64]float64
}

func (s *SummaryData) AddSummary(sum *dto.Summary) {
	s.sampleCount += sum.GetSampleCount()
	s.sampleSum += sum.GetSampleSum()

	qs := sum.GetQuantile()
	if len(qs) > 0 && s.quantiles == nil {
		s.quantiles = map[float64]float64{}
	}

	for _, q := range qs {
		// we assume that all summaries have same quantiles
		s.quantiles[q.GetQuantile()] += q.GetValue()
	}
}

func (s *SummaryData) Metric(desc *prometheus.Desc, labelValues ...string) prometheus.Metric {
	return prometheus.MustNewConstSummary(desc, s.sampleCount, s.sampleSum, s.quantiles, labelValues...)
}

type HistogramData struct {
	sampleCount uint64
	sampleSum   float64
	buckets     map[float64]uint64
}

func (d *HistogramData) AddHistogram(histo *dto.Histogram) {
	d.sampleCount += histo.GetSampleCount()
	d.sampleSum += histo.GetSampleSum()

	histoBuckets := histo.GetBucket()
	if len(histoBuckets) > 0 && d.buckets == nil {
		d.buckets = map[float64]uint64{}
	}

	for _, b := range histoBuckets {
		// we assume that all histograms have same buckets
		d.buckets[b.GetUpperBound()] += b.GetCumulativeCount()
	}
}

func (d *HistogramData) AddHistogramData(histo HistogramData) {
	d.sampleCount += histo.sampleCount
	d.sampleSum += histo.sampleSum

	if len(histo.buckets) > 0 && d.buckets == nil {
		d.buckets = map[float64]uint64{}
	}

	for bound, count := range histo.buckets {
		// we assume that all histograms have same buckets
		d.buckets[bound] += count
	}
}

func (d *HistogramData) Metric(desc *prometheus.Desc) prometheus.Metric {
	return prometheus.MustNewConstHistogram(desc, d.sampleCount, d.sampleSum, d.buckets)
}

func NewHistogramDataCollector(desc *prometheus.Desc) *HistogramDataCollector {
	return &HistogramDataCollector{
		desc: desc,
		data: &HistogramData{},
	}
}

type HistogramDataCollector struct {
	desc *prometheus.Desc

	dataMu sync.RWMutex
	data   *HistogramData
}

func (h *HistogramDataCollector) Describe(out chan<- *prometheus.Desc) {
	out <- h.desc
}

func (h *HistogramDataCollector) Collect(out chan<- prometheus.Metric) {
	h.dataMu.RLock()
	defer h.dataMu.RUnlock()

	out <- h.data.Metric(h.desc)
}

func (h *HistogramDataCollector) Add(hd HistogramData) {
	h.dataMu.Lock()
	defer h.dataMu.Unlock()

	h.data.AddHistogramData(hd)
}
