package utils

import "testing"

func TestVersionAllow(t *testing.T) {
	t.Log(VersionAllow("3.6.1.6", "3.6.2"))
}

func TestVersionExtract(t *testing.T) {
	t.Log(VersionExtract("5.3.0-61-generic"))
}

func TestHasCommand(t *testing.T) {
	t.Log(HasCommand("iptables"))
	t.Log(HasCommand("fakecmd"))
}

func TestCommandVersion(t *testing.T) {
	t.Log(CommandVersion("iptables"))
	t.Log(CommandVersion("gdb"))
}

func TestHasFile(t *testing.T) {
	t.Log(HasFile("/usr/bin", "docker"))
}