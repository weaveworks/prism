package storage

import (
	"context"
	"fmt"
	"math/rand"
	"sort"
	"strconv"
	"testing"
	"time"

	"github.com/prometheus/common/model"
	"github.com/stretchr/testify/require"

	"github.com/weaveworks/cortex/pkg/chunk"
)

func TestChunksBasic(t *testing.T) {
	forAllFixtures(t, func(t *testing.T, client chunk.StorageClient) {
		const batchSize = 50
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// Write a few batches of chunks.
		written := []string{}
		for i := 0; i < 50; i++ {
			chunks := []chunk.Chunk{}
			for j := 0; j < batchSize; j++ {
				chunk := dummyChunkFor(model.Now(), model.Metric{
					model.MetricNameLabel: "foo",
					"index":               model.LabelValue(strconv.Itoa(i*batchSize + j)),
				})
				chunks = append(chunks, chunk)
				_, err := chunk.Encode() // Need to encode it, side effect calculates crc
				require.NoError(t, err)
				written = append(written, chunk.ExternalKey())
			}

			err := client.PutChunks(ctx, chunks)
			require.NoError(t, err)
		}

		// Get a few batches of chunks.
		for i := 0; i < 50; i++ {
			keysToGet := map[string]struct{}{}
			chunksToGet := []chunk.Chunk{}
			for len(chunksToGet) < batchSize {
				key := written[rand.Intn(len(written))]
				if _, ok := keysToGet[key]; ok {
					continue
				}
				keysToGet[key] = struct{}{}
				chunk, err := chunk.ParseExternalKey(userID, key)
				require.NoError(t, err)
				chunksToGet = append(chunksToGet, chunk)
			}

			chunksWeGot, err := client.GetChunks(ctx, chunksToGet)
			require.NoError(t, err)
			require.Equal(t, len(chunksToGet), len(chunksWeGot))

			sort.Sort(chunk.ByKey(chunksToGet))
			sort.Sort(chunk.ByKey(chunksWeGot))
			for j := 0; j < len(chunksWeGot); j++ {
				require.Equal(t, chunksToGet[i].ExternalKey(), chunksWeGot[i].ExternalKey(), strconv.Itoa(i))
			}
		}
	})
}

func TestIndexBasic(t *testing.T) {
	forAllFixtures(t, func(t *testing.T, client chunk.StorageClient) {
		// Write out 30 entries, into different hash and range values.
		batch := client.NewWriteBatch()
		for i := 0; i < 30; i++ {
			batch.Add(tableName, fmt.Sprintf("hash%d", i), []byte(fmt.Sprintf("range%d", i)), nil)
		}
		err := client.BatchWrite(context.Background(), batch)
		require.NoError(t, err)

		// Make sure we get back the correct entries by hash value.
		for i := 0; i < 30; i++ {
			entry := chunk.IndexQuery{
				TableName: tableName,
				HashValue: fmt.Sprintf("hash%d", i),
			}
			var have []chunk.IndexEntry
			err := client.QueryPages(context.Background(), entry, func(read chunk.ReadBatch, lastPage bool) bool {
				for j := 0; j < read.Len(); j++ {
					have = append(have, chunk.IndexEntry{
						RangeValue: read.RangeValue(j),
					})
				}
				return !lastPage
			})
			require.NoError(t, err)
			require.Equal(t, []chunk.IndexEntry{
				{RangeValue: []byte(fmt.Sprintf("range%d", i))},
			}, have)
		}
	})
}

var entries = []chunk.IndexEntry{
	{
		TableName:  "table",
		HashValue:  "foo",
		RangeValue: []byte("bar:1"),
		Value:      []byte("10"),
	},
	{
		TableName:  "table",
		HashValue:  "foo",
		RangeValue: []byte("bar:2"),
		Value:      []byte("20"),
	},
	{
		TableName:  "table",
		HashValue:  "foo",
		RangeValue: []byte("bar:3"),
		Value:      []byte("30"),
	},
	{
		TableName:  "table",
		HashValue:  "foo",
		RangeValue: []byte("baz:1"),
		Value:      []byte("10"),
	},
	{
		TableName:  "table",
		HashValue:  "foo",
		RangeValue: []byte("baz:2"),
		Value:      []byte("20"),
	},
	{
		TableName:  "table",
		HashValue:  "flip",
		RangeValue: []byte("bar:1"),
		Value:      []byte("abc"),
	},
	{
		TableName:  "table",
		HashValue:  "flip",
		RangeValue: []byte("bar:2"),
		Value:      []byte("abc"),
	},
	{
		TableName:  "table",
		HashValue:  "flip",
		RangeValue: []byte("bar:3"),
		Value:      []byte("abc"),
	},
}

func TestQueryPages(t *testing.T) {
	forAllFixtures(t, func(t *testing.T, client chunk.StorageClient) {
		batch := client.NewWriteBatch()
		for _, entry := range entries {
			batch.Add(entry.TableName, entry.HashValue, entry.RangeValue, entry.Value)
		}

		err := client.BatchWrite(context.Background(), batch)
		require.NoError(t, err)

		tests := []struct {
			name           string
			query          chunk.IndexQuery
			provisionedErr int
			want           []chunk.IndexEntry
		}{
			{
				"check HashValue only",
				chunk.IndexQuery{
					TableName: "table",
					HashValue: "flip",
				},
				0,
				[]chunk.IndexEntry{entries[5], entries[6], entries[7]},
			},
			{
				"check RangeValueStart",
				chunk.IndexQuery{
					TableName:       "table",
					HashValue:       "foo",
					RangeValueStart: []byte("bar:2"),
				},
				0,
				[]chunk.IndexEntry{entries[1], entries[2], entries[3], entries[4]},
			},
			{
				"check RangeValuePrefix",
				chunk.IndexQuery{
					TableName:        "table",
					HashValue:        "foo",
					RangeValuePrefix: []byte("baz:"),
				},
				0,
				[]chunk.IndexEntry{entries[3], entries[4]},
			},
			{
				"check ValueEqual",
				chunk.IndexQuery{
					TableName:        "table",
					HashValue:        "foo",
					RangeValuePrefix: []byte("bar"),
					ValueEqual:       []byte("20"),
				},
				0,
				[]chunk.IndexEntry{entries[1]},
			},
			{
				"check retry logic",
				chunk.IndexQuery{
					TableName:        "table",
					HashValue:        "foo",
					RangeValuePrefix: []byte("bar"),
					ValueEqual:       []byte("20"),
				},
				2,
				[]chunk.IndexEntry{entries[1]},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				var have []chunk.IndexEntry
				err = client.QueryPages(context.Background(), tt.query, func(read chunk.ReadBatch, lastPage bool) bool {
					for i := 0; i < read.Len(); i++ {
						have = append(have, chunk.IndexEntry{
							TableName:  tt.query.TableName,
							HashValue:  tt.query.HashValue,
							RangeValue: read.RangeValue(i),
							Value:      read.Value(i),
						})
					}
					return !lastPage
				})
				require.NoError(t, err)
				require.Equal(t, tt.want, have)
			})
		}
	})
}
