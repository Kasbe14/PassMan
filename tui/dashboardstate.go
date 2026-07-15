package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type DashBoardModel struct {
    tabs       []string
    activeTab  int
    quit       bool
}
type DashboardMsg struct {
}

func initialDModel() DashBoardModel{
    tabs := []string{"AddProfile", "OpenVault"}
    dbm := DashBoardModel{
        tabs: tabs,
        quit:false,
        activeTab: 1,
    }
    return dbm
}

func (dbm DashBoardModel)  Init() tea.Cmd {
    return nil
}

func (dbm DashBoardModel) Update(msg tea.Msg) (tea.Model,tea.Cmd){
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.String() {
        case "esc":
            dbm.quit = true
            return dbm, nil
        case "tab","left","right","l","d","a","shift+tab", "up", "down":
            //navigate
            dbm.activeTab = 1 - dbm.activeTab
        }
    }
    return dbm,nil
}

func (dbm DashBoardModel) View() string {
    title := titleStyle.Render(title)
    greet := titleStyle.Render("WELCOME !")
     var renderTabs []string
    for i,tabName := range dbm.tabs {
        if i == dbm.activeTab {
            renderTabs = append(renderTabs,activeTabStyle.Render(tabName))
        }else {
            renderTabs = append(renderTabs,inactiveTabStyle.Render(tabName))
        }
    }
 
    tabRow := lipgloss.JoinHorizontal(lipgloss.Top,renderTabs[0],"     ",renderTabs[1])
    finalUi := lipgloss.JoinVertical(lipgloss.Center,title,greet,tabRow,"\n",help)

    return lipgloss.Place(
        10,10, lipgloss.Center,lipgloss.Center,finalUi,
    )
        // greet := 
}

