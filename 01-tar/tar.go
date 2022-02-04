package main

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
)

func main() {
	pathToTar := os.Args[1]
	outFile := os.Args[2]

	allFiles := gimmeAllFiles(pathToTar)

	for _, file := range allFiles {
		fmt.Println("Current path:", file)
	}

	if err := createTarArchive(allFiles, outFile); err != nil {
		log.Panic(err)
	}

	if err := gzipTarArchive(outFile); err != nil {
		log.Panic(err)
	}
}

func gzipTarArchive(tarFile string) error {
	outFile := tarFile + ".gz"
	fd, err := os.Create(outFile)
	if err != nil {
		return err
	}
	defer fd.Close()

	gz := gzip.NewWriter(fd)
	defer gz.Close()

	td, err := os.Open(tarFile)
	if err != nil {
		return err
	}

	_, err = io.Copy(gz, td)
	if err != nil {
		return err
	}
	return nil
}

func createTarArchive(allFiles []string, outFile string) error {
	out, err := os.Create(outFile)
	if err != nil {
		return err
	}
	defer out.Close()

	tw := tar.NewWriter(out)
	defer tw.Close()

	if err != nil {
		return err
	}

	for _, file := range allFiles {
		if err = addFileToTarArchive(tw, file); err != nil {
			return err
		}
	}
	return nil
}

func addFileToTarArchive(tw *tar.Writer, file string) error {
	fd, err := os.Open(file)
	if err != nil {
		return err
	}
	stats, _ := fd.Stat()

	header, err := tar.FileInfoHeader(stats, stats.Name())
	if err != nil {
		return err
	}
	header.Name = file

	if err = tw.WriteHeader(header); err != nil {
		return err
	}
	_, err = io.Copy(tw, fd)
	if err != nil {
		return err
	}
	return nil
}

func gimmeAllFiles(root string) []string {
	allFiles := make([]string, 0)
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}
		if path == os.Args[2] {
			return nil
		}
		allFiles = append(allFiles, path)
		return nil
	})
	if err != nil {
		log.Panic(err)
	}
	return allFiles
}
