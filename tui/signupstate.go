package tui

import (
	// "fmt"
	// "strings"
	// "github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

//styles
const (
    width = 60
    charLim = 32
   )
var (
    
    titleStyle_              = lipgloss.NewStyle().Foreground(lipgloss.Color("#9FA1FF")).Align(lipgloss.Center).Bold(true).
                               MarginBottom(2).Width(width)
    focusedStyle             = lipgloss.NewStyle().Foreground(lipgloss.Color("#AEE2FF"))
    blurredStyle            = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
    cursorStyle             = focusedStyle
    focusedButton           = focusedStyle.Render("[ Submit  ]")
    helpStyle               = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
    inputStyle              = lipgloss.NewStyle().Align(lipgloss.Left).
                              Border(lipgloss.NormalBorder(),false,false,true,false).
                              BorderBottom(true).BorderForeground(lipgloss.Color("#AEE2FF")).Width(width)
    textStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#D9F9DF"))
)


type SignUpModel struct {
    title      string
    inputs    []textinput.Model //handles username and password
    activeInput int
    quitSignUp bool
    
}
//message for error and to go back to home
type BackMsg struct {
    backMsg   string
    quit bool
}

func (su SignUpModel)  Init() tea.Cmd {
    // return su.inputs[0].Cursor.BlinkCmd()
    
    return  textinput.Blink
}
func initialSModel() SignUpModel {
    su:= SignUpModel{
        inputs: make([]textinput.Model,3),
    }
    su.title = "PassMan"
    su.quitSignUp = false
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
func (su SignUpModel) updateInputs(msg tea.Msg) tea.Cmd{
    cmds := make([]tea.Cmd,len(su.inputs))

    for i := range su.inputs {
        //calling textinput models to update
        su.inputs[i], cmds[i] = su.inputs[i].Update(msg)
    }
    return tea.Batch(cmds...)
}
func (su SignUpModel)  Update(msg tea.Msg) (tea.Model,tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        inputLen := len(su.inputs)
        switch msg.String(){
        case "q":
            su.quitSignUp = true
            return su,nil
        case "enter":
            //TODO handling inputs
            if su.activeInput == inputLen-1{
                //TODO get the inputs values using Value method and input validation then passing to backend
                su.quitSignUp = true
                return su,nil
            }
            //go to next input if not on last input
            su.activeInput++
            cmds := toggleActiveInputs(su)
            return su, tea.Batch(cmds...)
            // fallthrough
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
            cmds := toggleActiveInputs(su)
            return su,tea.Batch(cmds...)
    }
    // case cursor.BlinkMsg:
        cmd := su.updateInputs(msg)
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
    title := titleStyle_.Render(su.title)
    u := su.inputs[0].View()
    p := su.inputs[1].View()
    cp := su.inputs[2].View()
    
    styledU := inputStyle.Render(u)
    styledP := inputStyle.Render(p)
    styledC :=inputStyle.Render(cp)


    help := (helpStyle.Render("[ Next: tab/up/down | Submit: enter | Quit: q ]"))
    
    s := lipgloss.JoinVertical(lipgloss.Center,title,styledU,styledP,styledC,"\n",help)
  
    return lipgloss.Place(10,10,lipgloss.Center,lipgloss.Center,s)
    // return su.title
}

func toggleActiveInputs(su SignUpModel) []tea.Cmd {
            cmds := make([]tea.Cmd, len(su.inputs))
            for i := 0; i<len(su.inputs); i++ {
                if i == su.activeInput {
                    cmds[i] = su.inputs[i].Focus()
                    su.inputs[i].PromptStyle = focusedStyle
                    su.inputs[i].TextStyle = textStyle
                    
                    cmds = append(cmds,su.inputs[i].Cursor.BlinkCmd())
                    continue
                }
                //if not active blur it
               // if i!= su.activeInput{
                   su.inputs[i].Blur()
                   su.inputs[i].PromptStyle = blurredStyle
                   su.inputs[i].TextStyle = blurredStyle
               // } 

            }
            return cmds
}
