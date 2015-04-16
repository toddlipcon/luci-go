// Copyright 2015 The Chromium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package archiver

import (
	"io"
	"os"
	"sync"
	"time"

	"github.com/luci/luci-go/client/internal/common"
	"github.com/luci/luci-go/client/isolateserver"
	"github.com/maruel/interrupt"
)

// Archiver is an high level interface to an isolateserver.IsolateServer.
type Archiver interface {
	io.Closer
	PushFile(path string)
	Flush()
	Stats() *Stats
	HasFailed() bool
}

type UploadStat struct {
	Duration time.Duration
	Size     int64
}

// Statistics from the Archiver.
type Stats struct {
	Hits   []int64
	Misses []int64
	Pushed []*UploadStat
}

type archiver struct {
	is                    isolateserver.IsolateServer
	maxConcurrentHash     int
	maxConcurrentContains int
	maxConcurrentUpload   int
	containsBatchingDelay time.Duration
	filesToHash           chan string
	itemsToLookup         chan *hashedItem
	itemsToUpload         chan *lookedUpItem
	wg                    sync.WaitGroup

	lock  sync.Mutex
	err   error
	stats Stats
}

// New returns a thread-safe Archiver instance.
//
// TODO(maruel): Cache hashes and server cache presence.
func New(is isolateserver.IsolateServer) Archiver {
	a := &archiver{
		is:                    is,
		maxConcurrentHash:     5,
		maxConcurrentContains: 16,
		maxConcurrentUpload:   8,
		containsBatchingDelay: 100 * time.Millisecond,
		filesToHash:           make(chan string, 10240),
		itemsToLookup:         make(chan *hashedItem, 10240),
		itemsToUpload:         make(chan *lookedUpItem, 10240),
	}
	a.wg.Add(1)
	go func() {
		defer a.wg.Done()
		a.hashLoop()
	}()
	a.wg.Add(1)
	go func() {
		defer a.wg.Done()
		a.containsLoop()
	}()
	a.wg.Add(1)
	go func() {
		defer a.wg.Done()
		a.uploadLoop()
	}()
	return a
}

func (a *archiver) Close() error {
	close(a.filesToHash)
	a.Flush()
	a.wg.Wait()
	return nil
}

func (a *archiver) PushFile(path string) {
	a.filesToHash <- path
}

func (a *archiver) Flush() {
}

func (a *archiver) Stats() *Stats {
	s := &Stats{}
	a.lock.Lock()
	defer a.lock.Unlock()
	*s = a.stats
	s.Pushed = make([]*UploadStat, len(a.stats.Pushed))
	copy(s.Pushed, a.stats.Pushed)
	return s
}

func (a *archiver) HasFailed() bool {
	a.lock.Lock()
	defer a.lock.Unlock()
	return a.err == nil
}

// Private details.

type hashedItem struct {
	digest isolateserver.DigestItem
	path   string
}

type lookedUpItem struct {
	state *isolateserver.PushState
	path  string
	size  int64
}

func (a *archiver) hashLoop() {
	defer close(a.itemsToLookup)
	var wg sync.WaitGroup
	s := common.NewSemaphore(a.maxConcurrentHash)

	for file := range a.filesToHash {
		wg.Add(1)
		go func(f string) {
			s.Wait()
			defer func() {
				wg.Done()
				s.Signal()
			}()
			d, e := isolateserver.HashFile(a.is.GetHashAlgo(), f)
			if e != nil {
				//log.Printf("%s: %s\n", f, e)
				//err = e
			} else {
				//log.Printf("%s: %s\n", f, d.Digest)
				a.itemsToLookup <- &hashedItem{d, f}
			}
		}(file)
	}
	wg.Wait()
}

func (a *archiver) containsLoop() {
	defer close(a.itemsToUpload)
	items := []*hashedItem{}
	never := make(<-chan time.Time)
	timer := never
	loop := true
	for loop && !interrupt.IsSet() {
		select {
		case <-timer:
			tmp := make([]*isolateserver.DigestItem, len(items))
			for i, e := range items {
				tmp[i] = &e.digest
			}
			states, err := a.is.Contains(tmp)
			if err != nil {
				return
			}
			for index, state := range states {
				size := items[index].digest.Size
				if state == nil {
					a.lock.Lock()
					a.stats.Hits = append(a.stats.Hits, size)
					a.lock.Unlock()
				} else {
					a.lock.Lock()
					a.stats.Misses = append(a.stats.Misses, size)
					a.lock.Unlock()
					a.itemsToUpload <- &lookedUpItem{state, items[index].path, size}
				}
			}

			// Reset.
			items = []*hashedItem{}
			timer = never

		case item, ok := <-a.itemsToLookup:
			if !ok {
				loop = false
				break
			}
			items = append(items, item)
			if timer == never {
				timer = time.After(a.containsBatchingDelay)
			}
		}
	}

	if len(items) != 0 {
		tmp := make([]*isolateserver.DigestItem, len(items))
		for i, e := range items {
			tmp[i] = &e.digest
		}
		states, err := a.is.Contains(tmp)
		if err != nil {
			return
		}
		for index, state := range states {
			size := items[index].digest.Size
			if state == nil {
				a.lock.Lock()
				a.stats.Hits = append(a.stats.Hits, size)
				a.lock.Unlock()
			} else {
				a.lock.Lock()
				a.stats.Misses = append(a.stats.Misses, size)
				a.lock.Unlock()
				a.itemsToUpload <- &lookedUpItem{state, items[index].path, size}
			}
		}
	}
}

func (a *archiver) uploadLoop() {
	var wg sync.WaitGroup
	s := common.NewSemaphore(a.maxConcurrentUpload)

	for state := range a.itemsToUpload {
		wg.Add(1)
		go func(h *lookedUpItem) {
			s.Wait()
			defer func() {
				wg.Done()
				s.Signal()
			}()
			_ = a.push(h)
		}(state)
	}
	wg.Wait()
}

func (a *archiver) push(l *lookedUpItem) error {
	src, err := os.Open(l.path)
	if err != nil {
		return err
	}
	defer src.Close()
	start := time.Now()
	err = a.is.Push(l.state, src)
	duration := time.Now().Sub(start)
	u := &UploadStat{duration, l.size}
	a.lock.Lock()
	a.stats.Pushed = append(a.stats.Pushed, u)
	a.lock.Unlock()
	return err
}
