// Copyright 2015 The Chromium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package lhttp

import (
	"errors"
	"net/url"
)

// URLToHTTPS ensures the url is https://.
func URLToHTTPS(s string) (string, error) {
	u, err := url.Parse(s)
	if err != nil {
		return "", err
	}
	if u.Scheme != "" && u.Scheme != "https" && u.Scheme != "http" {
		return "", errors.New("Only http:// or https:// scheme is accepted.")
	}
	if u.Scheme == "" {
		s = "https://" + s
	}
	if _, err = url.Parse(s); err != nil {
		return "", err
	}
	return s, nil
}
