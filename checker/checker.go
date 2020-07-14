package checker

import (
	"fmt"
	"github.com/Catelemmon/oneKeyMediawiki/utils"
	"github.com/zcalusic/sysinfo"
	"log"
	"os"
	"os/user"
	"reflect"
)

type SupportInfo struct {
	hasChecked    bool
	KernelAllow   bool `comment:"内核是否允许"`
	IpTablesAllow bool `comment:"iptable版本是否支持"`
	HasApache     bool `comment:"是否有apache服务器"`
	HasNginx      bool `comment:"是否有nginx服务器"`
	HasDocker     bool `comment:"是否有docker"`
	GrubOption    string `comment:"启动项信息"`
	HasGit        bool `comment:"是否有git"`
	GitAllow      bool `comment:"git版本是否符合需求"`
	HasXz         bool `comment:"是否有xz"`
	XzAllow       bool `comment:"xz程序是否允许"`
	HasMysql      bool `comment:"是否有mysql数据库软件"`
	MysqlAllow    bool `comment:"mysql版本是否符合要求"`
	HasPHP        bool `comment:"是否有php"`
	PHPAllow      bool `comment:"php版本是否符合要求"`
	HasNode       bool `comment:"是否有nodejs"`
	PackageAllow  bool `comment:"其他一系列包的支持情况"`
}

var supportInfo *SupportInfo

func  NewSupportInfo() *SupportInfo{
	if supportInfo == nil{
		supportInfo = &SupportInfo{
			KernelAllow:   false,
			IpTablesAllow: false,
			HasApache:     false,
			HasNginx:      false,
			HasDocker:     false,
			GrubOption:    "",
			HasGit:        false,
			GitAllow:      false,
			HasXz:         false,
			XzAllow:       false,
		}
	}
	return supportInfo
}

func getSystemInfoCheck() sysinfo.SysInfo{
	currentUser, err := user.Current()
	if err != nil {
		log.Fatalf("%v", err)
	}
	if currentUser.Uid != "0"{
		fmt.Println("Please use root account!")
		os.Exit(3)
	}
	var systemInfo sysinfo.SysInfo
	systemInfo.GetSysInfo()

	return systemInfo
}

func (si *SupportInfo) GetSupportInfo() error {
	if (*supportInfo).hasChecked{
		return nil
	}
	currSysInfo := getSystemInfoCheck()
	var err error
	err = checkKernel(currSysInfo)
	if err != nil {
		return err
	}
	err = checkIpTables()
	if err != nil {
		return err
	}
	err = checkGit()
	if err != nil {
		return err
	}
	err = checkApache()
	if err != nil {
		return err
	}
	err = checkNginx()
	if err != nil {
		return err
	}
	err = checkXz()
	if err != nil{
		return err
	}
	err = checkPHP()
	if err != nil{
		return nil
	}
	err = checkHasNode()
	if err != nil {
		return nil
	}
	err = checkPackage()
	if err != nil{
		return nil
	}
	err = checkMysql()
	if err != nil {
		return nil
	}
	return nil
}

func (si *SupportInfo) ShowSupportInfo(){
	fmt.Println("系统支持信息如下：")
	waitChan := make(chan struct{})
	go func(travelItem interface{}){
		typ := reflect.TypeOf(travelItem)
		val := reflect.ValueOf(travelItem)
		num := val.NumField()
		for i := 0 ; i < num; i++ {
			var outInfo string
			if itemKnd := val.Field(i).Kind(); itemKnd == reflect.Bool && val.Field(i).Bool(){
				outInfo = "是"
			} else if itemKnd == reflect.Bool && !val.Field(i).Bool(){
				outInfo = "否"
			} else if itemKnd == reflect.String{
				continue
			}
			tag := typ.Field(i).Tag.Get("comment")
			if tag != ""{
				fmt.Printf("%v: %s\n", typ.Field(i).Tag.Get("comment"), outInfo)
			}

		}
		waitChan <- struct{}{}
	}(*si)
	<- waitChan
}

func checkKernel(info sysinfo.SysInfo) error{

	kernelString := info.Kernel.Release
	version, err := utils.VersionExtract(kernelString)
	if  err != nil {
		(*supportInfo).KernelAllow = false
	}
	if utils.VersionAllow(version, "3.10"){
		(*supportInfo).KernelAllow = true
	} else{
		(*supportInfo).KernelAllow = false
	}
	return nil
}

func checkIpTables() error{

	ver, err := utils.CommandVersion("iptables")
	if err != nil {
		panic("unsupported iptables")
	}
	if utils.VersionAllow(ver, "1.4"){
		(*supportInfo).IpTablesAllow = true
	} else {
		(*supportInfo).IpTablesAllow = false
	}
	return nil
}

func checkGit() error{
	if utils.HasCommand("git") {
		(*supportInfo).HasGit = true
	} else{
		(*supportInfo).HasGit = false
		(*supportInfo).GitAllow = false
		return nil
	}
	ver, err := utils.CommandVersion("git")
	if err != nil {
		panic("unsupported git")
	}
	if utils.VersionAllow(ver, "1.7") {
		(*supportInfo).GitAllow = true
	} else {
		(*supportInfo).GitAllow = false
	}
	return nil
}

func checkXz() error {

	if utils.HasCommand("xz"){
		(*supportInfo).HasXz = true
	} else {
		(*supportInfo).HasXz = false
		(*supportInfo).XzAllow = false
	}
	ver, err := utils.CommandVersion("xz")
	if err != nil {
		panic("unsupported xz")
	}
	if utils.VersionAllow(ver, "4.9"){
		(*supportInfo).XzAllow = true
	} else {
		(*supportInfo).XzAllow = false
	}
	return nil
}

func checkNginx() error{
	if utils.HasCommand("nginx"){
		(*supportInfo).HasNginx = true
	} else {
		(*supportInfo).HasNginx =false
	}
	return nil
}

func checkApache() error{
	if utils.HasCommand("httpd"){
		(*supportInfo).HasApache = true
	} else {
		(*supportInfo).HasApache = false
	}
	return nil
}

func checkPHP() error{
	if utils.HasCommand("php"){
		(*supportInfo).HasPHP = true
	} else {
		(*supportInfo).HasPHP = false
	}
	ver, err := utils.CommandVersion("php")
	if err != nil{
		fmt.Println("php version is not supported")
	}
	if utils.VersionAllow(ver, "7.2.9"){
		(*supportInfo).HasPHP = true
	} else {
		(*supportInfo).XzAllow = false
	}
	return nil
}

func checkHasNode() error{
	if utils.HasCommand("node") {
		(*supportInfo).HasNode = true
	} else {
		(*supportInfo).HasNode = false
	}
	return nil
}

func checkMysql() error{
	if utils.HasCommand("mysqld"){
		(*supportInfo).HasMysql = true
	} else {
		(*supportInfo).HasMysql = false
	}
	ver, err := utils.CommandVersion( "mysqld")
	if err != nil {
		fmt.Println("mysql version is not supported")
		return err
	}
	if utils.VersionAllow(ver, "5.5.8"){
		(*supportInfo).MysqlAllow = true
	} else {
		(*supportInfo).MysqlAllow = false
	}
	return nil
}

func checkPackage() error{
	// TODO：add packages check
	(*supportInfo).PackageAllow = false
	return nil
}
