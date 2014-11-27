package main

import (
	"errors"
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
	err := checkJava()

	fatalIfError(err)

	homeDir, err := getHomeDir()

	fatalIfError(err)

	cfg, err := readConfig(homeDir)

	fatalIfError(err)

	cfg.save(homeDir)
}

func checkJava() error {
	javaHome := os.Getenv("JAVA_HOME")

	if len(javaHome) == 0 {
		return errors.New("The JAVA_HOME environment variable is not defined.")
	}

	names := []string{"java", "javac"}

	for _, n := range names {
		if runtime.GOOS == "windows" {
			n += ".exe"
		}

		n = filepath.Join(javaHome, "bin", n)

		if !fileExists(n) {
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
	u, err := user.Current()

	if err != nil {
		return "", err
	}

	homeDir := filepath.Join(u.HomeDir, homeName)

	err = os.MkdirAll(homeDir, 0777)

	return homeDir, err
}
