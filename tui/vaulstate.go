package tui

import (
	// "github.com/Kasbe14/PassMan/model"
	"time"

	"github.com/Kasbe14/PassMan/core"
	"github.com/charmbracelet/bubbles/list"

	// "github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)
type profileName struct {
    name string
    isRevealed bool
    decryptedPass []byte
}

func (p profileName) Title() string {return p.name}
//rendering pass
func (p profileName) Description() string {
    if(p.isRevealed)  {
        return string(p.decryptedPass)
    }
    return "••••••••••••••••"
}
func (p profileName) FilterValue() string {return p.name}

type RevealPassMsg struct {
     err     error
     pass    []byte
     proName string
     lockMsg string
}
type  HideMsg struct {
    proName string
}
type ViewPassMsg struct {
    profileName string
    // idx      int
}
type  AllProfilesMsg struct {}

//TODO swap lock for key bytes

type VaultModel struct {
    listProfile   list.Model
    key []byte
    userId int64
    profiles []profileName //lazy loading only get pass from db when asked
    revealing bool //msg to show password
    quit   bool
    visiblePass []byte
    
}

func (vm VaultModel) Init() tea.Cmd {
    return nil
}
func initialVModel(key []byte, userID int64) VaultModel {
dummyProfiles := []profileName{
	{
		name:          "Github",
		isRevealed:    false,
		decryptedPass: []byte("git_SuperSecret123!"),
	},
	{
		name:          "Google Account",
		isRevealed:    false,
		decryptedPass: []byte("goog_P@ssw0rd99"),
	},
	{
		name:          "ProtonMail",
		isRevealed:    false,
		decryptedPass: []byte("proton_Encrypted007"),
	},
	{
		name:          "HDFC Bank",
		isRevealed:    false,
		decryptedPass: []byte("bank_SecureMoney$$$"),
	},
	{
		name:          "Dattaniddhi Root",
		isRevealed:    false,
		decryptedPass: []byte("local_db_admin_2026"),
	},
}
	var items []list.Item
	for _, p := range dummyProfiles {
        items = append(items, p)
	}
    delegate := list.NewDefaultDelegate()
    //for rendering password
    delegate.ShowDescription = true
    delegate.Styles.NormalTitle = delegate.Styles.NormalTitle.
		Foreground(lipgloss.Color("242"))

	delegate.Styles.SelectedTitle = delegate.Styles.SelectedTitle.
		Foreground(lipgloss.Color("#D9F9DF")).            // The text color (Hot Pink)
		BorderLeftForeground(lipgloss.Color("#D9F9DF"))   // The left border/cursor color 
	delegate.Styles.SelectedDesc = delegate.Styles.SelectedTitle.
		Foreground(lipgloss.Color("#9FA1FF")).            
		BorderLeftForeground(lipgloss.Color("#D9F9DF"))   
    l := list.New(items,delegate,50,10)
    l.SetShowTitle(false)
	l.SetShowHelp(false)
	l.SetShowStatusBar(false)
    l.FilterInput.Prompt = "Search Vault: "
    l.FilterInput.Placeholder = "Profile name..."
    l.FilterInput.Cursor.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("#D9F9DF"))
    l.FilterInput.PlaceholderStyle = blurredStyle
    l.FilterInput.PromptStyle = focusedStyle
    l.FilterInput.TextStyle = textStyle
    vm := VaultModel{
        listProfile: l,
        key: key,
        userId: userID,
        profiles: dummyProfiles,
        revealing: false,
        quit: false,
    }
    return vm
}
func (vm VaultModel) Update(msg tea.Msg) (tea.Model,tea.Cmd) {
    switch msg := msg.(type) {
    case RevealPassMsg:
        if msg.err != nil {
            // TODO do somthing of error rendering
            return vm,nil
        }
        vm.revealing = true
        vm.visiblePass = msg.pass
        var renderBytes []byte //decide to show password or lock msg
        if msg.lockMsg == "" {
            renderBytes = msg.pass
        } else {
            renderBytes = []byte(msg.lockMsg)
        }
        //find selected profile by name
        var idx int
        var targetProfile profileName
        for i, item := range vm.listProfile.Items() {
            pro := item.(profileName) 
            if pro.name == msg.proName {
                idx = i
                targetProfile = pro
                break
            }
        }
        //modify the profile in the list
        targetProfile.isRevealed = true
        targetProfile.decryptedPass = renderBytes
        //updating the listmodel
        updateCmd := vm.listProfile.SetItem(idx,targetProfile)
        //timer for hiding the message passess hidemsg
        timerCmd := tea.Tick(time.Second*10, func (t time.Time) tea.Msg {
            return HideMsg{proName: msg.proName}
        })
        return vm, tea.Batch(updateCmd,timerCmd)
    case HideMsg:
        if !vm.revealing {return vm,nil}
        core.Wipe(vm.visiblePass)
        vm.visiblePass = nil
        vm.revealing = false
        //hiding the rendered pass in list
        var idx int
        var targetProfile profileName
        for i, item := range vm.listProfile.Items() {
            pro := item.(profileName)
            if pro.name == msg.proName {
                idx = i
                targetProfile = pro
            }
        }
        targetProfile.isRevealed = false
        targetProfile.decryptedPass = nil
        cmd := vm.listProfile.SetItem(idx, targetProfile)
        return vm,cmd
    case tea.KeyMsg:
        if vm.listProfile.FilterState() == list.Filtering {
            //do nothing
            break
        }
        switch msg.String(){
        case "esc":
            //esc will work for filter
             if vm.listProfile.FilterState() == list.Filtering {
                 break
             }
             if vm.listProfile.FilterState() == list.FilterApplied {
                 vm.listProfile.ResetFilter()
                 return vm,nil
             }
            vm.quit = true
            core.Wipe(vm.key)
            return vm,nil
        case "c":
            //TODO copy the password of this profile to clipboard and clear clipboard after 10 sec
        case "v":
            selectedProfile := vm.listProfile.SelectedItem()
            if selectedProfile == nil {return vm,nil}
            selectedProfileName := string(selectedProfile.(profileName).name)
            return vm, func() tea.Msg {

                return ViewPassMsg {profileName : selectedProfileName}
            }
        case "ctrl+d":
             //TODO delete the focused profile's data
         
        }
    }
    var cmd tea.Cmd
    //case blinking
    //case list updating [pass msg to the list]
    vm.listProfile, cmd = vm.listProfile.Update(msg)
    return vm,cmd
}
func (vm VaultModel) View() string {
    
    title := titleStyle.Render(title)

    var s string
    h := helpStyle.Render("[To Search: / |Copy: c | View: v | Delete: ctrl+d | Quit: esc ]")
    s = lipgloss.JoinVertical(lipgloss.Left,title,vm.listProfile.View(),h)


    return lipgloss.Place(10,10,lipgloss.Center,lipgloss.Center,s)
}
