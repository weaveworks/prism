package queryrange

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"
	"time"

	"github.com/cortexproject/cortex/pkg/chunk"
	"github.com/cortexproject/cortex/pkg/util"
	"github.com/go-kit/kit/log"
	"github.com/prometheus/prometheus/promql"
	"github.com/stretchr/testify/require"
	"github.com/weaveworks/common/middleware"
	"github.com/weaveworks/common/user"
)

func TestRoundTrip(t *testing.T) {
	s := httptest.NewServer(
		middleware.AuthenticateUser.Wrap(
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				var err error
				if r.RequestURI == query {
					_, err = w.Write([]byte(responseBody))
				} else {
					_, err = w.Write([]byte("bar"))
				}
				if err != nil {
					t.Fatal(err)
				}
			}),
		),
	)
	defer s.Close()

	u, err := url.Parse(s.URL)
	require.NoError(t, err)

	downstream := singleHostRoundTripper{
		host: u.Host,
		next: http.DefaultTransport,
	}

	tw, _, err := NewTripperware(Config{},
		util.Logger,
		fakeLimits{},
		PrometheusCodec,
		nil,
		chunk.SchemaConfig{},
		promql.EngineOpts{
			Logger:        util.Logger,
			Reg:           nil,
			MaxConcurrent: 10,
			MaxSamples:    1000,
			Timeout:       time.Minute,
		},
		0,
	)

	if err != nil {
		t.Fatal(err)
	}

	for i, tc := range []struct {
		path, expectedBody string
	}{
		{"/foo", "bar"},
		{query, responseBody},
	} {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			req, err := http.NewRequest("GET", tc.path, http.NoBody)
			require.NoError(t, err)

			// query-frontend doesn't actually authenticate requests, we rely on
			// the queriers to do this.  Hence we ensure the request doesn't have a
			// org ID in the ctx, but does have the header.
			ctx := user.InjectOrgID(context.Background(), "1")
			req = req.WithContext(ctx)
			err = user.InjectOrgIDIntoHTTPRequest(ctx, req)
			require.NoError(t, err)

			resp, err := tw(downstream).RoundTrip(req)
			require.NoError(t, err)
			require.Equal(t, 200, resp.StatusCode)

			bs, err := ioutil.ReadAll(resp.Body)
			require.NoError(t, err)
			require.Equal(t, tc.expectedBody, string(bs))
		})
	}
}

type singleHostRoundTripper struct {
	host string
	next http.RoundTripper
}

func (s singleHostRoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	r.URL.Scheme = "http"
	r.URL.Host = s.host
	return s.next.RoundTrip(r)
}

func Test_ShardingConfigError(t *testing.T) {
	_, _, err := NewTripperware(
		Config{SumShards: true},
		log.NewNopLogger(),
		nil,
		nil,
		nil,
		chunk.SchemaConfig{},
		promql.EngineOpts{},
		0,
	)

	require.EqualError(t, err, "a non-zero value is required for querier.query-ingesters-within when querier.sum-shards is enabled")
}
