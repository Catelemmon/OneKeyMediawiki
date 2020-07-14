package main

import (
	"fmt"
	"github.com/Catelemmon/oneKeyMediawiki/checker"
	"math/rand"
	"os"
	"strings"
	"time"
)

func main() {

	si := checker.NewSupportInfo()
	si.GetSupportInfo()
	si.ShowSupportInfo()
	fmt.Println("是否需要升级内核以安装Mediawiki?[Y/n]")
	var updateKernel string
	fmt.Scanln(&updateKernel)
	if strings.TrimSpace(updateKernel) == "n" {
		os.Exit(3)
	}
	fmt.Println("updating kernel...")
	time.Sleep(time.Duration(rand.Intn(10)) * time.Second)
	fmt.Println("不支持目标内核")
}


