package chunk

import "context"

// IndexClient is a client for the storage of the index (e.g. DynamoDB or Bigtable).
type IndexClient interface {
	Stop()

	// For the write path.
	NewWriteBatch() WriteBatch
	BatchWrite(context.Context, WriteBatch) error

	// For the read path.
	QueryPages(ctx context.Context, queries []IndexQuery, callback func(IndexQuery, ReadBatch) (shouldContinue bool)) error
}

// ObjectClient is for storing and retrieving chunks.
type ObjectClient interface {
	Stop()

	PutChunks(ctx context.Context, chunks []Chunk) error
	GetChunks(ctx context.Context, chunks []Chunk) ([]Chunk, error)

	NewStreamer() Streamer
}

// WriteBatch represents a batch of writes.
type WriteBatch interface {
	Add(tableName, hashValue string, rangeValue []byte, value []byte)
}

// ReadBatch represents the results of a QueryPages.
type ReadBatch interface {
	Iterator() ReadBatchIterator
}

// ReadBatchIterator is an iterator over a ReadBatch.
type ReadBatchIterator interface {
	Next() bool
	RangeValue() []byte
	Value() []byte
}

// Streamer represents the configuration for streaming chunks
type Streamer interface {
	Add(table string, users string, from int, to int)
	Stream(ctx context.Context, forwardChan chan []Chunk) error
	Size(ctx context.Context) (int, error)
}
