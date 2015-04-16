// Copyright 2015 The Chromium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package isolateserver

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	"github.com/maruel/ut"
)

type jsonAPI func(body io.Reader) interface{}

// handlerJSON converts a jsonAPI http handler to a proper http.Handler.
func handlerJSON(t *testing.T, handler jsonAPI) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		contentType := "application/json; charset=utf-8"
		if r.Header.Get("Content-Type") != contentType {
			t.Fatalf("invalid content type: %s", r.Header.Get("Content-Type"))
		}
		defer r.Body.Close()
		out := handler(r.Body)
		w.Header().Set("Content-Type", contentType)
		j := json.NewEncoder(w)
		ut.AssertEqual(t, nil, j.Encode(out))
	})
}

func handleJSON(t *testing.T, mux *http.ServeMux, path string, handler jsonAPI) {
	mux.Handle(path, handlerJSON(t, handler))
}

type isolateServerFake struct {
	lock     sync.Mutex
	contents map[HexDigest][]byte
}

func newIsolateServerFake(t *testing.T) (http.Handler, *isolateServerFake) {
	mux := http.NewServeMux()
	server := &isolateServerFake{
		contents: map[HexDigest][]byte{},
	}

	handleJSON(t, mux, "/_ah/api/isolateservice/v1/server_details", func(body io.Reader) interface{} {
		content, err := ioutil.ReadAll(body)
		ut.AssertEqual(t, nil, err)
		ut.AssertEqual(t, []byte("{}"), content)
		return &ServerCapabilities{"v1"}
	})

	handleJSON(t, mux, "/_ah/api/isolateservice/v1/preupload", func(body io.Reader) interface{} {
		data := &digestCollection{}
		ut.AssertEqual(t, nil, json.NewDecoder(body).Decode(data))
		ut.AssertEqual(t, "default", data.Namespace.Namespace)
		ut.AssertEqual(t, "flate", data.Namespace.Compression)
		ut.AssertEqual(t, "sha-1", data.Namespace.DigestAlgo)
		out := &urlCollection{}

		server.lock.Lock()
		defer server.lock.Unlock()
		for i, d := range data.Items {
			if _, ok := server.contents[d.Digest]; !ok {
				ticket := "ticket:" + string(d.Digest)
				out.Items = append(out.Items, preuploadStatus{"", ticket, Int(i)})
			}
		}
		return out
	})

	handleJSON(t, mux, "/_ah/api/isolateservice/v1/finalize_gs_upload", func(body io.Reader) interface{} {
		data := &finalizeRequest{}
		ut.AssertEqual(t, nil, json.NewDecoder(body).Decode(data))

		server.lock.Lock()
		defer server.lock.Unlock()
		return map[string]string{"ok": "true"}
	})

	handleJSON(t, mux, "/_ah/api/isolateservice/v1/store_inline", func(body io.Reader) interface{} {
		data := &storageRequest{}
		ut.AssertEqual(t, nil, json.NewDecoder(body).Decode(data))
		prefix := "ticket:"
		ut.AssertEqual(t, true, strings.HasPrefix(data.UploadTicket, prefix))
		digest := HexDigest(data.UploadTicket[len(prefix):])

		// TODO(maruel): This information is constructed from the ticket.
		n := namespaceSpec{DigestAlgo: "sha-1", Compression: "flate"}

		algo, err := n.getHashAlgo()
		ut.AssertEqual(t, nil, err)
		ut.AssertEqual(t, true, digest.Validate(algo))
		comp, err := n.getDecompressor(bytes.NewBuffer(data.Content))
		ut.AssertEqual(t, nil, err)
		raw, err := ioutil.ReadAll(comp)
		ut.AssertEqual(t, nil, err)
		ut.AssertEqual(t, digest, Hash(algo, raw))

		server.lock.Lock()
		defer server.lock.Unlock()
		server.contents[digest] = raw
		return map[string]string{"ok": "true"}
	})

	// Fail on anything else.
	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		t.Fatal()
	})
	return mux, server
}

func TestIsolateServerCaps(t *testing.T) {
	t.Parallel()
	mux, _ := newIsolateServerFake(t)
	ts := httptest.NewServer(mux)
	defer ts.Close()
	client := New(ts.URL, "default", "sha-1", "flate")
	caps, err := client.ServerCapabilities()
	ut.AssertEqual(t, nil, err)
	ut.AssertEqual(t, &ServerCapabilities{"v1"}, caps)
}

type items struct {
	digests  []*DigestItem
	contents [][]byte
}

func makeItems(contents ...string) items {
	out := items{}
	for _, content := range contents {
		// TODO(maruel): Remove hardcoding.
		h := sha1.New()
		c := []byte(content)
		h.Write(c)
		hex := HexDigest(hex.EncodeToString(h.Sum(nil)))
		out.digests = append(out.digests, &DigestItem{hex, false, int64(len(content))})
		out.contents = append(out.contents, c)
	}
	return out
}

func TestIsolateServer(t *testing.T) {
	t.Parallel()
	mux, server := newIsolateServerFake(t)
	ts := httptest.NewServer(mux)
	defer ts.Close()
	client := New(ts.URL, "default", "sha-1", "flate")

	files := makeItems("foo", "bar")
	states, err := client.Contains(files.digests)
	ut.AssertEqual(t, nil, err)
	ut.AssertEqual(t, len(files.digests), len(states))
	for index, state := range states {
		err = client.Push(state, bytes.NewBuffer(files.contents[index]))
		ut.AssertEqual(t, nil, err)
	}
	// foo and bar.
	expected := map[HexDigest][]byte{
		"0beec7b5ea3f0fdbc95d0dd47f3c5bc275da8a33": {0x66, 0x6f, 0x6f},
		"62cdb7020ff920e5aa642c3d4066950dd1f01f4d": {0x62, 0x61, 0x72},
	}
	ut.AssertEqual(t, expected, server.contents)
	states, err = client.Contains(files.digests)
	ut.AssertEqual(t, nil, err)
	ut.AssertEqual(t, len(files.digests), len(states))
	for _, state := range states {
		ut.AssertEqual(t, (*PushState)(nil), state)
	}
}
