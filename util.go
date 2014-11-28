package main

import (
	"archive/zip"
	"crypto/md5"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func checkMD5(path string, md5Path string) error {
	f, err := os.Open(path)

	if err != nil {
		return err
	}

	defer f.Close()

	checkMD5, err := readFile(md5Path)

	if err != nil {
		return err
	}

	h := md5.New()

	io.Copy(h, f)

	md5 := fmt.Sprintf("%x", h.Sum(nil))

	if checkMD5 != md5 {
		return errors.New("MD5 checksum failed")
	}

	return nil
}

func downloadFile(url string, dest string) error {
	f, err := os.Create(dest)

	if err != nil {
		return err
	}

	defer f.Close()

	resp, err := http.Get(url)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New(resp.Status)
	}

	_, err = io.Copy(f, resp.Body)

	return err
}

func pathExists(path string) bool {
	_, err := os.Stat(path)

	if os.IsNotExist(err) {
		return false
	}

	return true
}

func readFile(path string) (string, error) {
	b, err := ioutil.ReadFile(path)

	if err != nil {
		return "", err
	}

	return string(b), nil
}

func unzip(zipPath string, dest string) error {
	err := os.RemoveAll(dest)

	if err != nil {
		return err
	}

	zipReader, err := zip.OpenReader(zipPath)

	if err != nil {
		return err
	}

	defer zipReader.Close()

	for _, entry := range zipReader.File {
		destEntryPath := filepath.Join(dest, entry.Name)

		isDir := entry.FileInfo().IsDir()

		var destEntryDir string

		if isDir {
			destEntryDir = destEntryPath
		} else {
			pos := strings.LastIndex(destEntryPath, string(os.PathSeparator))

			if pos > -1 {
				destEntryDir = destEntryPath[:pos]
			}
		}

		if len(destEntryDir) > 0 {
			err := os.MkdirAll(destEntryDir, 0777)

			if err != nil {
				return err
			}
		}

		if !isDir {
			entryReader, err := entry.Open()

			if err != nil {
				return err
			}

			defer entryReader.Close()

			destEntryWriter, err := os.Create(destEntryPath)

			if err != nil {
				return err
			}

			defer destEntryWriter.Close()

			_, err = io.Copy(destEntryWriter, entryReader)

			if err != nil {
				return err
			}
		}
	}

	return nil
}
