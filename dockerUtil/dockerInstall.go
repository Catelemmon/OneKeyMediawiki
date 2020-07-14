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
	if err != nil{
		panic("check support env failed")
	}
	if si.HasDocker{
		return
	}
	fmt.Println("正在安装docker 守护程序")
	if ! utils.HasFile(execDir, "docker-19.03.9.tgz"){
		return errors.New("dockker package is not found")
	}
	// compress tar gz
	cmd := exec2.Command("tar", "-xf", "docker")
	_, err = cmd.Output()
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



	return nil
 }

func systemctlAdapt() error{

}

func serviceAdapt() error{

 }