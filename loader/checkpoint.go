// Copyright 2016 PingCAP, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"os"
	"sync"

	"bufio"
	"fmt"
	"github.com/juju/errors"
	"github.com/ngaut/log"
	"io"
	"strings"
)

type CheckPoint struct {
	sync.RWMutex
	path              string
	restoredFiles     map[string]struct{}
	restoreFromLastCP bool // restore from last checkpoint
}

func newCheckPoint(filename string) *CheckPoint {
	cp := &CheckPoint{
		path:              filename,
		restoredFiles:     make(map[string]struct{}),
		restoreFromLastCP: false,
	}
	if err := cp.load(); err != nil {
		log.Fatal("recover from check point failed, %v", err)
	}
	return cp
}

func (cp *CheckPoint) IsRestoreFromLastCheckPoint() bool {
	return cp.restoreFromLastCP
}

func (cp *CheckPoint) load() error {
	f, err := os.Open(cp.path)
	if err != nil && !os.IsNotExist(err) {
		return errors.Trace(err)
	}
	if os.IsNotExist(err) {
		return nil
	}
	defer f.Close()

	cp.restoreFromLastCP = true
	br := bufio.NewReader(f)
	for {
		line, err := br.ReadString('\n')
		if err == io.EOF {
			break
		}
		l := strings.TrimSpace(line[:len(line)-1])
		if len(l) == 0 {
			continue
		}

		cp.restoredFiles[l] = struct{}{}
	}

	return nil
}

func (cp *CheckPoint) Save(filename string) error {
	cp.Lock()
	defer cp.Unlock()

	cp.restoredFiles[filename] = struct{}{}

	// add to checkpoint file
	f, err := os.OpenFile(cp.path, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		return errors.Trace(err)
	}
	defer f.Close()

	if _, err := fmt.Fprintf(f, "%s\n", filename); err != nil {
		return errors.Trace(err)
	}

	return nil
}

func (cp *CheckPoint) Dump() map[string]struct{} {
	cp.RLock()
	defer cp.RUnlock()

	m := make(map[string]struct{})
	for file, _ := range cp.restoredFiles {
		m[file] = struct{}{}
	}

	return m
}
