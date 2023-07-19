package redihash

import (
	"context"

	redis "github.com/redis/go-redis/v9"
)

// HashReadWriter reads and writes messagepack-formatted data to specific Redis
// hash records.
type HashReadWriter[T any] struct {
	*HashWriter[T]
	*HashReader[T]
}

// NewHashReadWriter returns a HashReadWriter of the specified type.
func NewHashReadWriter[T any](rdb *redis.Client, recordKey string) *HashReadWriter[T] {
	return &HashReadWriter[T]{
		HashWriter: NewHashWriter[T](rdb, recordKey),
		HashReader: NewHashReader[T](rdb, recordKey),
	}
}

// Update modifies a specific key in a Reds hash record by executing the
// 'updater' function from its parameters.
func (hrw *HashReadWriter[T]) Update(ctx context.Context, dataKey string, updater func(*T)) error {

	val, err := hrw.Read(ctx, dataKey)
	if err != nil {
		return err
	}

	updater(val)

	err = hrw.Write(ctx, dataKey, val)

	return err
}
