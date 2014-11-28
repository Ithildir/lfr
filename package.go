package main

import (
	"fmt"
	"path/filepath"
	"runtime"
)

const (
	defaultURL    string = "http://localhost:8080/lfr/"
	packagePrefix string = "v"
)

func downloadPackage(homeDir string, url string, version string) error {
	if packageExist(homeDir, version) {
		fmt.Println("Package version " + version + " is already installed")

		return nil
	}

	fmt.Println("Downloading package version " + version)

	zipPath := getPackagePath(homeDir, version, ".zip")
	zipURL := getPackageURL(url, version, ".zip")

	err := downloadFile(zipURL, zipPath)

	if err != nil {
		return err
	}

	md5Path := getPackagePath(homeDir, version, ".md5")
	md5URL := getPackageURL(url, version, ".md5")

	err = downloadFile(md5URL, md5Path)

	if err != nil {
		return err
	}

	err = checkMD5(zipPath, md5Path)

	if err != nil {
		return err
	}

	path := getPackagePath(homeDir, version, "")

	return unzip(zipPath, path)
}

func getPackagePath(homeDir string, version string, ext string) string {
	return filepath.Join(homeDir, (packagePrefix + version + ext))
}

func getPackageURL(url string, version string, ext string) string {
	return fmt.Sprint(url, packagePrefix, version, "-", runtime.GOOS, ext)
}

func packageExist(homeDir string, version string) bool {
	path := getPackagePath(homeDir, version, "")

	return pathExists(path)
}
