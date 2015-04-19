// Copyright 2015 The Chromium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package main

import (
	"errors"
	"fmt"
	"github.com/luci/luci-go/client/isolate"
	"github.com/maruel/subcommands"
)

var cmdCheck = &subcommands.Command{
	UsageLine: "check <options>",
	ShortDesc: "checks that all the inputs are present and generates .isolated",
	LongDesc:  "",
	CommandRun: func() subcommands.CommandRun {
		c := checkRun{}
		c.commonFlags.Init(&c.CommandRunBase)
		c.isolateFlags.Init(&c.CommandRunBase)
		return &c
	},
}

type checkRun struct {
	subcommands.CommandRunBase
	commonFlags
	isolateFlags
}

func (c *checkRun) Parse(a subcommands.Application, args []string) error {
	if err := c.isolateFlags.Parse(); err != nil {
		return err
	}
	if err := c.isolateFlags.RequireIsolateFile(); err != nil {
		return err;
	}
	if err := c.isolateFlags.RequireIsolatedFile(); err != nil {
		return err;
	}
	if len(args) != 0 {
		return errors.New("position arguments not expected")
	}
	return nil
}

func (c *checkRun) main(a subcommands.Application, args []string) error {
	if c.verbose {
		fmt.Printf("Isolate:   %s\n", c.Isolate)
		fmt.Printf("Isolated:  %s\n", c.Isolated)
		fmt.Printf("Blacklist: %s\n", c.Blacklist)
		fmt.Printf("Config:    %s\n", c.ConfigVariables)
		fmt.Printf("Path:      %s\n", c.PathVariables)
		fmt.Printf("Extra:     %s\n", c.ExtraVariables)
	}

	tree := isolate.Tree{
		Cwd: ".",
		Opts: c.ArchiveOptions,
	}

	_, err := isolate.IsolateAndArchive([]isolate.Tree{tree}, "", "")
	return err
}

func (c *checkRun) Run(a subcommands.Application, args []string) int {
	if err := c.Parse(a, args); err != nil {
		fmt.Fprintf(a.GetErr(), "%s: %s\n", a.GetName(), err)
		return 1
	}
	if err := c.main(a, args); err != nil {
		fmt.Fprintf(a.GetErr(), "%s: %s\n", a.GetName(), err)
		return 1
	}
	return 0
}
