package tui

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)


type AddProfileModel struct {
    addInput []textinput.Model
    activeInput int
    submit bool
    quit bool
}
 
type AddProMsg struct {
    Err string
}


func (apm  AddProfileModel) Init() tea.Cmd {
    return nil
}

func (apm  AddProfileModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.String(){
        case "esc":
            apm.quit = true
            return apm,nil
        }
    }

    return apm, nil
}

func (apm  AddProfileModel) View() string {
    return "hello from add profile"
}
