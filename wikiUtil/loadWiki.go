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
	"regexp"
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
		composeFile = path.Join(scriptsDir, "docker-compose-v1-default.yml")

	} else {
		composeFile = path.Join(scriptsDir, "docker-compose.yml")
	}
	fs, err := ioutil.ReadFile(composeFile)
	if err != nil {
		return errors.New("compose file not found")
	}
	err = ioutil.WriteFile(path.Join(wikiRoot, path.Base(composeFile)), fs, 0666)
	if err != nil {
		return err
	}
	// init mediawiki

	cmd := exec2.Command("docker-compose", "up", "-d")
	cmd.Dir = wikiRoot
	out, err := cmd.Output()
	if err != nil {
		fmt.Println("启动mediawiki失败")
		fmt.Printf("err as follow:\n%s\noutput as follow:\n%s", err, string(out))
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
		if strings.HasSuffix(img.Name(), "gz"){
			cmd := exec2.Command("bash", "-c", "gunzip -c " + img.Name() + " | docker load")
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

func InitWiki(wikiRoot, adminEmail, adminName, adminPassword, dbName, dbUser, dbPassword string) error {
	err := RenderEnv(path.Join(wikiRoot, ".env"), adminEmail, "", "", dbName, dbUser, dbPassword, false)
	if err != nil{
		return err
	}
	err = LaunchWiki(wikiRoot)
	if err != nil{
		return err
	}
	// install wiki in container
	cmd := exec2.Command("docker", "exec", "-it", "mediawiki_wiki", "/script/install.sh", adminName, adminPassword)
	cmd.Dir = wikiRoot
	out, err := cmd.Output()
	if err != nil{
		return err
	}
	outs := string(out)
	secretKey, upgradeKey := "", ""
	secretPtn := regexp.MustCompile(`\$wgSecretKey=(\w+)`)
	upgradePtn := regexp.MustCompile(`\$wgUpgradeKey=(\w+)`)
	matches := secretPtn.FindStringSubmatch(outs)
	if len(matches) > 0{
		secretKey = matches[1]
	} else {
		fmt.Printf("failed to install wiki, err content is follow: \n%s", out)
		return errors.New("cannot get secret key")
	}
	matches = upgradePtn.FindStringSubmatch(outs)
	if len(matches) > 0{
		upgradeKey = matches[1]
	} else {
		return errors.New("cannot get upgrade key")
	}
	err = RenderEnv(path.Join(wikiRoot, ".env"), "", secretKey, upgradeKey, "", "", "", false)
	if err != nil{
		return err
	}
	return nil
}

func RenderEnv(envFile,
	adminEmail,
	secretKey,
	upgradeKey,
	dbName,
	dbUser,
	dbPassword string,
	overwrite bool) error{
	envStore := make(map[string]string)

	envStore["EMERGENCY_CONTACT"] = adminEmail
	envStore["PASSWORD_SENDER"] = adminEmail
	envStore["SECRET_KEY"] = secretKey
	envStore["UPGRADE_KEY"] = upgradeKey
	envStore["DB_NAME"] = dbName
	envStore["DB_USER"] = dbUser
	envStore["DB_PASSWORD"] = dbPassword
	envStore["SERVER_HOST"] =

	if err, envExist := utils.FileExist(envFile); err == nil && envExist && !overwrite{
		envFi, _ := os.Open(envFile)
		defer envFi.Close()
		reader := bufio.NewReader(envFi)
		for true {
			bl, _, err := reader.ReadLine()
			if err == io.EOF{
				break
			}
			line := string(bl)
			pieces := strings.Split(line, "=")
			envName, envValue := pieces[0], pieces[1]
			envStore[envName] = envValue // add secret key upgrade key
			envStore["SECRET_KEY"] = secretKey
			envStore["UPGRADE_KEY"] = upgradeKey
		}
	}
	cntLine := make([]string, 0, 7)
	for key, item := range envStore{
		cntLine = append(cntLine, key+"="+item)
	}
	envContent := strings.Join(cntLine, "\n")
	envFi, err := os.OpenFile(envFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	defer envFi.Close()
	if err != nil{
		return err
	}
	_, err = envFi.Write([]byte(envContent))
	if err != nil{
		return err
	}
	return nil
}



