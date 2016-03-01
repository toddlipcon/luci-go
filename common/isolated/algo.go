// Copyright 2015 The Chromium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

// Defines the algorithms used by the isolated server.

package isolated

import (
	"compress/zlib"
	"crypto/sha1"
	"hash"
	"io"
	// https://crbug.com/552697
	//"github.com/klauspost/compress/zlib"
)

var hashLength = sha1.New().Size()

// GetHash returns a fresh instance of the hashing algorithm to be used to
// calculate the HexDigest.
//
// It is currently hardcoded to sha-1.
func GetHash() hash.Hash {
	return sha1.New()
}

// GetDecompressor returns a fresh instance of the decompression algorithm.
//
// It must be closed after use.
//
// It is currently hardcoded to RFC 1950 (zlib).
func GetDecompressor(in io.Reader) (io.ReadCloser, error) {
	d, err := zlib.NewReader(in)
	if err != nil {
		return nil, err
	}
	return d, nil
}

// GetCompressor returns a fresh instance of the compression algorithm.
//
// It must be closed after use.
//
// It is currently hardcoded to RFC 1950 (zlib).
func GetCompressor(out io.Writer) io.WriteCloser {
	c, _ := zlib.NewWriterLevel(out, 7)
	return c
}

// HexDigest is the hash of a file that is hex-encoded. Only lower case letters
// are accepted.
type HexDigest string

// Validate returns true if the hash is valid.
func (d HexDigest) Validate() bool {
	if len(d) != hashLength*2 {
		return false
	}
	for _, c := range d {
		if ('0' <= c && c <= '9') || ('a' <= c && c <= 'f') {
			continue
		}
		return false
	}
	return true
}

// HexDigests is a slice of HexDigest that implements sort.Interface.
type HexDigests []HexDigest

func (h HexDigests) Len() int           { return len(h) }
func (h HexDigests) Less(i, j int) bool { return h[i] < h[j] }
func (h HexDigests) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
