package redihash

import (
	"context"

	redis "github.com/redis/go-redis/v9"
	"github.com/vmihailenco/msgpack/v5"
)

// A HashWriter can write messagepack-formatted data to Redis hash records.
type HashWriter[T any] struct {
	rdb       *redis.Client
	recordKey string
}

// NewHashWriter returns a HashWriter of the specified type.
func NewHashWriter[T any](rdb *redis.Client, recordKey string) *HashWriter[T] {
	return &HashWriter[T]{
		rdb:       rdb,
		recordKey: recordKey,
	}
}

// Write marshals a value to messagepack and writes it to a specific key in a Redis hash record.
func (hw *HashWriter[T]) Write(ctx context.Context, dataKey string, val *T) error {

	encodedData, err := msgpack.Marshal(val)
	if err != nil {
		return err
	}

	_, err = hw.rdb.HSet(ctx, hw.recordKey, dataKey, encodedData).Result()
	return err
}

// Del deletes  specific key from a Redis hash record.
func (hw *HashWriter[T]) Del(ctx context.Context, dataKey string) error {

	_, err := hw.rdb.HDel(ctx, hw.recordKey, dataKey).Result()
	if err != nil {
		return err
	}
	return nil
}
