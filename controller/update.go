package controller

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/blockassets/bam_agent/tool"
	"github.com/json-iterator/go"
	"github.com/mholt/archiver"
)

const (
	defaultUpdateScriptName = "update.sh"
	execTimeout = time.Duration(60) * time.Second
)

func NewUpdateCtrl() Result {
	return Result{
		Controller: &Controller{
			Path:    "/update",
			Methods: []string{http.MethodPost},
			Handler: tool.JsonHandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				var status = http.StatusOK
				var response []byte
				var err error

				script := r.FormValue("script")
				if len(script) == 0 {
					script = defaultUpdateScriptName
				}

				var gzFileData multipart.File
				gzFileData, _, err = r.FormFile("file")
				if err == nil {
					response, err = unZipAndUpdate(gzFileData, script)
					if err == nil {
						w.Header().Set("Content-Type", "text/plain")
					}
				}

				if err != nil {
					status = http.StatusInternalServerError
					response, _ = jsoniter.Marshal(BAMStatus{Status: "Error", Error: err, Message: fmt.Sprintf("%s\n%s", response, err)})
				}

				w.WriteHeader(status)
				w.Write(response)
			}),
		},
	}
}

func unZipAndUpdate(file io.Reader, script string) ([]byte, error) {
	var err error
	var updateSh *string

	tempDir, err := ioutil.TempDir("", "update-unzip")
	defer os.RemoveAll(tempDir)

	if err == nil {
		err = archiver.TarGz.Read(file, tempDir)
		if err == nil {
			updateSh, err = findUpdateScript(tempDir, script)
			if err == nil {
				ctx, cancel := context.WithTimeout(context.Background(), execTimeout)
				defer cancel()

				cmd := exec.CommandContext(ctx, "/bin/sh", *updateSh)
				dir, _ := filepath.Split(*updateSh)
				cmd.Dir = dir

				return cmd.CombinedOutput()
			}
		}
	}

	return nil, err
}

func findUpdateScript(dir string, script string) (*string, error) {
	var updatePath *string

	isDir, err := tool.IsDirectory(dir)
	if err != nil || !isDir {
		return nil, errors.New(fmt.Sprintf("expected directory: %s  error: %v", dir, err))
	}

	filepath.Walk(dir, filepath.WalkFunc(func(path string, info os.FileInfo, err error) error {
		// exit early
		if updatePath != nil {
			return filepath.SkipDir
		}

		if info.Name() == script && tool.IsExecutable(info.Mode()) {
			updatePath = &path
			return filepath.SkipDir // exit early
		}

		return nil
	}))

	if updatePath == nil {
		return nil, errors.New(fmt.Sprintf("could not find executable update.sh in %s", dir))
	}

	return updatePath, nil
}
