package tui

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
) 

var (

)

type LoginModel struct {
    title string
    inputs  []textinput.Model
    quitLogin bool
}


func (lm LoginModel) Init() tea.Cmd {

    return textinput.Blink

}
func (lm LoginModel) Update(msg tea.Msg) (tea.Model,tea.Cmd) {
    return lm,nil
}
func (lm LoginModel) View() string {

    return "hello form Login model"

}
