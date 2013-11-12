// Copyright 2013 Jari Takkala. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"path/filepath"
	"launchpad.net/tomb"
	"log"
	"os"
)

type IO struct {
	metaInfo MetaInfo
	files map[string]*os.File
	t tomb.Tomb
}

func (io *IO) Stop() error {
	io.t.Kill(nil)
	return io.t.Wait()
}

func (io *IO) Init() {
	io.files = make(map[string]*os.File)
	if len(io.metaInfo.Info.Files) > 0 {
		// Multiple File Mode
		directory := io.metaInfo.Info.Name
		// Create the directory if it doesn't exist
		if _, err := os.Stat(directory); os.IsNotExist(err) {
			err = os.Mkdir(directory, os.ModeDir | os.ModePerm)
			if err != nil {
				log.Fatal(err)
			}
		}
		err := os.Chdir(directory)
		if err != nil {
			log.Fatal(err)
		}
		for _, file := range io.metaInfo.Info.Files {
			// Create any sub-directories if required
			if len(file.Path) > 1 {
				directory = filepath.Join(file.Path[1:]...)
				if _, err := os.Stat(directory); os.IsNotExist(err) {
					err = os.MkdirAll(directory, os.ModeDir | os.ModePerm)
					if err != nil {
						log.Fatal(err)
					}
				}
			}
			// Create the file if it doesn't exist
			path := filepath.Join(file.Path...)
			var fh *os.File
			if _, err := os.Stat(path); os.IsNotExist(err) {
				// Create the file and return a handle
				fh, err = os.Create(path)
				if err != nil {
					log.Fatal(err)
				}
			} else {
				// Open the file and return a handle
				fh, err = os.Open(path)
				if err != nil {
					log.Fatal(err)
				}
			}
			io.files[path] = fh
		}
	} else {
		// Single File Mode
	}
}

func (io *IO) Run() {
	log.Println("IO : Run : Started")
	defer io.t.Done()
	defer log.Println("IO : Run : Completed")

	io.Init()
	fmt.Println(io.files)

	for {
		select {
		case <- io.t.Dying():
			return
		}
	}
}