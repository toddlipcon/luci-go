// Copyright 2015 The Chromium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

// Package store contains code for storing and retreiving metrics.
package store

import (
	"time"

	"github.com/luci/luci-go/common/tsmon/types"
	"golang.org/x/net/context"
)

// MetricHandle is an opaque handle to a metric registered in the Store.
type MetricHandle interface{}

// Modification contains all the arguments for a Set or Incr call.
type Modification struct {
	Handle    MetricHandle
	ResetTime time.Time
	FieldVals []interface{}
	SetValue  interface{} // nil unless this is a Set.
	IncrDelta interface{} // nil unless this is an Incr.
}

// A Store is responsible for handling all metric data.
type Store interface {
	Register(m types.Metric) (MetricHandle, error)
	Unregister(h MetricHandle)

	Get(ctx context.Context, h MetricHandle, resetTime time.Time, fieldVals []interface{}) (value interface{}, err error)
	Set(ctx context.Context, h MetricHandle, resetTime time.Time, fieldVals []interface{}, value interface{}) error
	Incr(ctx context.Context, h MetricHandle, resetTime time.Time, fieldVals []interface{}, delta interface{}) error
	ModifyMulti(ctx context.Context, mods []Modification) error

	GetAll(ctx context.Context) []types.Cell

	ResetForUnittest()
}