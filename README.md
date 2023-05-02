# GoNotes

Welcome to GoNotes! This is a nice terminal application written in `Go` that utilizes the almighty `gnome-terminal` to populate
multiple terminals that act as digital 'sticky notes'. Each sticky note can be saved to and loaded from disk for easy retrieval (deletion of previous data
in progress). It will open a manager accessible with the rest of the system tray icons that will allow you to load stickies as necessary, or stop the program if you want.

## Prerequisites

This software can only be run on Linux (sorry not sorry!). To install the prerequisites, you can follow the following commands for Debian-like systems:

```bash
sudo apt-get install golang wmctrl gcc libgtk-3-dev libayatana-appindicator3-dev
```

## Installation

To run the software, simply run an execute the `make.sh` script in the root directory. It will compile and execute all the needed commands for you.

There's also a desktop launcher should you wish to copy it to your desktop.

#### Nota bene: these scripts have hard-coded absolute paths; adjust as necessary!

