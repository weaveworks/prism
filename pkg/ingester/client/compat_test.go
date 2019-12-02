package client

import (
	"fmt"
	"reflect"
	"sort"
	"testing"

	"github.com/prometheus/common/model"
	"github.com/prometheus/prometheus/pkg/labels"
)

func TestQueryRequest(t *testing.T) {
	from, to := model.Time(int64(0)), model.Time(int64(10))
	matchers := []*labels.Matcher{}
	matcher1, err := labels.NewMatcher(labels.MatchEqual, "foo", "1")
	if err != nil {
		t.Fatal(err)
	}
	matchers = append(matchers, matcher1)

	matcher2, err := labels.NewMatcher(labels.MatchNotEqual, "bar", "2")
	if err != nil {
		t.Fatal(err)
	}
	matchers = append(matchers, matcher2)

	matcher3, err := labels.NewMatcher(labels.MatchRegexp, "baz", "3")
	if err != nil {
		t.Fatal(err)
	}
	matchers = append(matchers, matcher3)

	matcher4, err := labels.NewMatcher(labels.MatchNotRegexp, "bop", "4")
	if err != nil {
		t.Fatal(err)
	}
	matchers = append(matchers, matcher4)

	req, err := ToQueryRequest(from, to, matchers)
	if err != nil {
		t.Fatal(err)
	}

	haveFrom, haveTo, haveMatchers, err := FromQueryRequest(req)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(haveFrom, from) {
		t.Fatalf("Bad from FromQueryRequest(ToQueryRequest) round trip")
	}
	if !reflect.DeepEqual(haveTo, to) {
		t.Fatalf("Bad to FromQueryRequest(ToQueryRequest) round trip")
	}
	if !reflect.DeepEqual(haveMatchers, matchers) {
		t.Fatalf("Bad have FromQueryRequest(ToQueryRequest) round trip - %v != %v", haveMatchers, matchers)
	}
}

func buildTestMatrix(numSeries int, samplesPerSeries int, offset int) model.Matrix {
	m := make(model.Matrix, 0, numSeries)
	for i := 0; i < numSeries; i++ {
		ss := model.SampleStream{
			Metric: model.Metric{
				model.MetricNameLabel: model.LabelValue(fmt.Sprintf("testmetric_%d", i)),
				model.JobLabel:        "testjob",
			},
			Values: make([]model.SamplePair, 0, samplesPerSeries),
		}
		for j := 0; j < samplesPerSeries; j++ {
			ss.Values = append(ss.Values, model.SamplePair{
				Timestamp: model.Time(i + j + offset),
				Value:     model.SampleValue(i + j + offset),
			})
		}
		m = append(m, &ss)
	}
	sort.Sort(m)
	return m
}

func TestQueryResponse(t *testing.T) {
	want := buildTestMatrix(10, 10, 10)
	have := FromQueryResponse(ToQueryResponse(want))
	if !reflect.DeepEqual(have, want) {
		t.Fatalf("Bad FromQueryResponse(ToQueryResponse) round trip")
	}

}

func TestFingerprintCollisions(t *testing.T) {
	// "8yn0iYCKYHlIj4-BwPqk" and "GReLUrM4wMqfg9yzV3KQ" have same FNV-1a hash.
	// If we use it as a single label name (for labels that have same value), we get colliding labels.
	c1 := labels.FromStrings("8yn0iYCKYHlIj4-BwPqk", "hello")
	c2 := labels.FromStrings("GReLUrM4wMqfg9yzV3KQ", "hello")
	verifyCollision(t, true, c1, c2)

	// Adding _="ypfajYg2lsv" or _="KiqbryhzUpn" respectively to most [*] metrics will produce collision.
	// It's because "_\xffypfajYg2lsv" and "_\xffKiqbryhzUpn" have same FNV-1a hash, and "_" label is sorted before
	// most other labels (except labels starting with upper-case letter)

	const _label1 = "ypfajYg2lsv"
	const _label2 = "KiqbryhzUpn"

	metric := labels.NewBuilder(labels.FromStrings("__name__", "logs"))
	c1 = metric.Set("_", _label1).Labels()
	c2 = metric.Set("_", _label2).Labels()
	verifyCollision(t, true, c1, c2)

	metric = labels.NewBuilder(labels.FromStrings("__name__", "up", "instance", "hello"))
	c1 = metric.Set("_", _label1).Labels()
	c2 = metric.Set("_", _label2).Labels()
	verifyCollision(t, true, c1, c2)

	// here it breaks, because "Z" label is sorted before "_" label.
	metric = labels.NewBuilder(labels.FromStrings("__name__", "up", "Z", "hello"))
	c1 = metric.Set("_", _label1).Labels()
	c2 = metric.Set("_", _label2).Labels()
	verifyCollision(t, false, c1, c2)

	// But A="K6sjsNNczPl" and A="cswpLMIZpwt" label has similar property.
	// (Again, "A\xffK6sjsNNczPl" and "A\xffcswpLMIZpwt" have same FNV-1a hash)
	// This time, "A" is the smallest possible label name, and is always sorted first.

	const Alabel1 = "K6sjsNNczPl"
	const Alabel2 = "cswpLMIZpwt"

	metric = labels.NewBuilder(labels.FromStrings("__name__", "up", "Z", "hello"))
	c1 = metric.Set("A", Alabel1).Labels()
	c2 = metric.Set("A", Alabel2).Labels()
	verifyCollision(t, true, c1, c2)
}

func verifyCollision(t *testing.T, collision bool, ls1 labels.Labels, ls2 labels.Labels) {
	if collision && Fingerprint(ls1) != Fingerprint(ls2) {
		t.Errorf("expected same fingerprints for %v (%016x) and %v (%016x)", ls1.String(), Fingerprint(ls1), ls2.String(), Fingerprint(ls2))
	} else if !collision && Fingerprint(ls1) == Fingerprint(ls2) {
		t.Errorf("expected different fingerprints for %v (%016x) and %v (%016x)", ls1.String(), Fingerprint(ls1), ls2.String(), Fingerprint(ls2))
	}
}
