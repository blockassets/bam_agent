package service

import (
	"io/ioutil"
	"log"
	"strings"
)

//
func readFileTrim(file string) string {
	dat, err := ioutil.ReadFile(file)
	if err != nil {
		log.Println(err)
	}

	return strings.TrimSpace(string(dat))
}

/*
	BW saves their cgminer version into a file.

	Future miners can just check alternative locations in this method.
*/
func ReadVersionFile() string {
	return readFileTrim("/usr/app/version.txt")
}
