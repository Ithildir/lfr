package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
)

const (
	homeName string = ".lfr"
)

func main() {
	url := *flag.String("url", defaultURL, "URL for downloading packages")
	version := *flag.String("version", "", "version of the package to use (empty for current)")

	flag.Parse()

	err := checkJava()

	fatalIfError(err, "")

	homeDir, err := getHomeDir()

	fatalIfError(err, "Unable to get home directory")

	cfg, err := readConfig(homeDir)

	fatalIfError(err, "Unable to read configuration")

	if len(version) > 0 {
		err = downloadPackage(homeDir, url, version)

		fatalIfError(err, ("Unable to download package version " + version))
	} else {
		err = update(&cfg, homeDir, url)

		if err != nil {
			if packageExist(homeDir, cfg.Version) {
				msg := fmt.Sprint("Unable to update to current version (", err.Error(), "), using ", cfg.Version, " instead")

				fmt.Println(msg)
			} else {
				fatalIfError(err, "Unable to update to current version")
			}
		}

		version = cfg.Version
	}

	fmt.Println("Using package version " + version)

	cfg.save(homeDir)
}

func checkJava() error {
	javaHome := os.Getenv("JAVA_HOME")

	if len(javaHome) == 0 {
		return errors.New("The JAVA_HOME environment variable is not defined.")
	}

	names := []string{"java", "javac"}

	for _, name := range names {
		if runtime.GOOS == "windows" {
			name += ".exe"
		}

		path := filepath.Join(javaHome, "bin", name)

		if !pathExists(path) {
			return errors.New("The JAVA_HOME environment variable is not defined correctly.")
		}
	}

	return nil
}

func fatalIfError(err error, msg string) {
	if err != nil {
		if len(msg) > 0 {
			fmt.Println(msg + ": " + err.Error())
		} else {
			fmt.Println(err)
		}

		os.Exit(1)
	}
}

func getHomeDir() (string, error) {
	user, err := user.Current()

	if err != nil {
		return "", err
	}

	homeDir := filepath.Join(user.HomeDir, homeName)

	err = os.MkdirAll(homeDir, 0777)

	return homeDir, err
}
