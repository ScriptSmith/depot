package main

import (
	"compress/flate"
	"errors"
	"fmt"
	"github.com/mholt/archiver"
	"net/http"
	"os"
	"path"
	"path/filepath"
)

type PathWalkerArgs struct {
	filePath    string
	info        os.FileInfo
	err         error
	z           archiver.Zip
	dirPath     string
	dirPathInfo os.FileInfo
}

// Add discovered file to zip
func pathWalker(args PathWalkerArgs) error {

	if args.err != nil {
		return args.err
	}

	internalName, err := archiver.NameInArchive(args.dirPathInfo, args.dirPath, args.filePath)
	if err != nil {
		return err
	}

	file, err := os.Open(args.filePath)
	if err != nil {
		return err
	}

	err = args.z.Write(archiver.File{
		FileInfo: archiver.FileInfo{
			FileInfo:   args.info,
			CustomName: internalName,
		},
		ReadCloser: file,
	})
	if err != nil {
		return err
	}

	err = file.Close()
	if err != nil {
		return err
	}

	return nil
}

// Create streamed zip file of a job
func ZipDir(dirName string, root string, reqW http.ResponseWriter, _ *http.Request) error {
	dirPath := path.Join(root, dirName)
	dirPathInfo, err := os.Stat(dirPath)
	if err != nil {
		return errors.New("reading job dir")
	}

	z := archiver.Zip{
		CompressionLevel:       flate.NoCompression,
		MkdirAll:               true,
		SelectiveCompression:   true,
		ContinueOnError:        false,
		OverwriteExisting:      false,
		ImplicitTopLevelFolder: false,
	}

	err = z.Create(reqW)
	if err != nil {
		return errors.New("creating zip file")
	}

	defer func() {
		_ = z.Close()
	}()

	reqW.Header().Set(
		"Content-Disposition",
		fmt.Sprintf("attachment; filename=%s.zip", dirName),
	)

	err = filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		return pathWalker(PathWalkerArgs{
			filePath:    path,
			info:        info,
			err:         err,
			z:           z,
			dirPath:     dirPath,
			dirPathInfo: dirPathInfo,
		})

	})
	if err != nil {
		return errors.New("reading job directory")
	}
	return nil
}
