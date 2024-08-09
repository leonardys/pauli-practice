package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type session struct {
	StartTime, Endtime time.Time
	Logs               []question
}

type question struct {
	Q1, Q2, Answer  int
	Shown, Answered time.Time
}

type model struct {
	currQuestion               question
	thisSession                session
	correctAnswer, totalAnswer int
	keymap                     keymap
	help                       help.Model
}

type keymap struct {
	quit key.Binding
}

func main() {
	m := model{
		currQuestion: question{Q1: rand.Intn(10), Q2: rand.Intn(10), Shown: time.Now()},
		thisSession:  session{StartTime: time.Now(), Logs: make([]question, 0)},
		keymap: keymap{
			quit: key.NewBinding(
				key.WithKeys("q"),
				key.WithHelp("q", "quit"),
			),
		},
		help: help.New(),
	}

	p := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) helpView() string {
	return "\n" + m.help.ShortHelpView([]key.Binding{
		m.keymap.quit,
	})
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keymap.quit):
			m.thisSession.Endtime = time.Now()
			return m, saveLogToJSON(m)
		default:
			if v, err := strconv.Atoi(msg.String()); err == nil {
				if v == (m.currQuestion.Q1+m.currQuestion.Q2)%10 {
					m.correctAnswer += 1
				}
				m.currQuestion.Answer = v
				m.currQuestion.Answered = time.Now()

				m.totalAnswer += 1

				m.thisSession.Logs = append(m.thisSession.Logs, m.currQuestion)
				m.currQuestion = question{Q1: m.currQuestion.Q2, Q2: rand.Intn(10), Shown: time.Now()}
			}
		}
	case saveOkMsg:
		return m, tea.Quit
	}

	return m, nil
}

func (m model) View() string {
	s := "Pauli Test Trainer\n\n"
	s += fmt.Sprintf("Correct Answers: %d / %d\n\n", m.correctAnswer, m.totalAnswer)
	s += fmt.Sprintf("%d\n%d\n", m.currQuestion.Q1, m.currQuestion.Q2)
	s += m.helpView()

	return s
}

type saveErrMsg struct{ err error }
type saveOkMsg int

func saveLogToJSON(m model) tea.Cmd {
	return func() tea.Msg {
		var sessions []session

		data, err := os.ReadFile("data.json")
		if err != nil {
			sessions = make([]session, 0)

		} else {
			json.Unmarshal(data, &sessions)
		}

		sessions = append(sessions, m.thisSession)

		json, _ := json.Marshal(sessions)
		err = os.WriteFile("data.json", json, 0644)
		if err != nil {
			return saveErrMsg{err}
		}

		return saveOkMsg(1)
	}
}
