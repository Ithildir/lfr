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

	fatalIfError(err)

	homeDir, err := getHomeDir()

	fatalIfError(err)

	cfg, err := readConfig(homeDir)

	fatalIfError(err)

	if len(version) > 0 {
		err = downloadPackage(homeDir, url, version)
	} else {

	}

	fatalIfError(err)

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

func fatalIfError(err error) {
	if err != nil {
		fmt.Println(err)

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
