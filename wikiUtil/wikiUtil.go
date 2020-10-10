package wikiUtil

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/Catelemmon/oneKeyMediawiki/checker"
	"github.com/Catelemmon/oneKeyMediawiki/utils"
	"io"
	"io/ioutil"
	"os"
	exec2 "os/exec"
	"path"
	"strings"
)

func LaunchWiki(wikiRoot string) error {
	si := checker.NewSupportInfo()
	si.GetSupportInfo()
	dockerCmpVersion := si.DockerComposeVersion
	var composeArg string
	version1Flag := false
	if utils.VerBEThan(dockerCmpVersion, "1.13.0"){
		composeArg = "up -d"
	} else if utils.VerBEThan(dockerCmpVersion, "1.10.0"){
		fmt.Println("cannot support docker-compose version 2")
		return errors.New("cannot support docker-compose version 2")
	} else {
		composeArg = "-f docker-compose-v1-default.yml up -d"
		version1Flag = true
	}
	if version1Flag{
		// TODO: render version1 compose file


	}
	cmd := exec2.Command("docker-compose",  strings.Split(composeArg, " ")...)
	cmd.Dir = wikiRoot
	out, err := cmd.Output()
	if err != nil {
		fmt.Printf("Failed to Launch Wiki\n")
		return err
	}
	outs := strings.TrimSpace(string(out))
	if outs != ""{
		fmt.Println("Launch log:", outs)
	}
	return errors.New("launch wiki failed")
}

func ShutDownWiki(wikiRoot string) error {
	cmd := exec2.Command("docker-compose", "down")
	cmd.Dir = wikiRoot
	out, err := cmd.Output()
	if err != nil {
		fmt.Printf("Failed to shutdown Wiki\n")
		return err
	}
	outs := strings.TrimSpace(string(out))
	if outs != ""{
		fmt.Println("Shutdown log:", outs)
	}
	return errors.New("shutdown wiki failed")
}

func CompatibleDockerCmpVersion1(wikiRoot string){

	envStore := make(map[string]string)

	envStore["EMERGENCY_CONTACT"] = ""
	envStore["PASSWORD_SENDER"] = ""
	envStore["SECRET_KEY"] = ""
	envStore["UPGRADE_KEY"] = ""
	envStore["DB_NAME"] = ""
	envStore["DB_USER"] = ""
	envStore["DB_PASSWORD"] = ""

	envFi, _ := os.Open(path.Join(wikiRoot, ".env"))
	defer envFi.Close()
	reader := bufio.NewReader(envFi)
	for true {
		bl, _, err := reader.ReadLine()
		if err != io.EOF{
			break
		}
		line := string(bl)
		pieces := strings.Split(line, "=")
		envName, envValue := pieces[0], pieces[1]
		envStore[envName] = envValue  // read .env
	}
	// render docker-compose version1
	cmpFi, _ := os.Open(path.Join(wikiRoot, "docker-compose-v1-default.yml"))
	defer cmpFi.Close()
	dcTempb, _ := ioutil.ReadAll(cmpFi)
	dcTemp := string(dcTempb)
	for k, v := range envStore{
		dcTemp = strings.Replace(dcTemp, fmt.Sprintf("${%s}", k), v, -1)
	}
	dcV1Fi, _ := os.OpenFile(path.Join(wikiRoot, "docker-compose-v1-default.yml"), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	defer dcV1Fi.Close()
	dcV1Fi.Write([]byte(dcTemp))
}