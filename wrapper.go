package main 

import (
	"os/exec"
	"os"
	"log"
    "github.com/getlantern/systray"
    "fmt"
    "time"
)
func onReady() {
	// systray.SetIcon(iconData)
	systray.SetTitle("GoNotes")
	systray.SetTooltip("Stickies")

	mNewWindow := systray.AddMenuItem("Open new sticky", "Open a new instance of a sticky note")
    systray.AddSeparator()
    mQuitOrig := systray.AddMenuItem("Quit", "Quit the whole app")
	go func() {
		<-mQuitOrig.ClickedCh
		fmt.Println("Requesting quit")
		systray.Quit()
		fmt.Println("Finished quitting")
	}()
	currentWorkingDirectory, errcwd := os.Getwd()
    if errcwd != nil {
        log.Fatal(errcwd)
    }

	go func() {
		for {
			select {
			case <-mNewWindow.ClickedCh:
                cmd := exec.Command("/usr/bin/gnome-terminal","--geometry", "80x10","--","bash","-c",currentWorkingDirectory+"/driver/gonotes_driver")
				if err := cmd.Start(); err != nil {
					fmt.Printf("Error starting command: %v\n", err)
				}
			}
		}
	}()
}

func main() {
    onExit := func() {
		now := time.Now()
		fmt.Printf(`We exited at %d.`, now.UnixNano())
	}
	systray.Run(onReady, onExit)
}
