package miner

import (
	"io/ioutil"
	"log"
	"strings"

	"github.com/blockassets/bam_agent/service/agent"
	"go.uber.org/fx"
)

type Version struct {
	V string
}

//
func readFileTrim(file string) (*string, error) {
	dat, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	res := strings.TrimSpace(string(dat))
	return &res, nil
}

/*
	BW saves their cgminer version into a file.
*/
func readBWVersionFile() string {
	str, err := readFileTrim("/usr/app/version.txt")
	if err != nil {
		log.Println(err)
		return ""
	}
	return *str
}

/*
	Config is currently unused, but in the future on other miners,
	we might get the version information from another source and we
	would likely pass that in via the agent config.
*/
func NewVersion(config agent.Config) Version {
	bwVersion := readBWVersionFile()

	return Version{
		V: bwVersion,
	}
}

var VersionModule = fx.Provide(NewVersion)
