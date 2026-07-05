package tui

import (
	"github.com/Kasbe14/PassMan/core"
	tea "github.com/charmbracelet/bubbletea"
    // "fmt"
    // "strings"
)

type sessionState int

const (
    stateStart sessionState = iota + 1
    stateSignUp
    stateLogin
    stateVault
)
 

type MainModel struct {
    vault *core.VaultService
    state sessionState
    sm     StartUpModel
    su     SignUpModel
    text string
}

func (m MainModel) Init() tea.Cmd {

    return tea.Batch(
        m.sm.Init(),
        m.su.Init(),
    )
}

func (m MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Global overrides (Ctrl+C to quit)
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		if keyMsg.Type == tea.KeyCtrlC /*|| keyMsg.String() == "q"*/{
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd

	switch m.state {
	case stateStart:
		// Intercept the "Enter" key BEFORE handing it down
		if keyMsg, ok := msg.(tea.KeyMsg); ok && keyMsg.String() == "enter" {
            
			if m.sm.activeTab == 0 {
				m.state = stateSignUp // Tab 0 is SignUp
                m.su = initialSModel()
                cmd = m.su.Init()
                return m,cmd
			} else if m.sm.activeTab == 1 {
				m.state = stateLogin  // Tab 1 is LogIn
			}
			return m, cmd 
        }

        //continue nav on startpage
		updatedModel, subCmd := m.sm.Update(msg)
		m.sm = updatedModel.(StartUpModel)
		cmd = subCmd

	case stateSignUp:
        //TODO
        if keyMsg, ok := msg.(tea.KeyMsg); ok && (keyMsg.String() == "q" && m.su.quitSignUp == true) {
            m.state = stateStart
            return m,nil
        }
        // m.su = initialSModel()
        // m.su.Init()
        updatedModel, subCmd := m.su.Update(msg)
        m.su = updatedModel.(SignUpModel)
        cmd = subCmd
	case stateLogin:
        //TODO
	}

	return m, cmd
}

func InitialMainModel(sm StartUpModel) MainModel {
    return MainModel{
        sm: sm,
        state:stateStart,
    }
}

func (m MainModel) View() string {
    switch m.state {
    case stateStart:
        return m.sm.View()
    case stateSignUp:
        //return signyp
        return m.su.View()
    }
    return "Loading.."
}

//tui entry point
func Start()  {
    p := tea.NewProgram(

        InitialMainModel(StartUpModel{Tabs: []string{"SignUp","LogIn"},ProgramTitle: "PassMan",activeTab:0}),
    )
    if _,err := p.Run(); err != nil {
        panic(err)
    }
}
