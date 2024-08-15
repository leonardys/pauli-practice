package internal

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type Session struct {
	StartTime, Endtime time.Time
	Logs               []Question
}

type Question struct {
	Q1, Q2, Answer  int
	Shown, Answered time.Time
}

type Model struct {
	CurrQuestion               Question
	ThisSession                Session
	CorrectAnswer, TotalAnswer int
	Keymap                     Keymap
	Help                       help.Model
}

type Keymap struct {
	Quit key.Binding
}

func New() Model {
	m := Model{
		CurrQuestion: Question{Q1: rand.Intn(10), Q2: rand.Intn(10), Shown: time.Now()},
		ThisSession:  Session{StartTime: time.Now(), Logs: make([]Question, 0)},
		Keymap: Keymap{
			Quit: key.NewBinding(
				key.WithKeys("q"),
				key.WithHelp("q", "quit"),
			),
		},
		Help: help.New(),
	}
	return m
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) helpView() string {
	return "\n" + m.Help.ShortHelpView([]key.Binding{
		m.Keymap.Quit,
	})
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.Keymap.Quit):
			m.ThisSession.Endtime = time.Now()
			return m, saveLogToJSON(m)
		default:
			if v, err := strconv.Atoi(msg.String()); err == nil {
				if v == (m.CurrQuestion.Q1+m.CurrQuestion.Q2)%10 {
					m.CorrectAnswer += 1
				}
				m.CurrQuestion.Answer = v
				m.CurrQuestion.Answered = time.Now()

				m.TotalAnswer += 1

				m.ThisSession.Logs = append(m.ThisSession.Logs, m.CurrQuestion)
				m.CurrQuestion = Question{Q1: m.CurrQuestion.Q2, Q2: rand.Intn(10), Shown: time.Now()}
			}
		}
	case saveOkMsg:
		return m, tea.Quit
	}

	return m, nil
}

func (m Model) View() string {
	s := "Pauli Test Practice\n\n"
	s += fmt.Sprintf("Correct Answers: %d / %d\n\n", m.CorrectAnswer, m.TotalAnswer)
	s += fmt.Sprintf("%d\n%d\n", m.CurrQuestion.Q1, m.CurrQuestion.Q2)
	s += m.helpView()

	return s
}

type saveErrMsg struct{ err error }
type saveOkMsg int

func saveLogToJSON(m Model) tea.Cmd {
	return func() tea.Msg {
		err := saveData("data.json", m.ThisSession)

		if err != nil {
			return saveErrMsg{err}
		}

		return saveOkMsg(1)
	}
}
