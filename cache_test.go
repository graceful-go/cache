package cache

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type sourceImpl struct{}

func (s *sourceImpl) Get(ctx context.Context, key string) (interface{}, error) {
	if key == "ok" {
		return "ok", nil
	}
	if key == "nil" {
		return nil, nil
	}
	if key == "int32" {
		return int32(1), nil
	}
	if key == "int64" {
		return int64(1), nil
	}
	if key == "string" {
		return "string", nil
	}
	return nil, errors.New("error")
}

func TestCache(t *testing.T) {
	si := &sourceImpl{}
	ce := NewCache(WithSource(si), WithTTL(time.Microsecond))
	ctx := context.Background()

	// get existed key
	rst, err := ce.Get(ctx, "ok")
	assert.Nil(t, err, "must be ok")
	assert.Equal(t, "ok", rst)
	time.Sleep(time.Microsecond * 2)
	ce.Get(ctx, "ok")

	// get invalid key, expected error
	rst, err = ce.Get(ctx, "err")
	assert.NotNil(t, err, "must be fail")
	assert.Nil(t, rst, "must be nil")

	rst, err = ce.Get(ctx, "nil")
	assert.Nil(t, err, "must be ok")
	assert.Nil(t, rst, "must be nil")

	ce = NewCache(WithSource(si), WithTTL(time.Second*100))
	ce.Get(ctx, "ok")
	rst, err = ce.Get(ctx, "ok")
	assert.Nil(t, err, "must be ok")
	assert.Equal(t, "ok", rst)

	rst, err = ce.Get(ctx, "int32")
	assert.Nil(t, err, "must be ok")
	assert.Equal(t, int32(1), rst)

	r := ce.RGet(ctx, "int32")
	assert.Nil(t, r.Error())
	assert.Equal(t, int32(1), r.Int32())
	assert.Equal(t, int32(1), r.Data())

	r = ce.RGet(ctx, "int64")
	assert.Nil(t, r.Error())
	assert.Equal(t, int64(1), r.Int64())
	assert.Equal(t, int64(1), r.Data())

	r = ce.RGet(ctx, "string")
	assert.Nil(t, r.Error())
	assert.Equal(t, "string", r.String())
	assert.Equal(t, "string", r.Data())
	assert.Equal(t, int64(0), r.Int64())

}
