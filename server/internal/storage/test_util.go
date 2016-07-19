// Copyright 2015 The Chromium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package storage

import (
	"bytes"
	"io"
	"io/ioutil"
	"testing"

	"github.com/luci/luci-go/common/isolated"
	"github.com/maruel/ut"
)

func makeBuf(data []byte) io.ReadCloser {
	return ioutil.NopCloser(bytes.NewBuffer(data))
}

func testReadWriteStorage(t *testing.T, s Storage) {
	// Prepare some test data.
	data := []byte("hello world")
	hash := isolated.HashBytes(data)
	data_z := isolated.CompressBytes(data)

	// Store should initially not contain the data.
	ut.AssertEqual(t, false, s.Contains(hash))

	// Write some data into the store.
	err := s.Write(hash, makeBuf(data_z))
	ut.AssertEqual(t, nil, err)

	// Read it back out.
	ut.AssertEqual(t, true, s.Contains(hash))
	r, err := s.Read(hash)
	ut.AssertEqual(t, nil, err)
	data_returned, _ := ioutil.ReadAll(r)
	ut.AssertEqual(t, data_z, data_returned)
}

func testInvalidHash(t *testing.T, s Storage) {
	// Prepare some test data.
	data := []byte("hello world")
	hash := isolated.HashBytes([]byte("other data"))
	data_z := isolated.CompressBytes(data)

	// Write some data into the store.
	err := s.Write(hash, makeBuf(data_z))
	ut.AssertEqual(t, true, err != nil)
	ut.AssertEqual(t, "calculated incorrect hash "+
		"2aae6c35c94fcfb415dbe95f408b9ce91ee846ed"+
		" (expected ddd9d41363a535aeb9a8178ed03ede5ca69fd438)",
		err.Error())

	ut.AssertEqual(t, false, s.Contains(hash))
}

func testInvalidZlibHeader(t *testing.T, s Storage) {
	// Prepare some test data.
	data := []byte("this is not a zlib string")
	hash := isolated.HashBytes(data)

	err := s.Write(hash, makeBuf(data))
	ut.AssertEqual(t, true, err != nil)
	ut.AssertEqual(t, "failed to calculate digest of decompressed data: "+
		"zlib: invalid header", err.Error())
}

func testInvalidZlibData(t *testing.T, s Storage) {
	m := NewMemory()
	// Prepare some test data.
	data := []byte("this is not a zlib string")
	data_z := isolated.CompressBytes(data)
	data_z[35] = 'x'
	hash := isolated.HashBytes(data_z)

	err := m.Write(hash, makeBuf(data_z))
	ut.AssertEqual(t, true, err != nil)
	ut.AssertEqual(t, "failed to calculate digest of decompressed data: "+
		"zlib: invalid checksum", err.Error())
}