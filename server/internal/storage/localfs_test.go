// Copyright 2015 The Chromium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

// TODO: share code with the other test
package storage

import (
	"io/ioutil"
	"testing"
)

func makeTestFs() *localfs {
	tmpDir, err := ioutil.TempDir("", "localfs_test")
	if err != nil {
		panic(err)
	}
	return newLocalFs(tmpDir)
}

// Test for "correct" cases.
func TestFsStorage(t *testing.T) {
	testReadWriteStorage(t, makeTestFs())
}

// Test when the inserted data doesn't match the expected hash.
func TestInvalidHashFS(t *testing.T) {
	testInvalidHash(t, makeTestFs())
}

// Test when the decompressed data is invalid at the beginning of the string
// (e.g. uncompressed data)
func TestInvalidZlibHeaderFS(t *testing.T) {
	testInvalidZlibHeader(t, makeTestFs())
}

// Test when the decompressed data is invalid in the beginning of the string.
// (e.g. truncated or corrupt data)
func TestInvalidZlibDataFS(t *testing.T) {
	testInvalidZlibData(t, makeTestFs())
}
