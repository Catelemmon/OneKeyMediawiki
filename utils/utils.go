package utils

import (
	"errors"
	"fmt"
	"io/ioutil"
	"math"
	"net"
	"os"
	exec2 "os/exec"
	"regexp"
	"strconv"
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

func VerBEThan(srcVer string, tgtVer string) bool{
	srcSubVer := strings.Split(srcVer, ".")
	tgtSubVer := strings.Split(tgtVer, ".")
	for i := 0; math.Min(float64(len(srcSubVer)), float64(len(tgtSubVer))) > float64(i); i++{
		v1, _ := strconv.Atoi(srcSubVer[i])
		v2, _ := strconv.Atoi(tgtSubVer[i])
		if v1 > v2{
			return true

		} else {
			continue
		}
	}
	return len(srcSubVer) > len(tgtSubVer)
}

func FileExist(file string) (error, bool){
	if _, err := os.Stat(file); err == nil{
		return nil, true
	} else if os.IsNotExist(err){
		return nil, false
	} else{
		return err, false
	}
}

func GetIp() ([]string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil{
		return []string{} , nil
	}
	strAddrs := make([]string, 0, 5)
	for _, addr := range addrs{
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback(){
			if ipnet.IP.To4() != nil{
				strAddrs = append(strAddrs, ipnet.IP.String())
			}
		}
	}
	return strAddrs, nil
}