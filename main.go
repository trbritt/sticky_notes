package main

import (
	"compress/gzip"
	"encoding/gob"
	"fmt"
	"os"
	"os/exec"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
)

const (
	initialInputs = 1
	maxInputs     = 6
	minInputs     = 1
	helpHeight    = 5
)
var (
	maxWidth      = 3
	maxHeight     = maxInputs / maxWidth
)

var (
	cursorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("212"))

	cursorLineStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("57")).
			Foreground(lipgloss.Color("230"))

	placeholderStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("238"))

	endOfBufferStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("235"))

	focusedPlaceholderStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("99"))

	focusedBorderStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("238"))

	blurredBorderStyle = lipgloss.NewStyle().
				Border(lipgloss.HiddenBorder())
)

type keymap = struct {
	next, prev, add, write, remove, quit key.Binding
}

func newTextarea() textarea.Model {
	t := textarea.New()
	t.Prompt = ""
	t.Placeholder = "Type something"
	t.ShowLineNumbers = true
	t.Cursor.Style = cursorStyle
	t.FocusedStyle.Placeholder = focusedPlaceholderStyle
	t.BlurredStyle.Placeholder = placeholderStyle
	t.FocusedStyle.CursorLine = cursorLineStyle
	t.FocusedStyle.Base = focusedBorderStyle
	t.BlurredStyle.Base = blurredBorderStyle
	t.FocusedStyle.EndOfBuffer = endOfBufferStyle
	t.BlurredStyle.EndOfBuffer = endOfBufferStyle
	t.KeyMap.DeleteWordBackward.SetEnabled(false)
	t.KeyMap.LineNext = key.NewBinding(key.WithKeys("down"))
	t.KeyMap.LinePrevious = key.NewBinding(key.WithKeys("up"))
	t.Blur()
	return t
}

type model struct {
	width  int
	height int
	keymap keymap
	help   help.Model
	inputs []textarea.Model
	focus  int
}
var fname = "/home/"+ os.Getenv("USER") +"/.cache/gonotes.gogz"

func newModel() model {
	m := model{
		inputs: make([]textarea.Model, initialInputs),
		help:   help.New(),
		keymap: keymap{
			next: key.NewBinding(
				key.WithKeys("tab"),
				key.WithHelp("tab", "next"),
			),
			prev: key.NewBinding(
				key.WithKeys("shift+tab"),
				key.WithHelp("shift+tab", "prev"),
			),
			add: key.NewBinding(
				key.WithKeys("ctrl+n"),
				key.WithHelp("ctrl+n", "add an editor"),
			),
			remove: key.NewBinding(
				key.WithKeys("ctrl+w"),
				key.WithHelp("ctrl+w", "remove an editor"),
			),
			quit: key.NewBinding(
				key.WithKeys("esc", "ctrl+c"),
				key.WithHelp("esc", "quit"),
			),
			write: key.NewBinding(
				key.WithKeys("ctrl+s"),
				key.WithHelp("ctrl+s", "write all stickies to file"),
			),
		},
	}
	for i := 0; i < initialInputs; i++ {
		m.inputs[i] = newTextarea()
	}
	// We need to see if this file isn't empty and
	// if theres stuff we need to put it into the stickies!
	if _, err := os.Stat(fname); os.IsNotExist(err) {
		// file does not exist
		println("File does not exist; creating clean notes :)")
	} else {
		// file exists
		// println("File exists")
		// read from the compressed binary file
		// open the gzip file for reading

		gzFile, err := os.Open(fname)
		if err != nil {
			panic(err)
		}
		defer gzFile.Close()

		// create a gzip reader
		gzReader, err := gzip.NewReader(gzFile)
		if err != nil {
			panic(err)
		}
		defer gzReader.Close()

		// create a binary decoder
		decoder := gob.NewDecoder(gzReader)

		// decode the slice of strings
		var decodedStrs []string
		err = decoder.Decode(&decodedStrs)
		if err != nil {
			panic(err)
		}

		// print the decoded slice of strings
		max_not_empty := 0
		for i:= range decodedStrs {
			if decodedStrs[i] != "" {
				if i > max_not_empty {
					max_not_empty = i
				}
				// fmt.Println(i, " is not empty!")
				// fmt.Println(decodedStrs[i])
			}
		}
		// fmt.Println(max_not_empty)
		for i:=0; i<max_not_empty; i++ {
			m.inputs = append(m.inputs, newTextarea())
		}
		for i:= range m.inputs {
			m.inputs[i].SetValue((decodedStrs[i]))
		}
	}
	
	m.inputs[m.focus].Focus()
	m.updateKeybindings()
	return m
}

func (m model) Init() tea.Cmd {
	return textarea.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keymap.quit):
			for i := range m.inputs {
				m.inputs[i].Blur()
			}
			return m, tea.Quit
		case key.Matches(msg, m.keymap.next):
			m.inputs[m.focus].Blur()
			m.focus++
			if m.focus > len(m.inputs)-1 {
				m.focus = 0
			}
			cmd := m.inputs[m.focus].Focus()
			cmds = append(cmds, cmd)
		case key.Matches(msg, m.keymap.prev):
			m.inputs[m.focus].Blur()
			m.focus--
			if m.focus < 0 {
				m.focus = len(m.inputs) - 1
			}
			cmd := m.inputs[m.focus].Focus()
			cmds = append(cmds, cmd)
		case key.Matches(msg, m.keymap.add):
			m.inputs = append(m.inputs, newTextarea())
		case key.Matches(msg, m.keymap.remove):
			m.inputs = m.inputs[:len(m.inputs)-1]
			if m.focus > len(m.inputs)-1 {
				m.focus = len(m.inputs) - 1
			}
		case key.Matches(msg, m.keymap.write):
			var total_contents []string 
			for i := range m.inputs {
				total_contents = append(total_contents, m.inputs[i].Value())
			}
			//write slice of strings to compressed binary
			file, err := os.Create(fname)
			if err != nil {
				panic(err)
			}
			defer file.Close()

			// create a gzip writer
			gzWriter := gzip.NewWriter(file)
			defer gzWriter.Close()

			// create a binary encoder
			encoder := gob.NewEncoder(gzWriter)

			// encode the slice of strings
			err = encoder.Encode(total_contents)
			if err != nil {
				panic(err)
			}

		}
	case tea.WindowSizeMsg:
		m.height = msg.Height
		m.width = msg.Width
		//check if terminal is portrait or landscape
		if (msg.Height > msg.Width){
			//portrait mode
			maxHeight = 3
			maxWidth = 2
		}	else {
			maxHeight = 2
			maxWidth = 3
		}
	}

	m.updateKeybindings()
	m.sizeInputs()

	// Update all textareas
	for i := range m.inputs {
		newModel, cmd := m.inputs[i].Update(msg)
		m.inputs[i] = newModel
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m *model) sizeInputs() {
	for i := range m.inputs {
		divisor := 0.0
		m.inputs[i].SetWidth((m.width) / maxWidth)
		if len(m.inputs) <= maxWidth {
			divisor = 1.1 //full height ish
		} else {
			divisor = 2.1
		}
		m.inputs[i].SetHeight(int(float64(m.height-helpHeight) / divisor))
	}
}

func (m *model) updateKeybindings() {
	m.keymap.add.SetEnabled(len(m.inputs) < maxInputs)
	m.keymap.remove.SetEnabled(len(m.inputs) > minInputs)
}

func (m model) View() string {
	help := m.help.ShortHelpView([]key.Binding{
		m.keymap.next,
		m.keymap.prev,
		m.keymap.add,
		m.keymap.remove,
		m.keymap.write,
		m.keymap.quit,
	})

	var viewsX []string
	var viewsY []string

	for i := range m.inputs {
		if i < maxWidth {
			viewsX = append(viewsX, m.inputs[i].View())
		} else {
			viewsY = append(viewsY, m.inputs[i].View())
		}
	}
	return lipgloss.JoinHorizontal(lipgloss.Top, viewsX...) + "\n" + lipgloss.JoinVertical(lipgloss.Bottom, lipgloss.JoinHorizontal(lipgloss.Top, viewsY...)) + "\n\n" + help
}

func main() {
	termenv.SetWindowTitle("gonotes")
	cmd := exec.Command("gnome-terminal")

	writer, err := cmd.StdinPipe()
	if err != nil {
		panic(err)
	}
	tea.WithOutput(writer)
	// start the new terminal window
	err = cmd.Start()
	if err != nil {
		panic(err)
	}
	if _, err := tea.NewProgram(newModel(), tea.WithAltScreen()).Run(); err != nil {
		fmt.Println("Error while running program:", err)
		os.Exit(1)
	}
}
