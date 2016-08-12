// Copyright 2015 The Chromium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.
package main

// TODO(todd): probably need to rename this directory so the binary
// isn't just called 'standalone'
// TODO(todd): remove ugly debug printfs and add real logging, add comments, etc.
// TODO(todd): fix error handling - the server shouldn't have a global error
// TODO(todd): add command line flags or some other configuration
// TODO(todd): add TLS support
// TODO(todd): factor out oauth stubs into a separate auth module, etc.

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"regexp"
	"strings"
	"sync"

	"github.com/fvbock/endless"
	"github.com/gorilla/handlers"
	"github.com/luci/luci-go/common/isolated"
	"github.com/luci/luci-go/server/internal/storage"
)

type saServer struct {
	mux     *http.ServeMux
	lock    sync.Mutex
	err     error
	storage storage.Storage
	log     io.Writer
}

type jsonAPI func(body io.Reader) (interface{}, error)

// handlerJSON converts a jsonAPI http handler to a proper http.Handler.
func handlerJSON(handler jsonAPI) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("Request for " + r.URL.String())
		contentType := "application/json; charset=utf-8"
		if r.Header.Get("Content-Type") != contentType {
			log.Printf("invalid content type: %s", r.Header.Get("Content-Type"))
			return
		}
		defer r.Body.Close()
		out, err := handler(r.Body)
		if err != nil {
			writeError(w, err)
			return
		}
		if err := writeJsonResponse(w, out); err != nil {
			writeError(w, err)
			return
		}
	})
}

func writeJsonResponse(w http.ResponseWriter, jsonData interface{}) error {
	contentType := "application/json; charset=utf-8"
	w.Header().Set("Content-Type", contentType)
	j := json.NewEncoder(w)
	if err := j.Encode(jsonData); err != nil {
		return err
	}
	return nil
}

func (s *saServer) handle(path string, handler http.Handler) {
	s.mux.Handle(path, handlers.LoggingHandler(s.log, handler))
}

func (s *saServer) handleFunc(path string, f func(http.ResponseWriter, *http.Request)) {
	s.handle(path, http.HandlerFunc(f))
}

func (s *saServer) handleJSON(path string, handler jsonAPI) {
	s.handle(path, handlerJSON(handler))
}

func (server *saServer) serverDetails(body io.Reader) (interface{}, error) {
	content, err := ioutil.ReadAll(body)
	if err != nil {
		return nil, err
	}
	if string(content) != "{}" {
		return nil, fmt.Errorf("unexpected content %#v", string(content))
	}
	return map[string]string{"server_version": "v1"}, nil
}

func (server *saServer) preupload(body io.Reader) (interface{}, error) {
	data := &isolated.DigestCollection{}
	if err := json.NewDecoder(body).Decode(data); err != nil {
		return nil, err
	}
	if data.Namespace.Namespace != "default-gzip" {
		return nil, fmt.Errorf("unexpected namespace %#v", data.Namespace.Namespace)
	}
	out := &isolated.URLCollection{}

	server.lock.Lock()
	defer server.lock.Unlock()
	for i, d := range data.Items {
		if !server.storage.Contains(d.Digest) {
			ticket := "ticket:" + string(d.Digest)
			out.Items = append(out.Items, isolated.PreuploadStatus{"", ticket, isolated.Int(i)})
		}
	}
	return out, nil
}

func (server *saServer) finalizeGSUpload(body io.Reader) (interface{}, error) {
	data := &isolated.FinalizeRequest{}
	if err := json.NewDecoder(body).Decode(data); err != nil {
		return nil, err
	}

	server.lock.Lock()
	defer server.lock.Unlock()
	return map[string]string{"ok": "true"}, nil
}

func (server *saServer) storeInline(body io.Reader) (interface{}, error) {
	data := &isolated.StorageRequest{}
	if err := json.NewDecoder(body).Decode(data); err != nil {
		return nil, err
	}

	prefix := "ticket:"
	if !strings.HasPrefix(data.UploadTicket, prefix) {
		return nil, fmt.Errorf("unexpected ticket %#v", data.UploadTicket)
	}

	digest := isolated.HexDigest(data.UploadTicket[len(prefix):])
	if !digest.Validate() {
		return nil, fmt.Errorf("invalid digest %#v", digest)
	}

	err := server.storage.Write(digest, bytes.NewBuffer(data.Content))
	if err != nil {
		return nil, err
	}

	return map[string]string{"ok": "true"}, nil
}

func (server *saServer) accountsSelf(w http.ResponseWriter, r *http.Request) {
	ret := map[string]string{"identity": "anonymous:todd@lipcon.org"}
	if err := writeJsonResponse(w, ret); err != nil {
		writeError(w, err)
		return
	}
}

func (server *saServer) oauthConfig(w http.ResponseWriter, r *http.Request) {
	ret := map[string]string{"additional_client_ids": "",
		"client_id":            "x",
		"client_not_so_secret": "y",
		"primary_url":          "http://localhost:4242/",
	}

	if err := writeJsonResponse(w, ret); err != nil {
		panic(err)
	}
}

func (server *saServer) retrieveContent(w http.ResponseWriter, r *http.Request) {
	re := regexp.MustCompile("^/content-gs/retrieve/default-gzip/(.+)$")
	hash := re.FindStringSubmatch(r.URL.Path)[1]
	data_reader, err := server.storage.Read(isolated.HexDigest(hash))
	defer data_reader.Close()
	if err != nil {
		writeError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/octet-stream")
	io.Copy(w, data_reader)
}

func writeError(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(err.Error()))
	log.Println(err.Error())
	return
}

func New() *saServer {
	s := &saServer{
		mux:     http.NewServeMux(),
		storage: storage.NewLocalFs("/tmp/storage"),
		log:     os.Stderr,
	}
	s.handleJSON("/_ah/api/isolateservice/v1/server_details", s.serverDetails)
	s.handleJSON("/_ah/api/isolateservice/v1/preupload", s.preupload)
	s.handleJSON("/_ah/api/isolateservice/v1/finalize_gs_upload", s.finalizeGSUpload)
	s.handleJSON("/_ah/api/isolateservice/v1/store_inline", s.storeInline)
	s.handleFunc("/auth/api/v1/accounts/self", s.accountsSelf)
	s.handleFunc("/auth/api/v1/server/oauth_config", s.oauthConfig)
	s.handleFunc("/content-gs/retrieve/default-gzip/", s.retrieveContent)
	// Fail on anything else.
	s.handleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		log.Printf("Request for unknown endpoint %s\n", req.URL)
	})
	return s
}

func StartPProfServer() {
	go func() {
		log.Println(http.ListenAndServe("0.0.0.0:4243", nil))
	}()
}

func main() {
	StartPProfServer()
	server := New()
	err := endless.ListenAndServe("0.0.0.0:4242", server.mux)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
	}
	log.Println("Server on 4242 stopped")

	os.Exit(0)
}
