package tui

import (
	"fmt"
	// "time"

	"github.com/Kasbe14/PassMan/core"
	"github.com/Kasbe14/PassMan/database"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	// "strings"
	"os"
)

//styles
const (
    title = "PassMan"
    width = 60
    charLim = 32
   )
var (
    
    titleStyle_              = lipgloss.NewStyle().Foreground(lipgloss.Color("#9FA1FF")).Align(lipgloss.Center).Bold(true).
                               MarginBottom(2).Width(width)
    focusedStyle             = lipgloss.NewStyle().Foreground(lipgloss.Color("#AEE2FF"))
    blurredStyle            = lipgloss.NewStyle().Foreground(lipgloss.Color("242"))
    // cursorStyle             = focusedStyle
    // focusedButton           = focusedStyle.Render("[ Submit  ]")
    helpStyle               = lipgloss.NewStyle().Foreground(lipgloss.Color("242")).Align(lipgloss.Left)
    inputStyle              = lipgloss.NewStyle().Align(lipgloss.Left).Width(width)
                              // Border(lipgloss.NormalBorder(),false,false,true,false).
                              // BorderBottom(true).BorderForeground(lipgloss.Color("#AEE2FF")).Width(width)
    textStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#D9F9DF"))
    errStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("9")).Align(lipgloss.Left)
    recStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#9FA1FF")).Align(lipgloss.Center).Bold(true)
    help = (helpStyle.Render("[ Next: tab/up/down | Submit: enter | Back: esc | Quit: ctrl+c ]"))
)

type sessionState int

const (
    stateStart sessionState = iota + 1
    stateSignUp
    stateLogin
    stateVault
    stateDashBoard //login dashboard
    stateAddProfile //adding profile
)

//custom msg from the core
type CoreRegisterUserMsg struct {
    recKey string
    err    error
}
type CoreLoginUserMsg struct {
    key []byte
    userId int64
    err   error
}

 

type MainModel struct {
    vault *core.VaultService
    state sessionState
    sm     StartUpModel
    su     SignUpModel
    lm     LoginModel
    vm     VaultModel
    dbm    DashBoardModel
    apm    AddProfileModel
    text string
}
type MainModelMsg struct {
    Err string
}

func (m MainModel) Init() tea.Cmd {

    return tea.Batch(
        m.sm.Init(),
        m.su.Init(),
        m.lm.Init(),
        m.vm.Init(),
        m.dbm.Init(),
    )
}

func (m MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    //Global msg for state change
    switch msg:= msg.(type) {
    case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
            core.Wipe(m.vm.key)
			return m, tea.Quit
		}

    case SignUpMsg://after signup gets valid input returns msg to call backend
        cmd := func () tea.Msg {
            generatedString, err := m.vault.RegisterUser(msg.User,msg.Pass)
            if err != nil {
                return CoreRegisterUserMsg{recKey: "", err: err}
            }
            return CoreRegisterUserMsg{recKey: generatedString, err: nil}
        }
        return m,cmd
    case CoreRegisterUserMsg://msg gets the result form backend and updates signup models acrodingly
        if msg.err != nil {
            m.su.suMsg.Err = msg.err.Error()
            m.su.success = false
        }else {
            m.su.userRecString = msg.recKey
            m.su.success = true
        }
        //clearing the msgreckey
        msg.recKey = ""
        return m,nil
        
    case LoginMsg:
        cmd := func () tea.Msg {
            key, userId, err := m.vault.LoginUser(msg.User,msg.Pass) 
            if err != nil {
                return CoreLoginUserMsg{
                    key: nil,
                    userId: 0,
                    err: err,
                }
            }
            return CoreLoginUserMsg{
                key: key,
                userId: userId,
                err: nil,
            }
        }
        return m,cmd
    case CoreLoginUserMsg:
        if msg.err != nil {
            m.lm.lmMsg.Err = msg.err.Error()
            m.lm.submit = false
        }else {
            m.lm.submit = true
            //goto dashboardstate
            m.state = /*stateVault*/ stateDashBoard
            //initialize the dashboard state
            m.dbm = initialDModel()
            //passing the keys and userId to vaultstate
            m.vm = initialVModel(msg.key,msg.userId)
            // m.vm.key = msg.key
            // m.vm.userId = msg.userId
            //overwritng the login msg variables[!only local copies of pass and user is cleaned]
            core.Wipe(msg.key) 
            msg.userId = 0
            m.lm.lmMsg.Pass = ""
            m.lm.lmMsg.User = ""
        }
    case AddProMsg:
        
    case ViewPassMsg:
        if msg.profileName == "" {
            return m,nil
        }
        //uncomment for database/connectivity
        //else {
        //    return m, func() tea.Msg {
        //        pro ,err := m.vault.GetProfileByName(msg.profile,m.vm.key)
        //        if err != nil {
        //            return RevealPassMsg{err: err, pass:nil,}
        //        }
        //        if pro.Locked {
        //            return RevealPassMsg{err:err,pass:nil,lockMsg:pro.LockedMsg}
        //        }
        //        //calling backned and just in time decrypt
        //        plaintext, err := m.vault.Decrypt(pro.Password,m.vm.key)
        //        if err != nil {
        //           return RevealPassMsg{err: err, pass:nil, lockMsg:pro.LockedMsg}
        //        }
        //        return RevealPassMsg{err:nil,pass:plaintext,lockMsg:pro.LockedMsg,listIndex: msg.idx}
        //    }
        //}
        // item := m.vm.listProfile.Items()[msg.idx].(profileName)
        var targetItem profileName
    for _, item := range m.vm.listProfile.Items() {
        p := item.(profileName)
        if p.name == msg.profileName {
            targetItem = p
            break
        }
    }
    
         return m, func() tea.Msg {
             return RevealPassMsg{
                err:       nil,
                pass:      targetItem.decryptedPass, // Sending the dummy bytes directly
               // listIndex: msg.idx,
               proName: msg.profileName,
        }}
    case RevealPassMsg:
        if msg.err != nil {
            //TODO render the error on the screen if any error form backend
            return m, nil
        }
        var cmd tea.Cmd
        updatedModel ,subCmd := m.vm.Update(msg)
        m.vm = updatedModel.(VaultModel)
        cmd = subCmd
        return m,cmd


    }

	var cmd tea.Cmd

	switch m.state {
	case stateStart:
		// Intercept the enter key BEFORE handing it down
		if keyMsg, ok := msg.(tea.KeyMsg); ok && keyMsg.String() == "enter" {
            
			if m.sm.activeTab == 0 {
				m.state = stateSignUp // Tab 0 is SignUp
                m.su = initialSModel()
                cmd = m.su.Init()
                return m,cmd
			} else if m.sm.activeTab == 1 {
				m.state = stateLogin  // Tab 1 is LogIn
                m.lm = initialLModel()
                cmd = m.lm.Init()
                return m,cmd
			}
			return m, cmd 
        }

        //continue nav on startpage
		updatedModel, subCmd := m.sm.Update(msg)
		m.sm = updatedModel.(StartUpModel)
		cmd = subCmd

	  case stateSignUp:
          if keyMsg, ok := msg.(tea.KeyMsg); ok && (keyMsg.String() == "esc" && m.su.quitSignUp == true) {
              m.su.userRecString = ""
              m.state = stateStart
              return m,nil
          }
        // m.su = initialSModel()
        // m.su.Init()
          updatedModel, subCmd := m.su.Update(msg)
           m.su = updatedModel.(SignUpModel)
         cmd = subCmd
	  case stateLogin:
          if keyMsg,ok := msg.(tea.KeyMsg);ok && (keyMsg.String() == "esc" && m.lm.quitLogin == true){
              m.state = stateStart
          }
          updatedModel, subCmd := m.lm.Update(msg)
          m.lm = updatedModel.(LoginModel)
          cmd = subCmd
      case stateVault :
          //TODO vault state
          if keyMsg, ok := msg.(tea.KeyMsg); ok && (keyMsg.String() == "esc" && m.vm.quit == true) {
              m.state = stateDashBoard
              // core.Wipe(m.vm.key)
              return m,nil
          }
          updatedModel, subCmd := m.vm.Update(msg)
          m.vm = updatedModel.(VaultModel)
          cmd = subCmd
      case stateDashBoard :
          if keyMsg, ok := msg.(tea.KeyMsg); ok && (keyMsg.String() == "esc" && m.dbm.quit == true) {
              m.state = stateStart
              core.Wipe(m.vm.key)
              return m,nil
          }
          if keyMsg, ok := msg.(tea.KeyMsg); ok && (keyMsg.String() == "enter") {
              if m.dbm.activeTab == 0 {
                  //TODO go to add profile model
                  m.state = stateAddProfile
                  return m,nil
              }else if m.dbm.activeTab == 1 {
                  m.state = stateVault
                  return m,nil
              }
          }
          updatedModel, subCmd := m.dbm.Update(msg)
          m.dbm = updatedModel.(DashBoardModel)
          cmd = subCmd
      case stateAddProfile:
          if keyMsg, ok := msg.(tea.KeyMsg); ok && (keyMsg.String() == "esc" && m.apm.quit == true) {
              m.state = stateDashBoard
              return m,nil
          }
          updatedModel, subCmd := m.apm.Update(msg)
          m.apm = updatedModel.(AddProfileModel)
          cmd = subCmd
	  }

	return m, cmd
}

func InitialMainModel(sm StartUpModel,vs *core.VaultService) MainModel {
    return MainModel{
        vault: vs,
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
    case stateLogin:
        return m.lm.View()
    case stateVault:
        return m.vm.View()
    case stateDashBoard:
        return m.dbm.View()
    case stateAddProfile:
        return m.apm.View()
    }
    
    return "Loading.."
}

//tui entry point
func Start()  {
    f, err := tea.LogToFile("debug.log", "debug")
    if err != nil {
        fmt.Println("fatal:", err)
        os.Exit(1)
    }
    defer f.Close()
    db,err := database.NewDatabase()
    if err != nil {
        fmt.Println("fatal:",err)
        os.Exit(1)
    }
    err = database.InitializeSchema(db)
    if err != nil {
        fmt.Println("fatal:",err)
        os.Exit(1)
    }

    vs := core.NewVaultService(db)
    p := tea.NewProgram(

        InitialMainModel(StartUpModel{Tabs: []string{"SignUp","LogIn"},activeTab:0},vs),
        tea.WithAltScreen(),//renders new terminal
    )
    if _,err := p.Run(); err != nil {
        fmt.Printf("Alas, there's been an error: %v", err)
        os.Exit(1)
    }
}
