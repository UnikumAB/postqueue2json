/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package main

import (
	"bufio"
	"encoding/json"
	"os"

	"github.com/UnikumAB/postqueue2json/filter"
	"github.com/sirupsen/logrus"
)

func main() {
	var filename string
	if len(os.Args) == 2 {
		filename = os.Args[1]
	}
	lines := make(chan string)
	go func() {
		defer close(lines)
		var file *os.File
		if filename != "" {
			file, err := os.Open(filename)
			if err != nil {
				return
			}
			defer file.Close()
		} else {
			file = os.Stdin
		}
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			lines <- scanner.Text()
		}
	}()
	queueItems, err := filter.ConvertPostqueueToQueueItem(lines)
	if err != nil {
		logrus.Fatalf("Failed to start converting: %v", err)
	}
	var items []filter.QueueItem
	items = []filter.QueueItem{}
	for item := range queueItems {
		items = append(items, item)
	}
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	err = encoder.Encode(items)
	if err != nil {
		logrus.Fatalf("failed to encode item %v: %v", items, err)
	}
}
