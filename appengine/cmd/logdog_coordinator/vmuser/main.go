// Copyright 2015 The Chromium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/GoogleCloudPlatform/go-endpoints/endpoints"
	"github.com/julienschmidt/httprouter"
	"github.com/luci/luci-go/appengine/ephelper"
	"github.com/luci/luci-go/appengine/ephelper/epfrontend"
	gaeauthServer "github.com/luci/luci-go/appengine/gaeauth/server"
	"github.com/luci/luci-go/appengine/gaemiddleware"
	"github.com/luci/luci-go/appengine/logdog/coordinator/config"
	"github.com/luci/luci-go/appengine/logdog/coordinator/endpoints/admin"
	"github.com/luci/luci-go/appengine/logdog/coordinator/endpoints/logs"
	"github.com/luci/luci-go/appengine/logdog/coordinator/endpoints/service"
	"github.com/luci/luci-go/server/auth"
	"github.com/luci/luci-go/server/middleware"
	"golang.org/x/net/context"
	"google.golang.org/appengine"
)

func authenticator(scopes ...string) auth.Authenticator {
	return auth.Authenticator{
		&gaeauthServer.OAuth2Method{Scopes: scopes},
		gaeauthServer.CookieAuth,
		&gaeauthServer.InboundAppIDAuthMethod{},
	}
}

// base is the root of the middleware chain.
func base(h middleware.Handler) httprouter.Handle {
	a := authenticator(gaeauthServer.EmailScope)
	h = auth.Use(h, a)
	if !appengine.IsDevAppServer() {
		h = middleware.WithPanicCatcher(h)
	}
	return gaemiddleware.BaseProd(h)
}

func endpointMiddleware(c context.Context) (context.Context, error) {
	return config.WithConfig(c), nil
}

func configureEndpoints(h *ephelper.Helper, s *endpoints.Server, sb *ephelper.ServiceBase) error {
	// Admin endpoint.
	if err := h.Register(s, &admin.Admin{ServiceBase: *sb}, &admin.Info, admin.MethodInfoMap); err != nil {
		return fmt.Errorf("failed to register 'admin' endpoint: %v", err)
	}

	// Service endpoint.
	err := h.Register(s, &service.Service{ServiceBase: *sb}, &service.Info, service.MethodInfoMap)
	if err != nil {
		return fmt.Errorf("failed to install 'service' endpoint: %v", err)
	}

	// Logs endpoint.
	if err := h.Register(s, &logs.Logs{ServiceBase: *sb}, &logs.Info, logs.MethodInfoMap); err != nil {
		return fmt.Errorf("failed to register 'logs' endpoint: %v", err)
	}

	return nil
}

// Run installs and executes this site.
func main() {
	router := httprouter.New()

	// Setup Cloud Endpoints.
	ep := endpoints.NewServer("")
	epfe := epfrontend.New("/api/", ep)
	h := ephelper.Helper{
		Frontend: epfe,
	}
	sb := ephelper.ServiceBase{
		Middleware: []ephelper.Middleware{
			ephelper.DefaultMiddleware(nil),
			endpointMiddleware,
		},
	}
	if err := configureEndpoints(&h, ep, &sb); err != nil {
		log.Fatalf("Failed to configure endpoints: %v", err)
	}

	// Standard HTTP endpoints.
	gaeauthServer.InstallHandlers(router, base)

	ep.HandleHTTP(nil)
	epfe.HandleHTTP(nil)
	http.Handle("/", router)
	appengine.Main()
}