package main

import (
	"io/ioutil"
	"strings"
)

func parseFile(fileName string) ([]string, error) {
	file, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	var fileSlice stringsAlias
	fileSlice = strings.Split(string(file), "\n")
	return fileSlice.removeEmpty(), nil
}

func (s stringsAlias) removeEmpty() []string {
	for i, element := range s {
		if i < len(s)-1 {
			if element == "" {
				s = append(s[:i], s[i+1:]...)
			} else {
				s = s[:i]
			}
		}
	}

	return s
}
