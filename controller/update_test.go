package controller

import (
	"bytes"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mholt/archiver"
)

func TestDefaultUpdateScriptName(t *testing.T) {
	if defaultUpdateScriptName != "update.sh" {
		t.Fatalf("expected update.sh and got %s", defaultUpdateScriptName)
	}
}

func TestNewUpdateCtrl(t *testing.T) {
	result := NewUpdateCtrl()
	ctrl := result.Controller

	if ctrl.Path != "/update" {
		t.Fatalf("expected path /update, got %s", ctrl.Path)
	}

	if len(ctrl.Methods) != 1 {
		t.Fatalf("expected 1 method, got %d", len(ctrl.Methods))
	}

	if ctrl.Methods[0] != http.MethodPost {
		t.Fatalf("expected method post, got %s", ctrl.Methods[0])
	}

	req := httptest.NewRequest("POST", "/doesnotmatter", nil)
	w := httptest.NewRecorder()
	ctrl.Handler.ServeHTTP(w, req)

	if w.Result().StatusCode != http.StatusInternalServerError {
		t.Fatalf("expected internal server error")
	}
}

func TestFileUpload(t *testing.T) {
	result := NewUpdateCtrl()
	ctrl := result.Controller

	tgz, err := makeTarGz(defaultUpdateScriptName)
	defer os.Remove(tgz.Name())
	if err != nil {
		t.Fatal(err)
	}

	body, contentType, err := createFormUpload(tgz, defaultUpdateScriptName)
	if err != nil || contentType == nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest("POST", "/doesnotmatter", body)
	req.Header.Set("Content-Type", *contentType)

	w := httptest.NewRecorder()
	ctrl.Handler.ServeHTTP(w, req)

	resResult := w.Result()
	bodyData, _ := ioutil.ReadAll(resResult.Body)
	bodyStr := string(bodyData)
	if !strings.Contains(bodyStr, "got here") {
		t.Fatalf("expected 'got here' and got %s", bodyStr)
	}
}

func TestFileUploadDefaultScriptName(t *testing.T) {
	commonFileUploadTest(t, defaultUpdateScriptName)
}

func TestFileUploadDifferentScriptName(t *testing.T) {
	commonFileUploadTest(t, "foo.sh")
}

func TestFileUploadScriptExecutionFailure(t *testing.T) {
	commonFileUploadTest(t, "fail.sh")
}

func commonFileUploadTest(t *testing.T, scriptName string) {
	result := NewUpdateCtrl()
	ctrl := result.Controller

	tgz, err := makeTarGz(scriptName)
	defer os.Remove(tgz.Name())
	if err != nil {
		t.Fatal(err)
	}

	body, contentType, err := createFormUpload(tgz, scriptName)
	if err != nil || contentType == nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest("POST", "/doesnotmatter", body)
	req.Header.Set("Content-Type", *contentType)

	w := httptest.NewRecorder()
	ctrl.Handler.ServeHTTP(w, req)

	resResult := w.Result()
	bodyData, _ := ioutil.ReadAll(resResult.Body)
	bodyStr := string(bodyData)

	if scriptName == "fail.sh" {
		if !strings.Contains(bodyStr, "failed") {
			t.Fatalf("expected failed, got %s", bodyStr)
		}

		if !strings.Contains(bodyStr, "exit status 1") {
			t.Fatalf("expected exit status 1, got %s", bodyStr)
		}

		if resResult.StatusCode != http.StatusInternalServerError {
			t.Fatalf("expected 500, got %v", resResult.StatusCode)
		}

		if resResult.Header.Get("Content-Type") != "application/json; charset=utf-8" {
			t.Fatalf("expected application/json; charset=utf-8, got %v", resResult.Header.Get("Content-Type"))
		}
	} else {
		if !strings.Contains(bodyStr, "got here") {
			t.Fatalf("expected 'got here' and got %s", bodyStr)
		}

		if resResult.StatusCode != http.StatusOK {
			t.Fatalf("expected 200, got %v", resResult.StatusCode)
		}

		if resResult.Header.Get("Content-Type") != "text/plain" {
			t.Fatalf("expected text/plain, got %v", resResult.Header.Get("Content-Type"))
		}
	}
}

func createFormUpload(file *os.File, script string) (*bytes.Buffer, *string, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	defer writer.Close()

	part, err := writer.CreateFormFile("file", file.Name())
	if err != nil {
		return nil, nil, err
	}
	_, err = io.Copy(part, file)
	if err != nil {
		return nil, nil, err
	}

	err = writer.WriteField("script", script)
	if err != nil {
		return nil, nil, err
	}

	ct := writer.FormDataContentType()
	return body, &ct, nil
}

func TestUnzipAndUpdate(t *testing.T) {
	tgz, err := makeTarGz(defaultUpdateScriptName)
	defer os.Remove(tgz.Name())
	if err != nil {
		t.Fatal(err)
	}

	response, err := unZipAndUpdate(tgz, defaultUpdateScriptName)
	if err != nil {
		t.Fatal(err)
	}

	responseStr := string(response)

	// Got this output from the script
	if !strings.Contains(responseStr, "got here") {
		t.Fatalf("expected 'got here' and got %s", responseStr)
	}

	// Test that we are setting the current working directory of the update script.
	// Should be the same as the temp directory that is created for the .tar file
	if !strings.Contains(responseStr, "update-test") {
		t.Fatalf("expected pwd output to have update-test in it and got %s", responseStr)
	}
}

func TestFindUpdateScript(t *testing.T) {
	tgz, err := makeTarGz(defaultUpdateScriptName)
	defer os.Remove(tgz.Name())
	if err != nil {
		t.Fatal(err)
	}

	tempDir, err := ioutil.TempDir("", "update-test")
	defer os.RemoveAll(tempDir)

	err = archiver.TarGz.Open(tgz.Name(), tempDir)
	if err != nil {
		t.Fatal(err)
	}

	script, err := findUpdateScript(tempDir, defaultUpdateScriptName)
	if err != nil {
		t.Fatal(err)
	}

	_, file := path.Split(*script)
	if file != defaultUpdateScriptName {
		t.Fatalf("expected script to be named update.sh and got %s", file)
	}
}

const (
	updateSh = `#!/usr/bin/env bash

echo $(pwd)
echo "got here"
`

	failSh = `#!/usr/bin/env bash

echo "failed"
exit 1
`
)

func makeTarGz(script string) (*os.File, error) {
	tempDir, err := ioutil.TempDir("", "update-test")
	defer os.RemoveAll(tempDir)
	if err != nil {
		return nil, err
	}

	update := filepath.Join(tempDir, script)

	scriptContents := []byte(updateSh)
	// Special test case for script failure
	if script == "fail.sh" {
		scriptContents = []byte(failSh)
	}

	err = ioutil.WriteFile(update, scriptContents, 0511)
	if err != nil {
		return nil, err
	}

	file, err := ioutil.TempFile("", "update-tgz")
	if err != nil {
		return nil, err
	}
	err = archiver.TarGz.Make(file.Name(), []string{tempDir})
	if err != nil {
		return nil, err
	}

	return file, nil
}
