package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"time"
	"github.com/getlantern/systray"
)

var accumulator int = 0

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
		panic(errcwd)
	}

	go func() {
		for {
			select {
			case <-mNewWindow.ClickedCh:
				fmt.Printf("Accumulator: %d\n", accumulator)
				var terminal_cmd [6]string
				terminal_cmd[0] = "/usr/bin/gnome-terminal"
				terminal_cmd[1] = "--geometry"
				terminal_cmd[2] = "90x30"
				terminal_cmd[3] = "--"
				terminal_cmd[4] = currentWorkingDirectory+"/driver/gonotes_driver"
				terminal_cmd[5] = "-id="+strconv.Itoa(accumulator)

                cmd := exec.Command(terminal_cmd[0], terminal_cmd[1:]...)
				fmt.Printf("Command: %v\n", cmd.String())

				err := cmd.Start()
				if err != nil {
					fmt.Printf("Error starting command: %v\n", err)
				}

                err = cmd.Process.Release()
                if err != nil {
                    fmt.Println("cmd.Process.Release failed: ", err)
                }
				accumulator++
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
