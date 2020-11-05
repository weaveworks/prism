package ring

import (
	"context"
	"errors"
	"sort"
	"time"
)

var errorTooManyZoneFailures = errors.New("too many zones failed")

// ReplicationSet describes the ingesters to talk to for a given key, and how
// many errors to tolerate.
type ReplicationSet struct {
	Ingesters           []IngesterDesc
	MaxErrors           int
	MaxUnavailableZones int
}

// Do function f in parallel for all replicas in the set, erroring is we exceed
// MaxErrors and returning early otherwise.
func (r ReplicationSet) Do(ctx context.Context, delay time.Duration, f func(context.Context, *IngesterDesc) (interface{}, error)) ([]interface{}, error) {
	type instanceError struct {
		err      error
		instance *IngesterDesc
	}

	var (
		errs        = make(chan instanceError, len(r.Ingesters))
		resultsChan = make(chan interface{}, len(r.Ingesters))
		minSuccess  = len(r.Ingesters) - r.MaxErrors
		forceStart  = make(chan struct{}, r.MaxErrors)
	)
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	for i := range r.Ingesters {
		go func(i int, ing *IngesterDesc) {
			// wait to send extra requests
			if i >= minSuccess && delay > 0 {
				after := time.NewTimer(delay)
				defer after.Stop()
				select {
				case <-ctx.Done():
					return
				case <-forceStart:
				case <-after.C:
				}
			}
			result, err := f(ctx, ing)
			if err != nil {
				errs <- instanceError{
					err:      err,
					instance: ing,
				}
			} else {
				resultsChan <- result
			}
		}(i, &r.Ingesters[i])
	}

	var (
		numErrs          int
		numSuccess       int
		results          = make([]interface{}, 0, len(r.Ingesters))
		zoneFailureCount = make(map[string]struct{})
	)
	for numSuccess < minSuccess {
		select {
		case err := <-errs:
			if r.MaxUnavailableZones > 0 {
				zoneFailureCount[err.instance.Zone] = struct{}{}

				if len(zoneFailureCount) > r.MaxUnavailableZones {
					return nil, errorTooManyZoneFailures
				}
			} else {
				numErrs++
				if numErrs > r.MaxErrors {
					return nil, err.err
				}
			}

			// force one of the delayed requests to start
			forceStart <- struct{}{}

		case result := <-resultsChan:
			numSuccess++
			results = append(results, result)

		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}

	return results, nil
}

// Includes returns whether the replication set includes the replica with the provided addr.
func (r ReplicationSet) Includes(addr string) bool {
	for _, instance := range r.Ingesters {
		if instance.GetAddr() == addr {
			return true
		}
	}

	return false
}

// GetAddresses returns the addresses of all instances within the replication set. Returned slice
// order is not guaranteed.
func (r ReplicationSet) GetAddresses() []string {
	addrs := make([]string, 0, len(r.Ingesters))
	for _, desc := range r.Ingesters {
		addrs = append(addrs, desc.Addr)
	}
	return addrs
}

// HasReplicationSetChanged returns true if two replications sets are the same (with possibly different timestamps),
// false if they differ in any way (number of instances, instance states, tokens, zones, ...).
func HasReplicationSetChanged(before, after ReplicationSet) bool {
	beforeInstances := before.Ingesters
	afterInstances := after.Ingesters

	if len(beforeInstances) != len(afterInstances) {
		return true
	}

	sort.Sort(ByAddr(beforeInstances))
	sort.Sort(ByAddr(afterInstances))

	for i := 0; i < len(beforeInstances); i++ {
		b := beforeInstances[i]
		a := afterInstances[i]

		// Exclude the heartbeat timestamp from the comparison.
		b.Timestamp = 0
		a.Timestamp = 0

		if !b.Equal(a) {
			return true
		}
	}

	return false
}
