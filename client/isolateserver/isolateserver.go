// Copyright 2015 The Chromium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package isolateserver

import (
	"compress/zlib"
	"crypto/sha1"
	"fmt"
	"hash"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/luci/luci-go/client/internal/common"
)

// IsolateServer is the low-level client interface to interact with an Isolate
// server.
type IsolateServer interface {
	GetHashAlgo() hash.Hash

	ServerCapabilities() (*ServerCapabilities, error)
	// Contains looks up cache pensence on the server of multiple items.
	//
	// The returned list is in the same order than 'items', with entries nil for
	// items that were present.
	Contains(items []*DigestItem) ([]*PushState, error)
	Push(state *PushState, src io.Reader) error
}

// ServerCapabilities is the server details as exposed by the server.
type ServerCapabilities struct {
	ServerVersion string `json:"server_version"`
}

// DigestItem is one item to look up on the server.
type DigestItem struct {
	Digest     HexDigest `json:"digest"`
	IsIsolated bool      `json:"is_isolated"`
	Size       int64     `json:"size"`
}

// PushState per-item state passed from IsolateServer.Contains() to
// IsolateServer.Push().
//
// It's content is implementation specific.
type PushState struct {
	status    preuploadStatus
	uploaded  bool
	finalized bool
}

// New returns a new IsolateServer client.
func New(url, namespace, digestAlgo, compression string) IsolateServer {
	return &isolateServer{
		url: url,
		namespace: namespaceSpec{
			Namespace:   namespace,
			DigestAlgo:  digestAlgo,
			Compression: compression,
		},
	}
}

// Private details.

type namespaceSpec struct {
	Namespace   string `json:"namespace"`
	DigestAlgo  string `json:"digest_hash"`
	Compression string `json:"compression"`
}

// getHashAlgo returns the valid hash.Hash instance for this namespace.
func (n *namespaceSpec) getHashAlgo() (hash.Hash, error) {
	switch n.DigestAlgo {
	case "sha-1":
		return sha1.New(), nil
	default:
		return nil, fmt.Errorf("unknown hash algo \"%s\"", n.DigestAlgo)
	}
}

// getDecompressor returns a valid decompressor for the namespace. It must be
// closed after use.
func (n *namespaceSpec) getDecompressor(in io.Reader) (io.ReadCloser, error) {
	switch n.Compression {
	case "":
		return ioutil.NopCloser(in), nil
	case "flate":
		// The name is a misnomer. It's neither flate, neither gzip, it's zlib with
		// RFC 1950 wrapping.
		return zlib.NewReader(in)
	case "zlib":
		return zlib.NewReader(in)
	default:
		return nil, fmt.Errorf("unknown compression algo \"%s\"", n.Compression)
	}
}

// getCompressor returns a valid compressor for the namespace. It must be
// closed after use.
func (n *namespaceSpec) getCompressor(out io.Writer) (io.WriteCloser, error) {
	switch n.Compression {
	case "":
		return nopWriteCloser{out}, nil
	case "flate":
		// The name is a misnomer. It's neither flate, neither gzip, it's zlib with
		// RFC 1950 wrapping.
		return zlib.NewWriterLevel(out, 7)
	case "zlib":
		return zlib.NewWriterLevel(out, 7)
	default:
		return nil, fmt.Errorf("unknown compression algo \"%s\"", n.Compression)
	}
}

type isolateServer struct {
	url       string
	namespace namespaceSpec
}

type digestCollection struct {
	Items     []*DigestItem `json:"items"`
	Namespace namespaceSpec `json:"namespace"`
}

type preuploadStatus struct {
	GSUploadURL  string `json:"gs_upload_url"`
	UploadTicket string `json:"upload_ticket"`
	Index        Int    `json:"index"`
}

type urlCollection struct {
	Items []preuploadStatus `json:"items"`
}

type finalizeRequest struct {
	UploadTicket string `json:"upload_ticket"`
}

type storageRequest struct {
	UploadTicket string `json:"upload_ticket"`
	Content      []byte `json:"content"`
}

func (i *isolateServer) GetHashAlgo() hash.Hash {
	h, _ := i.namespace.getHashAlgo()
	return h
}

func (i *isolateServer) ServerCapabilities() (*ServerCapabilities, error) {
	url := i.url + "/_ah/api/isolateservice/v1/server_details"
	out := &ServerCapabilities{}
	if _, err := common.PostJSON(nil, url, nil, out); err != nil {
		return nil, err
	}
	return out, nil
}

func (i *isolateServer) Contains(items []*DigestItem) ([]*PushState, error) {
	in := digestCollection{Items: items, Namespace: i.namespace}
	data := &urlCollection{}
	url := i.url + "/_ah/api/isolateservice/v1/preupload"
	if _, err := common.PostJSON(nil, url, in, data); err != nil {
		return nil, err
	}
	out := make([]*PushState, len(items))
	for _, e := range data.Items {
		index := int(e.Index)
		out[index] = &PushState{
			status: e,
		}
	}
	return out, nil
}

func (i *isolateServer) Push(state *PushState, src io.Reader) error {
	// This push operation may be a retry after failed finalization call below,
	// no need to reupload contents in that case.
	if !state.uploaded {
		// PUT file to uploadURL.
		if err := i.doPush(state, src); err != nil {
			return err
		}
		state.uploaded = true
	}

	// Optionally notify the server that it's done.
	if state.status.GSUploadURL != "" {
		// TODO(vadimsh): Calculate MD5 or CRC32C sum while uploading a file and
		// send it to isolated server. That way isolate server can verify that
		// the data safely reached Google Storage (GS provides MD5 and CRC32C of
		// stored files).
		in := finalizeRequest{state.status.UploadTicket}
		url := i.url + "/_ah/api/isolateservice/v1/finalize_gs_upload"
		_, err := common.PostJSON(nil, url, in, nil)
		if err != nil {
			return err
		}
	}
	state.finalized = true
	return nil
}

func (i *isolateServer) doPush(state *PushState, src io.Reader) error {
	reader, writer := io.Pipe()
	defer reader.Close()
	compressor, err := i.namespace.getCompressor(writer)
	if err != nil {
		return err
	}

	go func() {
		io.Copy(compressor, src)
		compressor.Close()
		writer.Close()
	}()

	// DB upload.
	if state.status.GSUploadURL == "" {
		url := i.url + "/_ah/api/isolateservice/v1/store_inline"
		content, err := ioutil.ReadAll(reader)
		if err != nil {
			return err
		}
		in := &storageRequest{state.status.UploadTicket, content}
		_, err = common.PostJSON(nil, url, in, nil)
		return err
	}

	// Upload to GCS.
	client := &http.Client{}
	request, err := http.NewRequest("PUT", state.status.GSUploadURL, reader)
	request.Header.Set("Content-Type", "application/octet-stream")
	// TODO(maruel): For relatively small file, set request.ContentLength so the
	// TCP connection can be reused.
	resp, err := client.Do(request)
	if err == nil {
		io.Copy(ioutil.Discard, resp.Body)
		resp.Body.Close()
	}
	return err
}
