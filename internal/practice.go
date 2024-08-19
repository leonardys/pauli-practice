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

type practiceModel struct {
	currQuestion               Question
	ThisSession                Session
	correctAnswer, totalAnswer int
	keymap                     keymap
	help                       help.Model
}

type keymap struct {
	quit key.Binding
}

func NewPracticeModel() practiceModel {
	m := practiceModel{
		currQuestion: Question{Q1: rand.Intn(10), Q2: rand.Intn(10), Shown: time.Now()},
		ThisSession:  Session{StartTime: time.Now(), Logs: make([]Question, 0)},
		keymap: keymap{
			quit: key.NewBinding(
				key.WithKeys("q"),
				key.WithHelp("q", "quit"),
			),
		},
		help: help.New(),
	}
	return m
}

func (m practiceModel) Init() tea.Cmd {
	return nil
}

func (m practiceModel) helpView() string {
	return "\n" + m.help.ShortHelpView([]key.Binding{
		m.keymap.quit,
	})
}

func (m practiceModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keymap.quit):
			m.ThisSession.Endtime = time.Now()
			return m, saveLogToJSON(m)
		default:
			if v, err := strconv.Atoi(msg.String()); err == nil {
				if v == (m.currQuestion.Q1+m.currQuestion.Q2)%10 {
					m.correctAnswer += 1
				}
				m.currQuestion.Answer = v
				m.currQuestion.Answered = time.Now()

				m.totalAnswer += 1

				m.ThisSession.Logs = append(m.ThisSession.Logs, m.currQuestion)
				m.currQuestion = Question{Q1: m.currQuestion.Q2, Q2: rand.Intn(10), Shown: time.Now()}
			}
		}
	case saveOkMsg:
		return m, tea.Quit
	}

	return m, nil
}

func (m practiceModel) View() string {
	s := "Pauli Test Practice\n\n"
	s += fmt.Sprintf("Correct Answers: %d / %d\n\n", m.correctAnswer, m.totalAnswer)
	s += fmt.Sprintf("%d\n%d\n", m.currQuestion.Q1, m.currQuestion.Q2)
	s += m.helpView()

	return s
}

type saveErrMsg struct{ err error }
type saveOkMsg int

func saveLogToJSON(m practiceModel) tea.Cmd {
	return func() tea.Msg {
		if len(m.ThisSession.Logs) > 0 {
			err := saveData("data.json", m.ThisSession)

			if err != nil {
				return saveErrMsg{err}
			}
		}

		return saveOkMsg(1)
	}
}
