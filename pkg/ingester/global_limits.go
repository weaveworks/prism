package ingester

import "fmt"

// GlobalLimits describes global limits used by ingester. Reaching any of these will result in Push method to return
// (internal) error.
type GlobalLimits struct {
	MaxIngestionRate  float64 `yaml:"max_ingestion_rate"`
	MaxInMemoryUsers  int64   `yaml:"max_users"`
	MaxInMemorySeries int64   `yaml:"max_series"`
}

// Sets default limit values for unmarshalling.
var defaultGlobalLimits *GlobalLimits = nil

// UnmarshalYAML implements the yaml.Unmarshaler interface. If give
func (l *GlobalLimits) UnmarshalYAML(unmarshal func(interface{}) error) error {
	if defaultGlobalLimits != nil {
		*l = *defaultGlobalLimits
	}
	type plain GlobalLimits // type indirection to make sure we don't go into recursive loop
	return unmarshal((*plain)(l))
}

type errMaxSamplesPushRateLimitReached struct {
	rate  float64
	limit float64
}

func (e errMaxSamplesPushRateLimitReached) Error() string {
	return fmt.Sprintf("cannot push more samples: ingester's max samples push rate reached, rate=%g, limit=%g", e.rate, e.limit)
}

type errMaxUsersLimitReached struct {
	users int64
	limit int64
}

func (e errMaxUsersLimitReached) Error() string {
	return fmt.Sprintf("cannot create TSDB: ingesters's max users limit reached, users=%d, limit=%d", e.users, e.limit)
}

type errMaxSeriesLimitReached struct {
	series int64
	limit  int64
}

func (e errMaxSeriesLimitReached) Error() string {
	return fmt.Sprintf("cannot add series: ingesters's max series limit reached, series=%d, limit=%d", e.series, e.limit)
}
