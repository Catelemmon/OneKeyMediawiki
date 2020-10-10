package dockerUtil

import (
	"errors"
	"fmt"
	"github.com/Catelemmon/oneKeyMediawiki/checker"
	"github.com/Catelemmon/oneKeyMediawiki/utils"
	"io/ioutil"
	"os"
	exec2 "os/exec"
	"path"
	"sync"
)


func InstallDocker(execDir string) error{
	si := checker.NewSupportInfo()
	err := si.GetSupportInfo()
	if err != si.GetSupportInfo(){
		panic("check support env failed")
	}
	if (si.SystemVendor == "centos" || si.SystemVendor == "rhel") && si.SystemVersion < "7"  {
		return InstallDockerOldCentos(execDir)
	} else {
		panic("未做centos7以上系统的兼容...\n")
	}
}

func InstallDockerOldCentos(execDir string) error{
	fmt.Println("正在安装docker守护程序")
	if !utils.HasFile(execDir, "docker-engine-1.7.1-1.el6.x86_64.rpm"){
		return errors.New("docker package is not found")
	}
	cmd := exec2.Command("rpm", "-ivh", "docker-engine-1.7.1-1.el6.x86_64.rpm")
	cmd.Dir = execDir
	out, err := cmd.Output()
	if err != nil{
		fmt.Println("failed to install docker when install docker")
		fmt.Printf("install docker output: %s\n", out)
		return err
	}
	err = InstallDockerComposeOldCentos(execDir)
	if err != nil{
		return err
	}
	// launch docker
	cmd = exec2.Command("service", "docker", "start")
	if err != nil {
		return err
	}
	_, err = cmd.Output()
	if err != nil{
		return err
	}
	fmt.Println("docker 安装启动完毕")
	return nil
}

func InstallDockerComposeOldCentos(execDir string) error{
	// install docker compose
	fs, err := ioutil.ReadFile(path.Join(execDir, "docker-compose-Linux-x86_64-1.5.1"))
	if err != nil {
		fmt.Println("failed to install docker-compose when reading docker-compose")
		return err
	}
	err = ioutil.WriteFile(path.Join("/usr/bin", "docker-compose"), fs, 0755)
	if err != nil{
		fmt.Println("failed to install docker-compose when writing /usr/bin")
		return  err
	}
	return nil
}

func InstallDockerCommon(execDir string, si *checker.SupportInfo) error{
	if si.HasDocker{
		return nil
	}
	fmt.Println("正在安装docker 守护程序")
	if ! utils.HasFile(execDir, "docker-19.03.9.tgz"){
		return errors.New("docker package is not found")
	}
	// compress tar gz
	cmd := exec2.Command("tar", "-xf", "docker")
	_, err := cmd.Output()
	if err != nil{
		return err
	}

	//travel docker dir
	dockerCmps, err := ioutil.ReadDir(path.Join(execDir, "docker"))
	if err != nil {
		return err
	}
	tgtFile := make([]string, 0, 10)
	var wg sync.WaitGroup
	failedFlag := make(chan struct{}, len(dockerCmps))
	for _, cmp := range dockerCmps{
		wg.Add(1)
		go func(filename string, _wg *sync.WaitGroup) {
			// execute copy
			defer _wg.Done()
			filepath := path.Join(execDir, filename)
			fs, err := ioutil.ReadFile(filepath)
			if err != nil{
				failedFlag <- struct{}{}
				return
			}
			err = ioutil.WriteFile("/usr/bin/"+filename, fs, 0751)
			if err != nil {
				failedFlag <- struct{}{}
				return
			}
			tgtFile = append(tgtFile, filepath)
		}(cmp.Name(), &wg)
	}
	wg.Wait()
	close(failedFlag)
	if len(failedFlag) != 0{
		// copy failed execute backup
		for _, src := range tgtFile{
			_ = os.Remove(src)

		}
		return errors.New("copy docker files failed")
	}

	if si.HasSystemd {
		err := systemctlAdapt(execDir)
		if err !=  nil {
			return  err
		}
	}

	return nil
 }

func systemctlAdapt(execDir string) error{
	if _, err := os.Stat("/etc/systemd/system"); os.IsNotExist(err){
		os.Mkdir("/etc/systemd/system", 0751)
	}
	if _, err := os.Stat(path.Join("/etc/systemd/system", "docker.service")); !os.IsNotExist(err){
		return nil
	}

	fs, err := ioutil.ReadFile(path.Join(execDir, "docker.service"))
	if err != nil{
		return err
	}
	err = ioutil.WriteFile(path.Join("/etc/systemd/system", "docker.service"), fs, 0644)
	if err != nil{
		return err
	}
	return nil
}

func serviceAdapt(execDir string) error{
	fs, err := ioutil.ReadFile(path.Join(execDir, "docker.script"))
	if err != nil{
		return err
	}
	err = ioutil.WriteFile(path.Join("/etc/init.d", "docker"), fs, 0751)
	if err != nil{
		return err
	}
	return nil
}