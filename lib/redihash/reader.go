// Package redihash implements features to read and write messagepack-formatted
// data to Redis hash records.
package redihash

import (
	"context"

	redis "github.com/redis/go-redis/v9"
	"github.com/vmihailenco/msgpack/v5"
)

// A HashReader can read messagepack-formatted data from Redis hash records.
type HashReader[T any] struct {
	rdb       *redis.Client
	recordKey string
}

// NewHashReader returns a HashReader of the specified type.
func NewHashReader[T any](rdb *redis.Client, recordKey string) *HashReader[T] {
	return &HashReader[T]{
		rdb:       rdb,
		recordKey: recordKey,
	}
}

// Read reads a value from a specific key in a Redis hash record and unmarshals it form messagepack.
func (hw *HashReader[T]) Read(ctx context.Context, dataKey string) (*T, error) {

	res, err := hw.rdb.HGet(ctx, hw.recordKey, dataKey).Result()
	if err != nil {
		return nil, err
	}

	var val T
	err = msgpack.Unmarshal([]byte(res), &val)
	if err != nil {
		return nil, err
	}
	return &val, err
}

// ReadAll reads all the keys in a Redis hash record.
func (hw *HashReader[T]) ReadAll(ctx context.Context) (map[string]*T, error) {

	res, err := hw.rdb.HGetAll(ctx, hw.recordKey).Result()
	if err != nil {
		return nil, err
	}

	numOfItems := len(res)
	dataMap := make(map[string]*T, numOfItems)

	for k, v := range res {

		var val T

		err = msgpack.Unmarshal([]byte(v), &val)
		if err != nil {
			return nil, err
		}

		dataMap[k] = &val
	}
	return dataMap, nil
}
