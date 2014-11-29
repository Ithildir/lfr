package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"runtime"
	"strings"
)

const (
	blank string = ""
	dash string = "-"
	homeName string = ".lfr"
	space string = " "
)

var (
	argVersion string
	argURL     string
)

func init() {
	flag.StringVar(&argURL, "url", defaultURL, "URL for downloading packages")
	flag.StringVar(&argVersion, "version", blank, "version of the package to use (empty for current)")
}

func main() {
	flag.Parse()

	err := checkJava()

	fatalIfError(err, blank)

	homeDir, err := getHomeDir()

	fatalIfError(err, "Unable to get home directory")

	cfg, err := readConfig(homeDir)

	fatalIfError(err, "Unable to read configuration")

	if len(argVersion) > 0 {
		err = downloadPackage(homeDir, argURL, argVersion)

		fatalIfError(err, ("Unable to download package version " + argVersion))
	} else {
		err = update(&cfg, homeDir, argURL)

		if err != nil {
			if packageExist(homeDir, cfg.Version) {
				msg := fmt.Sprint("Unable to update to current version (", err.Error(), "), using ", cfg.Version, " instead")

				fmt.Println(msg)
			} else {
				fatalIfError(err, "Unable to update to current version")
			}
		}

		argVersion = cfg.Version
	}

	err = cfg.save(homeDir)

	fatalIfError(err, "Unable to save configuration")

	fmt.Println("Using package version " + argVersion)

	err = execute(homeDir, argVersion, flag.Args())

	fatalIfError(err, "Unable to execute command")
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

func execute(homeDir, version string, args []string) error {
	packagePath := getPackagePath(homeDir, version, blank)
	prefixPath := filepath.Join(packagePath, "PREFIX")

	prefix, err := pathToString(prefixPath)

	if err != nil {
		return err
	}

	tokens := strings.Split(prefix, space)

	path := filepath.Join(packagePath, tokens[0])

	args = append(tokens[1:], args...)

	cmd := exec.Command(path, args...)

	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	return cmd.Run()
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
		return blank, err
	}

	homeDir := filepath.Join(user.HomeDir, homeName)

	err = os.MkdirAll(homeDir, 0777)

	return homeDir, err
}
