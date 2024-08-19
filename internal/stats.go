package internal

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type sessionItem struct {
	session Session
}

type statsModel struct {
	sessionList list.Model
}

func (i sessionItem) Title() string       { return i.session.StartTime.String() }
func (i sessionItem) Description() string { return i.session.Endtime.Sub(i.session.StartTime).String() }
func (i sessionItem) FilterValue() string { return i.session.StartTime.String() }

func NewStatsModel() statsModel {
	sessions := loadData("data.json")

	items := make([]list.Item, len(sessions))

	for i, s := range sessions {
		items[i] = sessionItem{session: s}
	}

	m := statsModel{
		sessionList: list.New(items, list.NewDefaultDelegate(), 0, 0),
	}
	m.sessionList.Title = "Previous Practice Sessions"
	m.sessionList.SetFilteringEnabled(false)

	return m
}

func (m statsModel) Init() tea.Cmd {
	return nil
}

func (m statsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.sessionList.SetSize(msg.Width, msg.Height)
	}
	var cmd tea.Cmd
	m.sessionList, cmd = m.sessionList.Update(msg)
	return m, cmd
}

func (m statsModel) View() string {
	return m.sessionList.View()
}
