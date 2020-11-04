package worker

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/atomic"
	"google.golang.org/grpc"

	"github.com/cortexproject/cortex/pkg/util"
	"github.com/cortexproject/cortex/pkg/util/test"
)

func TestRecvFailDoesntCancelProcess(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// We use random port here, hopefully without any gRPC server.
	cc, err := grpc.DialContext(ctx, "localhost:999", grpc.WithInsecure())
	require.NoError(t, err)

	cfg := Config{}
	mgr := newFrontendProcessor(cfg, nil, util.Logger)
	running := atomic.NewBool(false)
	go func() {
		running.Store(true)
		defer running.Store(false)

		mgr.processQueriesOnSingleStream(ctx, cc, "test:12345")
	}()

	time.Sleep(50 * time.Millisecond)
	assert.Equal(t, true, running.Load())

	cancel()
	test.Poll(t, 100*time.Millisecond, false, func() interface{} {
		return running.Load()
	})
}

func TestServeCancelStopsProcess(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// We use random port here, hopefully without any gRPC server.
	cc, err := grpc.DialContext(ctx, "localhost:999", grpc.WithInsecure())
	require.NoError(t, err)

	pm := newProcessorManager(ctx, &mockProcessor{}, cc, "test")
	pm.concurrency(1)

	test.Poll(t, 100*time.Millisecond, 1, func() interface{} {
		return int(pm.currentProcessors.Load())
	})

	cancel()

	test.Poll(t, 100*time.Millisecond, 0, func() interface{} {
		return int(pm.currentProcessors.Load())
	})

	pm.stop()
	test.Poll(t, 100*time.Millisecond, 0, func() interface{} {
		return int(pm.currentProcessors.Load())
	})
}