package kernalUtil

import (
	"fmt"
	"github.com/Catelemmon/oneKeyMediawiki/checker"
	"io/ioutil"
	"os"
	exec2 "os/exec"
	"regexp"
	"strings"
)

var kernelReg *regexp.Regexp = regexp.MustCompile(`(?m)vmlinuz-(([1-9]\d|[1-9])(\.(([1-9]\d{1,2})|\d)){1,2})`)

func CheckKernelPack(wd string) (bool, error){
	files, err := ioutil.ReadDir(wd)
	if err != nil{
		return false, err
	}
	rpmCnt := 0
	for _, f := range files{
		if f.Name() == "kernel-lt-devel-4.4.217-1.el6.elrepo.x86_64.rpm" ||
			f.Name() == "kernel-lt-4.4.217-1.el6.elrepo.x86_64.rpm"{
			rpmCnt += 1
		}
	}
	if rpmCnt == 2{
		return true, nil
	} else{
		return false, nil
	}
}

func InstallKernel(execDir string){

	si := checker.NewSupportInfo()
	err := si.GetSupportInfo()
	if err != nil {
		panic("check support support env failedly")
	}
	if si.KernelAllow{
		return
	}
	var updatekernel string
	fmt.Printf("是否需要更新内核文件到4.4.217[Y/n]:")
	fmt.Scanln(updatekernel)
	if strings.TrimSpace(updatekernel) == "n"{
		os.Exit(1)
	}
	cmd := exec2.Command("bash","-c", " sudo rpm -ivh kernel-lt-devel-4.4.217-1.el6.elrepo.x86_64.rpm")
	cmd.Dir = execDir
	out, err := cmd.Output()
	if err != nil {
		fmt.Printf("内核更新失败， 错误信息如下!\n%v", string(out))
		os.Exit(3)
	}
	cmd = exec2.Command("bash", "-c", "sudo rpm -ivh kernel-lt-4.4.217-1.el6.elrepo.x86_64.rpm")
	cmd.Dir = execDir
	_, err = cmd.Output()
	if err != nil {
		fmt.Printf("kernel install failed, output info as follow!\n%v", string(out))
		os.Exit(3)
	}
	fmt.Println("更新内核完毕")
}

func ChangeLaunchOption() error{
	grubConf := "/etc/grub.cfg"
	fmt.Println("正在修改内核启动项")
	_, err := os.Stat(grubConf)
	if err != nil || os.IsNotExist(err){
		return err
	}
	cnt, err := ioutil.ReadFile(grubConf)
	if err != nil{
		return nil
	}
	ms :=kernelReg.FindAllStringSubmatch(string(cnt), -1)
	var defVal int8
	for i, sms := range ms{
		if sms[1] == "4.4.217"{
			defVal = int8(i)
		}
	}
	subDef := regexp.MustCompile(`(?m)^default=\d+?$`)
	tgtS := subDef.ReplaceAllString(string(cnt), fmt.Sprintf("default=%d", defVal))
	err = ioutil.WriteFile(grubConf, []byte(tgtS), 0644)
	if err != nil{
		if err == os.ErrPermission{
			fmt.Println("建议使用sudo权限执行本程序")
			os.Exit(3)
		}
		return err
	}
	fmt.Println("修改内核启动项完毕，请重启...")
	return nil
}

