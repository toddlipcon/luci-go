// Copyright 2015 The Chromium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package storage

import (
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/luci/luci-go/common/isolated"
)

type localfs struct {
	rootDir string
}

func NewLocalFs(rootDir string) Storage {
	return newLocalFs(rootDir)
}

func newLocalFs(rootDir string) *localfs {
	return &localfs{
		rootDir: rootDir,
	}
}

func (fs *localfs) path(digest isolated.HexDigest) string {
	s := string(digest)
	return filepath.Join(fs.rootDir, s[0:2], s[2:4], s)
}

func (fs *localfs) Read(digest isolated.HexDigest) (io.ReadCloser, error) {
	if !digest.Validate() {
		return nil, errors.New("invalid digest")
	}
	p := fs.path(digest)
	return os.Open(p)
}

func (fs *localfs) Contains(digest isolated.HexDigest) bool {
	if !digest.Validate() {
		return false
	}

	if _, err := os.Stat(fs.path(digest)); err != nil {
		return false
	}
	return true
}

func (fs *localfs) Write(digest isolated.HexDigest, r io.Reader) error {
	if !digest.Validate() {
		return errors.New("invalid digest")
	}
	path := fs.path(digest)
	dir_path := filepath.Dir(path)
	if err := os.MkdirAll(dir_path, 0777); err != nil {
		return err
	}
	file_w, err := ioutil.TempFile(dir_path, filepath.Base(path)+".tmp")
	if err != nil {
		return err
	}
	tmp_path := file_w.Name()
	success := false
	defer func() {
		if !success {
			_ = os.Remove(tmp_path)
		}
		file_w.Close()
	}()

	// r -> file_w
	//   -> pipe -> decompress -> hasher
	pipe_r, pipe_w := io.Pipe()
	defer pipe_r.Close()
	tee := io.TeeReader(r, pipe_w)
	hasher := isolated.GetHash()
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
	if _, err := io.Copy(file_w, tee); err != nil {
		return fmt.Errorf("failed to copy data: %s", err.Error())
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

	// We fully wrote the tmp file, and the digest matched. Move it into place
	if err := os.Rename(tmp_path, path); err != nil {
		return fmt.Errorf("failed to commit uploaded file: %s", err.Error())
	}

	return nil
}
