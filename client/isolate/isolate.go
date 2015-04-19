// Copyright 2015 The Chromium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package isolate

import (
	"fmt"
	"path/filepath"
	"io/ioutil"
	"regexp"

	"github.com/luci/luci-go/client/internal/common"
)

// IsolatedGenJSONVersion is used in the batcharchive json format.
//
// TODO(tandrii): Migrate to batch_archive.go.
const IsolatedGenJSONVersion = 1

// ValidVariable is the regexp of valid isolate variable name.
const ValidVariable = "[A-Za-z_][A-Za-z_0-9]*"

var validVariableMatcher = regexp.MustCompile(ValidVariable)

// IsValidVariable returns true if the variable is a valid symbol name.
func IsValidVariable(variable string) bool {
	return validVariableMatcher.MatchString(variable)
}

// Tree to be isolated.
type Tree struct {
	Cwd  string
	Opts ArchiveOptions
}

// ArchiveOptions for achiving trees.
type ArchiveOptions struct {
	Isolate         string            `json:"isolate"`
	Isolated        string            `json:"isolated"`
	Blacklist       common.Strings    `json:"blacklist"`
	PathVariables   common.KeyValVars `json:"path_variables"`
	ExtraVariables  common.KeyValVars `json:"extra_variables"`
	ConfigVariables common.KeyValVars `json:"config_variables"`
}

// Init initializes with non-nil values.
func (a *ArchiveOptions) Init() {
	a.Blacklist = common.Strings{}
	a.PathVariables = common.KeyValVars{}
	a.ExtraVariables = common.KeyValVars{}
	a.ConfigVariables = common.KeyValVars{}
}

func replaceVars(str string, opts ArchiveOptions) string {
	r := regexp.MustCompile("<\\(" + ValidVariable + "?\\)")
	return r.ReplaceAllStringFunc(str, func(match string)string {
		var_name := match[2:len(match) - 1]
		if v, ok := opts.PathVariables[var_name]; ok {
			return v
		}
		if v, ok := opts.ExtraVariables[var_name]; ok {
			return v
		}
		if v, ok := opts.ConfigVariables[var_name]; ok {
			return v
		}
		panic("No value for variable " + var_name)
	})

}

type loadedIsolate struct {
	Dependencies []string
	IsolateDir string
}

func loadIsolate(tree Tree) (*loadedIsolate, error) {
	rel_path := filepath.Join(tree.Cwd, tree.Opts.Isolate)
	content, err := ioutil.ReadFile(rel_path)
	if err != nil {
		return nil, err
	}

	_, deps, _, isolate_dir, err := LoadIsolateForConfig(tree.Cwd, content, tree.Opts.ConfigVariables)
	if err != nil {
		return nil, err
	}

	loaded := new(loadedIsolate)
	loaded.IsolateDir = isolate_dir
	loaded.Dependencies = make([]string, len(deps))
	for i, dep := range deps {
		loaded.Dependencies[i] = replaceVars(dep, tree.Opts)
	}
	return loaded, nil
}

func IsolateAndArchive(trees []Tree, namespace string, server string) (
	map[string]string, error) {

	all_loaded := []*loadedIsolate{}
	for _, tree := range trees {
		loaded, err := loadIsolate(tree)
		if err != nil {
			return nil, err
		}
		all_loaded = append(all_loaded, loaded)
	}

	info_loader := LoadOrCreateCache()
	defer info_loader.Save()

	for _, loaded := range all_loaded {
		fmt.Printf("%+v\n", loaded)

		for _, dep := range loaded.Dependencies {
			infos, err := info_loader.LookupRecursive(filepath.Join(loaded.IsolateDir, dep))
			if err != nil {
				return nil, err
			}
			for _, info := range infos {
				fmt.Printf("%s\t%+v\n", dep, info)
			}
		}
	}


	return nil, nil
}
