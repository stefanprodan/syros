package main

import (
	"github.com/codeskyblue/go-sh"
	"os"
)

func main() {
	wd, _ := os.Getwd()
	session := sh.NewSession()
	session.SetEnv("BUILD_ID", "123")
	session.SetDir(wd)
	session.Command("ls", "-la").Run()
	session.ShowCMD = true
}
