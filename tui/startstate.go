package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)
var (
	//program name 
	titleStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#9FA1FF")). // Deep Pastel Purple
		MarginBottom(2)

	//highligh active tab and mute inactive tab 
	activeTabStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#AEE2FF")). // Icy Blue Border
		Foreground(lipgloss.Color("#D9F9DF")).       // Mint Green Text
		Bold(true).
		Padding(0, 3)

	// The Unselected "Inactive" Tab
	inactiveTabStyle = lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")). 
		Foreground(lipgloss.Color("240")).
		Padding(0, 3)
)



type StartUpModel struct {
    ProgramTitle string
    Tabs []string
    activeTab int
}

func (sm StartUpModel) Init() tea.Cmd {
    return nil
}
func (sm StartUpModel) Update(msg tea.Msg) (tea.Model,tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch keypress := msg.String(); keypress {
        case "ctrl+c", "q":
            return sm,tea.Quit

        case "l", "right", "tab","d":
            sm.activeTab =  1 - sm.activeTab
            return sm, nil
        case "r","left","shift+tab","a":
            sm.activeTab = 1 - sm.activeTab 
            return sm,nil

        }

    }
    return  sm,nil
}
func (sm StartUpModel) View() string {
    title := titleStyle.Render(sm.ProgramTitle)
    var renderTabs []string
    for i,tabName := range sm.Tabs {
        if i == sm.activeTab {
            renderTabs = append(renderTabs,activeTabStyle.Render(tabName))
        }else {
            renderTabs = append(renderTabs,inactiveTabStyle.Render(tabName))
        }
    }
    tabRow := lipgloss.JoinHorizontal(lipgloss.Top,renderTabs[0],"     ",renderTabs[1])
    finalUi := lipgloss.JoinVertical(lipgloss.Center,title,tabRow)

    return lipgloss.Place(
        10,10, lipgloss.Center,lipgloss.Center,finalUi,
    )
}
