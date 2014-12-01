package main

import (
	"fmt"
	"path/filepath"
	"runtime"
	"time"
)

const (
	defaultURL     string = "http://localhost:8080/lfr/"
	updateInterval int64  = 60 * 60 * 24
)

func downloadPackage(homeDir string, url string, version string) error {
	if packageExist(homeDir, version) {
		fmt.Println("Package version " + version + " is already installed")

		return nil
	}

	fmt.Println("Downloading package version " + version + "...")

	zipPath := getPackagePath(homeDir, version, ".zip")
	zipURL := getPackageURL(url, version, ".zip")

	err := urlToFile(zipURL, zipPath)

	if err != nil {
		return err
	}

	md5Path := getPackagePath(homeDir, version, ".md5")
	md5URL := getPackageURL(url, version, ".md5")

	err = urlToFile(md5URL, md5Path)

	if err != nil {
		return err
	}

	err = checkMD5(zipPath, md5Path)

	if err != nil {
		return err
	}

	path := getPackagePath(homeDir, version, blank)

	return unzip(zipPath, path)
}

func getPackagePath(homeDir string, version string, ext string) string {
	return filepath.Join(homeDir, version+ext)
}

func getPackageURL(url string, version string, ext string) string {
	return fmt.Sprint(url, version, dash, runtime.GOOS, ext)
}

func packageExist(homeDir string, version string) bool {
	if isNull(version) {
		return false
	}

	path := getPackagePath(homeDir, version, blank)

	return pathExists(path)
}

func update(cfg *config, homeDir string, url string) error {
	now := time.Now().Unix()

	if packageExist(homeDir, cfg.Version) && (cfg.LastUpdate >= (now - updateInterval)) {
		return nil
	}

	currentVersion, err := urlToString(url + "CURRENT")

	if err != nil {
		return err
	}

	err = downloadPackage(homeDir, url, currentVersion)

	if err != nil {
		return err
	}

	cfg.LastUpdate = now
	cfg.Version = currentVersion

	return nil
}
