// Copyright 2015 The Chromium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package storage

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"sync"

	"github.com/luci/luci-go/common/isolated"
)

type memoryStorage struct {
	lock     sync.Mutex
	contents map[isolated.HexDigest][]byte
}

func NewMemory() Storage {
	return newMemory()
}

func newMemory() *memoryStorage {
	return &memoryStorage{
		contents: map[isolated.HexDigest][]byte{},
	}
}

func (m *memoryStorage) Read(digest isolated.HexDigest) (io.ReadCloser, error) {
	m.lock.Lock()
	defer m.lock.Unlock()
	data, ok := m.contents[digest]
	if !ok {
		return nil, errors.New("not found")
	}
	return ioutil.NopCloser(bytes.NewReader(data)), nil
}

func (m *memoryStorage) Contains(digest isolated.HexDigest) bool {
	m.lock.Lock()
	defer m.lock.Unlock()
	_, ok := m.contents[digest]
	return ok
}

func (m *memoryStorage) Write(digest isolated.HexDigest, r io.Reader) error {
	// r -> buf
	//   -> pipe -> decompress -> hasher
	pipe_r, pipe_w := io.Pipe()
	tee := io.TeeReader(r, pipe_w)
	hasher := isolated.GetHash()
	buf := new(bytes.Buffer)
	err_ch := make(chan error)
	go func() {
		// If we hit an error before reading all of the data,
		// discard the rest. Otherwise the other side of the pipe
		// will block.
		defer io.Copy(ioutil.Discard, pipe_r)

		decomp, err := isolated.GetDecompressor(pipe_r)
		if err == nil {
			_, err = io.Copy(hasher, decomp)
		}
		if err != nil {
			err_ch <- err
			return
		}
		err_ch <- nil
	}()

	// Pulling data from the tee causes it to write to pipe_w.
	// The other goroutine pulls from pipe_r and writes into
	// hasher.
	if _, err := io.Copy(buf, tee); err != nil {
		return fmt.Errorf("failed to read data: %s", err.Error())
	}
	if err := <-err_ch; err != nil {
		return fmt.Errorf("failed to calculate digest of decompressed data: %s",
			err.Error())
	}

	calculated := isolated.HexDigest(hex.EncodeToString(hasher.Sum(nil)))
	if calculated != digest {
		return fmt.Errorf("calculated incorrect hash %s (expected %s)",
			calculated,
			digest)
	}

	m.lock.Lock()
	defer m.lock.Unlock()
	m.contents[digest] = buf.Bytes()
	return nil
}

func (m *memoryStorage) numItems() int {
	m.lock.Lock()
	defer m.lock.Unlock()
	return len(m.contents)
}
