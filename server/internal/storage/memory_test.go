// Copyright 2015 The Chromium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package storage

import (
	"testing"
)

// Test for "correct" cases.
func TestMemoryStorage(t *testing.T) {
	testReadWriteStorage(t, newMemory())
}

// Test when the inserted data doesn't match the expected hash.
func TestInvalidHash(t *testing.T) {
	testInvalidHash(t, newMemory())
}

// Test when the decompressed data is invalid at the beginning of the string
// (e.g. uncompressed data)
func TestInvalidZlibHeader(t *testing.T) {
	testInvalidZlibHeader(t, newMemory())
}

// Test when the decompressed data is invalid in the beginning of the string.
// (e.g. truncated or corrupt data)
func TestInvalidZlibData(t *testing.T) {
	testInvalidZlibData(t, newMemory())
}