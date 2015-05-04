// Copyright 2015 The Chromium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package storage

import (
	"io"

	"github.com/luci/luci-go/common/isolated"
)

type Storage interface {
	// Read data from the storage.
	// Returns an error if the entry is not found.
	Read(isolated.HexDigest) (io.ReadCloser, error)

	// Store data in the storage.
	// The reader expects zlib-compressed data. The data is uncompressed
	// on the fly and compared against the given digest. If it does not
	// match, the store is aborted and an error is returned.
	Write(isolated.HexDigest, io.Reader) error

	// Return true if the store contains an item with the given digest.
	Contains(isolated.HexDigest) bool
}
