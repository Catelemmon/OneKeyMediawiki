package wikiUtil

import "testing"

func TestInstallWiki(t *testing.T) {
	err := InstallWiki(
		"/home/cicada/workplaces/forGoSeries/src/github.com/Catelemmon/oneKeyMediawiki/statics/images",
		"/home/cicada/workplaces/forGoSeries/src/github.com/Catelemmon/oneKeyMediawiki/statics/scripts")
	if err != nil {
		t.Logf("Failed To Install Docker image %v \n", err)
	}
}