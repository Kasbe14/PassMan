package tui

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
) 

var (

)

type LoginModel struct {
    // title string
    inputs  []textinput.Model
    lmMsg    LoginMsg //mainmodle updates the loginmodel
    activeInput int
    submit bool   
    quitLogin bool
}
type LoginMsg struct {
    Err string
    User string
    Pass string
}


func (lm LoginModel) Init() tea.Cmd {
    return textinput.Blink
}
func initialLModel() LoginModel {
    lm := LoginModel{
        inputs: make([]textinput.Model,2),
    }
    lm.submit = false
    lm.activeInput = 0
    lm.quitLogin = false
    // lm.title = "PassMan"
    for i := range lm.inputs {
        t := textinput.New()
        t.CharLimit = charLim
        switch i {
        case 0:
            t.Placeholder = "Enter Username"
            t.TextStyle = textStyle
            t.Focus()
            t.Cursor.BlinkCmd()
            t.PromptStyle = focusedStyle
            lm.inputs[i] = t
        case 1:
            t.Placeholder = "Enter User Password"
            t.EchoMode = textinput.EchoPassword
            t.Blur()
            t.PromptStyle = blurredStyle
            t.EchoCharacter = '•'
            lm.inputs[i] = t
        }

    }
    return lm
}
func (lm LoginModel) Update(msg tea.Msg) (tea.Model,tea.Cmd) {
    switch msg := msg.(type){
    case tea.KeyMsg:
        inputLen := len(lm.inputs)
        switch msg.String() {
        case "esc" :
            lm.quitLogin = true
            return lm,nil
        case "tab","down","up":
            if msg.String()=="tab" || msg.String() == "down"{
                lm.activeInput++
            }else {
                lm.activeInput--
            }
            //toggle same
            if lm.activeInput > inputLen -1 {
                lm.activeInput = 0
            }
            if lm.activeInput < 0 {
                lm.activeInput++
            }
            cmds := toggleActiveInputs(lm.inputs,lm.activeInput)
            return lm, tea.Batch(cmds...)
        case "enter":
            if lm.activeInput == inputLen-1{
                 username := lm.inputs[0].Value()
                 password := lm.inputs[1].Value()
                 cmd := func () tea.Msg {
                     return LoginMsg{
                         Err: "",
                         User: username,
                         Pass: password,
                     }
                 }
                 lm.submit = true
                 return lm, cmd
            }
            lm.activeInput++
            cmds := toggleActiveInputs(lm.inputs,lm.activeInput)
            return lm,tea.Batch(cmds...)
        }

        //case keyblinking msg
        cmd := updateInputs(msg,lm.inputs)
        return lm,cmd
    }
     var cmds []tea.Cmd
     for i := range lm.inputs {
         var cmd tea.Cmd
         //updating the textinput model
         lm.inputs[i],cmd = lm.inputs[i].Update(msg)
         cmds = append(cmds,cmd)
     }
    return lm,tea.Batch(cmds...)
}
func (lm LoginModel) View() string {
    title := titleStyle_.Render(title)
    u := lm.inputs[0].View()
    p := lm.inputs[1].View()
    styledU := inputStyle.Render(u)
    styledP := inputStyle.Render(p)

    
    var s string
    e := errStyle.Render(lm.lmMsg.Err)
    if lm.lmMsg.Err != "" {
        s = lipgloss.JoinVertical(lipgloss.Center,title,styledU,styledP,e,"\n",help)
        
    }else{
        s = lipgloss.JoinVertical(lipgloss.Center,title,styledU,styledP,"\n",help)
    }
  
    return lipgloss.Place(10,10,lipgloss.Center,lipgloss.Center,s)

}
