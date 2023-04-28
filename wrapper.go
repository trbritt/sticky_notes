package main 

import (
	"os/exec"
	"os"

    
	"log"
    "fmt"
)
func main() {
    out, err := exec.Command("/usr/bin/dconf", "read", "/org/gnome/terminal/legacy/profiles:/:54cf545b-ecbf-4289-b9dd-56f3381de31b/background-color").Output()

    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("The color is %s\n", out)
	currentWorkingDirectory, errcwd := os.Getwd()
    if errcwd != nil {
        log.Fatal(errcwd)
    }
	cmd := exec.Command("/usr/bin/gnome-terminal", "--window-with-profile","Default","--","bash","-c",currentWorkingDirectory+"/driver/gonotes_driver")

    err = cmd.Run()

    if err != nil {
        log.Fatal(err)
    }

}
