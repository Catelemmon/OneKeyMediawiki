package wikiUtil

import (
	"testing"
)

func TestInstallWiki(t *testing.T) {
	err := InstallWiki(
		"/home/cicada/workplaces/forGoSeries/src/github.com/Catelemmon/oneKeyMediawiki/statics/images",
		"/home/cicada/workplaces/forGoSeries/src/github.com/Catelemmon/oneKeyMediawiki/statics/scripts")
	if err != nil {
		t.Logf("Failed To Install Docker image %v \n", err)
	}
}

func TestRenderEnv(t *testing.T) {
	err := RenderEnv("./.env", "1713856662a@gmail.com", "", "", "dbname", "cicada", "09170725", true)
	if err != nil{
		t.Logf("Failed to Render Env file")
	}
	err = RenderEnv("./.env", "", "dlshdlkshd", "dsdsdsfgdfsgfd", "", "", "", false)
	if err != nil {
		t.Logf("failed to Render Env file")
	}
}