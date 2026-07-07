package tui

import (
	// "fmt"
	// "strings"
	// "github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)



type SignUpModel struct {
    // title      string
    inputs    []textinput.Model //handles username and password
    activeInput int
    success   bool
    quitSignUp bool
    suMsg       SignUpMsg//passing and changin state of the SignUpModel
    userRecString string
    
}
//message for error and to go back to home
type SignUpMsg struct {
    Err  string
    User string
    Pass string
}

func (su SignUpModel)  Init() tea.Cmd {
    // return su.inputs[0].Cursor.BlinkCmd()
    
    return  textinput.Blink
}
func initialSModel() SignUpModel {
    su:= SignUpModel{
        inputs: make([]textinput.Model,3),
    }
    // su.title = "PassMan"
    su.quitSignUp = false
    su.success = false
    for i := range su.inputs {
        t := textinput.New()
        t.CharLimit = charLim
        switch i{
        case 0:
            t.Placeholder = "Enter Username"
            t.Focus()
            t.TextStyle = textStyle
            t.PromptStyle= focusedStyle
            t.Cursor.BlinkCmd()
            su.inputs[i] = t
        case 1:
            t.Placeholder = "Enter Password"
            t.Blur()
            t.PromptStyle = blurredStyle
            t.EchoMode = textinput.EchoPassword
            t.EchoCharacter = '•'
            su.inputs[i] = t
        case 2:
            t.Placeholder = "Confirm Password"
            t.Blur()
            t.PromptStyle = blurredStyle
            t.EchoMode = textinput.EchoPassword
            t.EchoCharacter = '•'
            su.inputs[i] = t
        }
    }
    
    return su
    
}
//takes the cursorBlinkmsg and handls for each focused input
func updateInputs(msg tea.Msg,inputs []textinput.Model) tea.Cmd{
    cmds := make([]tea.Cmd,len(inputs))

    for i := range inputs {
        //calling textinput models to update
        inputs[i], cmds[i] = inputs[i].Update(msg)
    }
    return tea.Batch(cmds...)
}
func (su SignUpModel)  Update(msg tea.Msg) (tea.Model,tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        inputLen := len(su.inputs)
        switch msg.String(){
        case "esc":
            su.quitSignUp = true
            return su,nil
        case "enter":
            //TODO handling inputs
            if su.activeInput == inputLen-1{
                username := su.inputs[0].Value()
                pass    := su.inputs[1].Value()
                passC := su.inputs[2].Value()
                ulen := len(username)
                if pass != passC && ulen > 6 {
                    su.success = false
                    su.suMsg.Err = "password doesn't match"
                    return su,nil
                    
                }else if ulen > 6 && len(pass) < 6 {
                    su.success = false
                    su.suMsg.Err = "password must be atleast 6 character"
                }
                if len(username) < 6 {
                    su.success = false
                    su.suMsg.Err = "username must be atleast 6 character"
                    return su,nil
                }
                cmd := func () tea.Msg {return SignUpMsg{
                    Err: "",
                    User: username,
                    Pass: pass,
                }}
                su.success = true
                su.suMsg.Err = ""
                return su,cmd
            }
        
            //go to next input if not on last input
            su.activeInput++
            cmds := toggleActiveInputs(su.inputs,su.activeInput)
            return su, tea.Batch(cmds...)
            // fallthroug
            //naivgating inputs 
          case "tab","up","down":
            if msg.String() == "down" || msg.String() == "tab" {
                su.activeInput++
            }else {
                su.activeInput--
            }
            //toggle 
            if su.activeInput > inputLen-1 {
                su.activeInput = 0
            }
            if su.activeInput < 0 {
                // su.activeInput = len(su.inputs)-1
                su.activeInput++ 
            }
            
            //updating the cmds to focus and blink the cursor for the active input
            cmds := toggleActiveInputs(su.inputs,su.activeInput)
            return su,tea.Batch(cmds...)
    }
    // case cursor.BlinkMsg:
        cmd := updateInputs(msg,su.inputs)
        return su,cmd
    
}
//individual inputs cmds
     var cmds []tea.Cmd
     for i := range su.inputs {
         var cmd tea.Cmd
         //updating the textinput model
         su.inputs[i],cmd = su.inputs[i].Update(msg)
         cmds = append(cmds,cmd)
     }
    return su,tea.Batch(cmds...)
}

func (su SignUpModel)  View() string {
    title := titleStyle_.Render(title)
    u := su.inputs[0].View()
    p := su.inputs[1].View()
    cp := su.inputs[2].View()
    
    styledU := inputStyle.Render(u)
    styledP := inputStyle.Render(p)
    styledC :=inputStyle.Render(cp)
    //if singup is submitted and sucess view signup sucess message with no inputs
    
    var s string
    e := errStyle.Render(su.suMsg.Err)
    //view err
    if(!su.success && su.suMsg.Err != "") {
        s = lipgloss.JoinVertical(lipgloss.Center,title,styledU,styledP,styledC,e,"\n",help)
    }
    suc := focusedStyle.Render("SignUp Successful !")
    warning := errStyle.Render("Important: Save this recovery key. Else forget your data if you forget your password")
    recString := recStyle.Render(su.userRecString)
    if (su.success && su.suMsg.Err == "") {
        s = lipgloss.JoinVertical(lipgloss.Center,title,suc,"\n",warning,recString,"\n",help)
    }
    if (!su.success && su.suMsg.Err == "") {
        s = lipgloss.JoinVertical(lipgloss.Center,title,styledU,styledP,styledC,"\n",help)
    }
  
    return lipgloss.Place(10,10,lipgloss.Center,lipgloss.Center,s)
    // return su.title
}

func toggleActiveInputs(inputs []textinput.Model, activeInput int) ([]tea.Cmd ){
            cmds := make([]tea.Cmd, len(inputs))
            for i := 0; i<len(inputs); i++ {
                if i == activeInput {
                    cmds[i] = inputs[i].Focus()
                    inputs[i].PromptStyle = focusedStyle
                    inputs[i].TextStyle = textStyle
                    
                    cmds = append(cmds,inputs[i].Cursor.BlinkCmd())
                    continue
                }
                //if not active blur it
               // if i!= su.activeInput{
                   inputs[i].Blur()
                   inputs[i].PromptStyle = blurredStyle
                   inputs[i].TextStyle = blurredStyle
               // } 

            }
            return cmds
}
