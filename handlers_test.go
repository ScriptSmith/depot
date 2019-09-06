package main

import (
	zip2 "archive/zip"
	"bytes"
	"fmt"
	"github.com/google/uuid"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"testing"
)

var RootTest = path.Join(os.TempDir(), "depot_test_root")

func TestRoot(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatalf("Failed sending GET request: %s", err)
	}

	handlers := &Handlers{root: RootTest}
	router := getRouter(handlers)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Fatalf("Incorrect status code %v: %s", status, rr.Body)
	}
}

func jobFilesSetup(t *testing.T) (string, string, int64) {
	jobId := uuid.New()
	t.Logf("Job id: %s", jobId)
	return jobId.String(), "http://speedcheck.cdn.on.net/10meg.test", int64(10000000)
}

func TestJobs(t *testing.T) {
	// Setup
	jobId, fileUploadUrl, fileSize := jobFilesSetup(t)
	fileDownloadPath := "newdir/test.out"
	handlers := &Handlers{root: RootTest, deletion: true}
	router := getRouter(handlers)

	// Get test file
	resp, err := http.Get(fileUploadUrl)
	if err != nil {
		t.Fatalf("Failed downloading test file: %s", err)
	}
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			t.Fatalf("Failed closing test file: %s", err)
		}
	}()

	// Create test file body stream
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed reading test file: %s", err)
	} else if size := resp.ContentLength; size != fileSize {
		t.Fatalf("Incorrect downloading file size: %d not %d", size, fileSize)
	}
	fileBody := bytes.NewReader(body)

	// Upload body stream
	newJobUrl := fmt.Sprintf("/jobs/%s/", jobId)
	newFileUrl := fmt.Sprintf("%s%s", newJobUrl, fileDownloadPath)
	t.Logf(newFileUrl)
	req, err := http.NewRequest("PUT", newFileUrl, fileBody)
	if err != nil {
		t.Fatalf("Failed sending PUT request: %s", err)
	}
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	// Check response
	if status := rr.Code; status != http.StatusOK {
		t.Fatalf("Incorrect status code %d: %s", status, rr.Body)
	}

	// Check file exists in job dir
	info, err := os.Stat(path.Join(RootTest, jobId, fileDownloadPath))
	if err != nil {
		t.Fatalf("Error reading file: %s", err)
	} else if size := info.Size(); size != fileSize {
		t.Fatalf("Incorrect saved file size: %d not %d", size, fileSize)
	}

	// Get the new file
	req, err = http.NewRequest("GET", newFileUrl, nil)
	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	rrBody, err := ioutil.ReadAll(rr.Body)
	if err != nil {
		t.Fatalf("Failed reading GET response body: %s", err)
	}

	// Check it's valid
	if status := rr.Code; status != http.StatusOK {
		t.Fatalf("Incorrect status code %d: %s", status, rr.Body)
	} else if size := int64(len(rrBody)); size != fileSize {
		t.Fatalf("Incorrect saved file size: %d not %d", size, fileSize)
	}

	// Delete it
	req, err = http.NewRequest("DELETE", newFileUrl, nil)
	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	rrBody, err = ioutil.ReadAll(rr.Body)
	if err != nil {
		t.Fatalf("Failed reading DELETE response body: %s", err)
	} else if status := rr.Code; status != http.StatusOK {
		t.Fatalf("Incorrect status code %d: %s", status, rr.Body)
	}

	// Check file doesn't exist
	info, err = os.Stat(path.Join(RootTest, jobId, fileDownloadPath))
	if err == nil {
		t.Fatalf("File still exists: %s", err)
	}

	// Delete job
	req, err = http.NewRequest("DELETE", newJobUrl, nil)
	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	rrBody, err = ioutil.ReadAll(rr.Body)
	if err != nil {
		t.Fatalf("Failed reading DELETE response body: %s", err)
	} else if status := rr.Code; status != http.StatusOK {
		t.Fatalf("Incorrect job status code %d: %s", status, rr.Body)
	}

	// Check job doesn't exist
	info, err = os.Stat(path.Join(RootTest, jobId))
	if err == nil {
		t.Fatalf("Job still exists: %s", err)
	}
}

func TestZip(t *testing.T) {
	// Setup
	jobId, fileUploadUrl, fileSize := jobFilesSetup(t)
	fileDir := "testdir"
	handlers := &Handlers{root: RootTest}
	router := getRouter(handlers)

	// Download three files
	for i := 1; i <= 3; i++ {
		resp, _ := http.Get(fileUploadUrl)
		body, _ := ioutil.ReadAll(resp.Body)
		_ = resp.Body.Close()
		fileBody := bytes.NewReader(body)

		newFileUrl := fmt.Sprintf("/jobs/%s/%s/%d.out", jobId, fileDir, i)
		req, _ := http.NewRequest("PUT", newFileUrl, fileBody)

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Fatalf("Incorrect status code %d: %s", status, rr.Body)
		}
	}

	// Check files are added
	for i := 1; i <= 3; i++ {
		filePath := path.Join(RootTest,
			jobId,
			fileDir,
			fmt.Sprintf("%d.out", i),
		)

		info, err := os.Stat(filePath)
		if err != nil {
			t.Fatalf("Error reading file %s: %s", filePath, err)
		} else if info.Size() != fileSize {
			t.Fatalf("File has incorrect size %d not %d", info.Size(), fileSize)
		}
	}

	// Get zipped copy
	zipUrl := fmt.Sprintf("/zip?id=%s", jobId)
	req, err := http.NewRequest("GET", zipUrl, nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	rrBody, err := ioutil.ReadAll(rr.Body)
	if err != nil {
		t.Fatalf("Failed reading zip body: %s", err)
	}

	// Read zip file
	zipReader, err := zip2.NewReader(bytes.NewReader(rrBody), int64(len(rrBody)))
	if err != nil {
		t.Fatalf("Failed creating zip reader: %s", err)
	}

	correctFileDir := path.Join(jobId, fileDir)
	for i := 2; i < len(zipReader.File); i++ {
		zippedFileDir, zippedFileName := path.Split(zipReader.File[i].Name)
		zippedFileDir = path.Clean(zippedFileDir)
		zippedFileSize := zipReader.File[i].UncompressedSize64
		correctFileName := fmt.Sprintf("%d.out", i-1)

		if zippedFileName != correctFileName {
			t.Fatalf("File names are incorrect: %s not %s", zippedFileName, correctFileName)
		} else if zippedFileDir != correctFileDir {
			t.Fatalf("File directory is incorrect: %s not %s", zippedFileDir, correctFileDir)
		} else if zippedFileSize != uint64(fileSize) {
			t.Fatalf("File sizes are incorrect: %d not %d", zippedFileSize, fileSize)
		}
	}
}

func TestMain(m *testing.M) {
	var err error
	RootTest, err = ioutil.TempDir("", "depot_test_")
	if err != nil {
		log.Fatalln("Couldn't create temp dir")
	}

	log.Printf("Using %s", RootTest)
	m.Run()

	err = os.RemoveAll(RootTest)
	if err != nil {
		log.Fatalln("Failed cleaning up temp dir")
	}
}
