package tui

import (
	// "github.com/Kasbe14/PassMan/model"
	"github.com/Kasbe14/PassMan/core"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)
type profileName struct {
    name string
}

func (p profileName) Title() string {return p.name}
func (p profileName) Description() string {return "Encrypted Vault"}
func (p profileName) FilterValue() string {return p.name}

type RevealPassMsg struct {
     pass string
}
type  HideMsg struct {}
type  AllProfilesMsg struct {}

//TODO swap lock for key bytes

type VaultModel struct {
    listProfile   list.Model
    key []byte
    userId int64
    profiles []string //lazy loading only get pass from db when asked
    revealing bool //msg to show password
    quit   bool
    visiblePass []byte
    
}

func (vm VaultModel) Init() tea.Cmd {
    return textinput.Blink
}
func initialVModel(key []byte, userID int64) VaultModel {
    vm := VaultModel{
        key: key,
        userId: userID,
    }
    return vm
}
func (vm VaultModel) Update(msg tea.Msg) (tea.Model,tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.String(){
        case "esc":
            vm.quit = true
            core.Wipe(vm.key)
            return vm,nil
        case "c":
            //TODO copy the password of this profile to clipboard and clear clipboard after 10 sec
        case "v":
             //TODO let user view passowrd for 20sec and hide it
        case "ctrl+d":
             //TODO delete the focused profile's data
        }
    }
    var cmd tea.Cmd
    //case blinking
    return vm,cmd
}
func (vm VaultModel) View() string {
    
    title := titleStyle.Render(title)

    var s string
    h := helpStyle.Render("[To Search: / |Copy: c | View: v | Delete: ctrl+d | Quit: esc ]")
    s = lipgloss.JoinVertical(lipgloss.Center,title,vm.listProfile.View(),h)


    return lipgloss.Place(10,10,lipgloss.Center,lipgloss.Center,s)
}
