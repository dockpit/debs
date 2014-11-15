// Copyright 2012 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package debs

import (
	"log"
	"os"
	"strings"
	"sync"
)

//@todo configure this:
var buildV = true
var buildN = true
var buildX = true

var exitStatus = 0
var exitMu sync.Mutex

// mergeEnvLists merges the two environment lists such that
// variables with the same name in "in" replace those in "out".
func mergeEnvLists(in, out []string) []string {
NextVar:
	for _, inkv := range in {
		k := strings.SplitAfterN(inkv, "=", 2)[0]
		for i, outkv := range out {
			if strings.HasPrefix(outkv, k) {
				out[i] = inkv
				continue NextVar
			}
		}
		out = append(out, inkv)
	}
	return out
}

// envForDir returns a copy of the environment
// suitable for running in the given directory.
// The environment is the current process's environment
// but with an updated $PWD, so that an os.Getwd in the
// child will be faster.
func envForDir(dir string) []string {
	env := os.Environ()
	// Internally we only use rooted paths, so dir is rooted.
	// Even if dir is not rooted, no harm done.
	return mergeEnvLists([]string{"PWD=" + dir}, env)
}

func setExitStatus(n int) {
	exitMu.Lock()
	if exitStatus < n {
		exitStatus = n
	}
	exitMu.Unlock()
}

var atexitFuncs []func()

func exit() {
	for _, f := range atexitFuncs {
		f()
	}
	os.Exit(exitStatus)
}

func exitIfErrors() {
	if exitStatus != 0 {
		exit()
	}
}

func errorf(format string, args ...interface{}) {
	log.Printf(format, args...)
	setExitStatus(1)
}
