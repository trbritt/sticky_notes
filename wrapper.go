package main 

import (
	"os/exec"
	"os"
	"log"
)
func main() {
	currentWorkingDirectory, errcwd := os.Getwd()
    if errcwd != nil {
        log.Fatal(errcwd)
    }
	cmd := exec.Command("/usr/bin/gnome-terminal", "--window-with-profile","Github","--","bash","-c",currentWorkingDirectory+"/driver/gonotes_driver")

    err := cmd.Run()

    if err != nil {
        log.Fatal(err)
    }
}
