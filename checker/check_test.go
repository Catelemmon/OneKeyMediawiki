package checker

import "testing"


func TestKernelCheck(t *testing.T) {
	t.Logf("System Info as Follow:\n %v \n ", getSystemInfoCheck())
}

func TestCheckKernel(t *testing.T){
	sysInfo := getSystemInfoCheck()
	err := checkSysMeta(sysInfo)
	if err != nil{
		t.Fatal(err)
	}
}

func TestCheckSupport(t *testing.T) {
	si := NewSupportInfo()
	err := si.GetSupportInfo()
	if err != nil{
		t.Log(err)
	}
	t.Log(si)
}

func TestSupportInfo_ShowSupportInfo(t *testing.T) {
	si := NewSupportInfo()
	err := si.GetSupportInfo()
	if err != nil{
		t.Log(err)
	}
	si.ShowSupportInfo()
}