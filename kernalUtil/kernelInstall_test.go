package kernalUtil

import "testing"

func TestCheckKernelPack(t *testing.T) {
	t.Log(CheckKernelPack("/home/cicada/workplaces/forGoSeries/src/github.com/" +
		"Catelemmon/oneKeyMediawiki/statics/kernel"))
}

func TestChangeLaunchOption(t *testing.T) {
	t.Log(ChangeLaunchOption())
}
