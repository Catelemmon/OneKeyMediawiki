package utils

import (
	"errors"
	"fmt"
	"io/ioutil"
	exec2 "os/exec"
	"regexp"
	"strings"
)

var versionReg, _ = regexp.Compile(`([1-9]\d|[1-9])(\.(([1-9]\d{1,2})|\d)){1,2}`)

func VersionAllow(srcVer string, tarVer string) bool {
	if srcVer < tarVer {
		return false
	} else {
		return true
	}
}

func VersionExtract(s string) (string, error) {
	version := versionReg.FindString(s)
	if versionReg.FindString(s) == "" {
		return "", errors.New("edition number is not found")
	} else {
		return version, nil
	}
}

func HasCommand(cmd string) bool {
	exec := exec2.Command("whereis", cmd)
	out, err := exec.Output()
	if err != nil {
		panic(err)
	}
	pieces := strings.Split(string(out), ":")
	if len(pieces) == 1 || len(strings.TrimSpace(pieces[1])) == 0 {
		return false
	} else {
		return true
	}
}

func CommandVersion(cmd string) (string, error) {
	for _, vArg := range []string{"-V", "-v", "--version"} {
		exec := exec2.Command(cmd, fmt.Sprintf("%s", vArg))
		out, _ := exec.Output()
		version, err := VersionExtract(string(out))
		if err != nil {
			continue
		} else {
			return version, nil
		}
	}
	return "", errors.New(fmt.Sprintf("command %s's version cannot be checked", cmd))
}

func HasFile(checkDir string, filename string) bool{
	files, err := ioutil.ReadDir(checkDir)
	if err != nil{
		return false
	}
	for _, file := range files{
		if file.Name() ==filename && !file.IsDir(){
			return true
		}
	}
	return false
}

