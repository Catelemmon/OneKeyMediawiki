package wikiUtil

import (
	"errors"
	"fmt"
	"github.com/Catelemmon/oneKeyMediawiki/checker"
	"github.com/Catelemmon/oneKeyMediawiki/utils"
	"io/ioutil"
	"os"
	exec2 "os/exec"
	"path"
	"strings"
)

func InstallWiki(imgDir string, scriptsDir string) error{

	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	// create wiki work directory
	wikiRoot := path.Join(cwd, "mediawikiRoot")
	if _, err := os.Stat(wikiRoot); os.IsNotExist(err){
		err = os.Mkdir(wikiRoot, 0755)
		if err != nil{
			return err
		}
	}
	// load images
	err = LoadDockerImage(imgDir)
	if err != nil{
		return err
	}
	// create essential files and directory
	wikiImgDir := path.Join(wikiRoot, "images")
	if _, err := os.Stat(wikiImgDir); os.IsNotExist(err){
		err = os.Mkdir(wikiImgDir, 0755)
		if err != nil{
			return err
		}
	}
	dbDataDir := path.Join(wikiRoot, "dbDataDir", "mariadb")
	if _, err := os.Stat(dbDataDir); os.IsNotExist(err){
		err = os.MkdirAll(dbDataDir, 0755)
		if err != nil{
			return err
		}
	}
	// write compose file
	si := checker.NewSupportInfo()
	si.GetSupportInfo()
	var composeFile string
	if utils.VerBEThan("1.10.1", si.DockerComposeVersion){
		// copy docker compose file v1
		composeFile = path.Join(scriptsDir, "docker-compose-v1.yml")

	} else {
		composeFile = path.Join(scriptsDir, "docker-compose-v3.yml")
	}
	fs, err := ioutil.ReadFile(composeFile)
	if err != nil {
		return errors.New("compose file not found")
	}
	err = ioutil.WriteFile(path.Join(wikiRoot, "docker-compose.yml"), fs, 0666)
	if err != nil {
		return err
	}
	// execute docker compose
	cmd := exec2.Command("docker-compose", "up", "-d")
	cmd.Dir = wikiRoot
	_, err = cmd.Output()
	if err != nil {
		fmt.Println("启动mediawiki失败")
		fmt.Printf("err as follow:\n%s\noutput as follow:\n%s")
		return err
	}
	fmt.Println("启动mediawiki成功，请访问本机IP地址或IP地址加上8080...")
	return nil
}

func LoadDockerImage(imgDir string) error{

	imgs, err := ioutil.ReadDir(imgDir)
	if err != nil {
		return err
	}
	for _, img := range imgs{
		if strings.HasSuffix(img.Name(), "tar"){
			cmd := exec2.Command("docker", "load", "-i" , img.Name())
			cmd.Dir = imgDir
			out, err := cmd.Output()
			if err != nil{
				fmt.Printf("Failed load image file\nerr: %v\noutput:%v\n", err, out)
				return err
			}
		}
	}

	return nil
}



