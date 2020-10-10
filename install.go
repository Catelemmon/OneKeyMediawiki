package main

import (
	"fmt"
	"github.com/Catelemmon/oneKeyMediawiki/checker"
	"github.com/Catelemmon/oneKeyMediawiki/dockerUtil"
	"github.com/Catelemmon/oneKeyMediawiki/kernelUtil"
	"github.com/Catelemmon/oneKeyMediawiki/wikiUtil"
	"os"
	"path"
	"strings"
)

func main() {

	si := checker.NewSupportInfo()
	err := si.GetSupportInfo()
	if err != nil{
		fmt.Println("获取基础信息失败...")
		fmt.Printf("%v", err)
		os.Exit(3 )
	}
	si.ShowSupportInfo()
	cwd, err := os.Getwd()
	if err != nil{
		fmt.Println("获取当前文件夹失败...")
	}
	staticsDir := path.Join(cwd, "statics")
	if si.KernelAllow{
		// installer docker directly
		if !si.HasDocker{
			err := dockerUtil.InstallDocker(path.Join(staticsDir, "docker"))
			if err != nil{
				fmt.Printf("安装docker失败，错误信息如下\n%v", err)
				os.Exit(3)
			}
		}
		// reinstall docker-compose
		fmt.Println("是否需要修复docker-compose[Y/n]")
		var repairDockerCompose string
		fmt.Scanln(&repairDockerCompose)
		if strings.TrimSpace(repairDockerCompose) == "Y" || strings.TrimSpace(repairDockerCompose) == "y" ||
			strings.TrimSpace(repairDockerCompose) == "" {
			err  = dockerUtil.InstallDockerComposeOldCentos(path.Join(staticsDir, "docker"))
			if err != nil{
				fmt.Printf("修复docker-compose失败，错误信息如下\n%v", err)
				os.Exit(3)
			}
		}
		fmt.Println("正在装载wiki镜像和启动中....")
		err = wikiUtil.InstallWiki(path.Join(staticsDir, "images"), path.Join(staticsDir, "scripts"))
		if err != nil {
			fmt.Printf("安装wiki镜像失败，错误信息如下\n%v", err)
			os.Exit(3)
		}
	} else{
		// kernel not allow
		fmt.Println("是否需要升级内核以安装Mediawiki?[Y/n]")
		var updateKernel string
		fmt.Scanln(&updateKernel)
		if strings.TrimSpace(updateKernel) == "n" {
			os.Exit(3)
		}

		fmt.Println("updating kernel...")
		kernelUtil.InstallKernel(path.Join(staticsDir, "kernel"))
		fmt.Println("修改内核启动项...")
		err := kernelUtil.ChangeLaunchOption()
		if err != nil{
			fmt.Printf("修改内核启动项失败, 错误信息如下\n%s", err)
		}
		fmt.Println("完成内核程序的安装，请重新启动")
	}
}


