package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// --- Design System ---

var (
	docStyle = lipgloss.NewStyle().Margin(1, 2)
	
	titleStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("62")).
			Foreground(lipgloss.Color("230")).
			Padding(0, 1)

	statusStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241")).
			Italic(true)

	selectedStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder(), false, false, false, true).
			BorderForeground(lipgloss.Color("170")).
			PaddingLeft(1).
			Foreground(lipgloss.Color("170"))

	helpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
)

// --- Data Model ---

type item struct {
	title, desc string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

type Note struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

// --- App State ---

type state int

const (
	viewList state = iota
	viewNote
	addNote
)

type model struct {
	list     list.Model
	notes    []Note
	state    state
	titleIn  textinput.Model
	contentIn textarea.Model
	selected int
}

func initialModel() model {
	// Initialize inputs
	ti := textinput.New()
	ti.Placeholder = "Note Title..."
	ti.Focus()

	ta := textarea.New()
	ta.Placeholder = "Note Content..."

	// Load notes
	notes, _ := loadNotes()
	items := make([]list.Item, len(notes))
	for i, n := range notes {
		items[i] = item{title: n.Title, desc: n.Content}
	}

	l := list.New(items, list.NewDefaultDelegate(), 0, 0)
	l.Title = "My Notes"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(true)
	l.Styles.Title = titleStyle

	return model{
		list:      l,
		notes:     notes,
		state:     viewList,
		titleIn:   ti,
		contentIn: ta,
	}
}

// --- Persistence ---

func loadNotes() ([]Note, error) {
	data, err := os.ReadFile("notes.json")
	if err != nil {
		return []Note{}, nil
	}
	var notes []Note
	err = json.Unmarshal(data, &notes)
	return notes, err
}

func (m *model) saveNotes() {
	data, _ := json.Marshal(m.notes)
	_ = os.WriteFile("notes.json", data, 0644)
}

// --- Lifecycle ---

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
		m.contentIn.SetWidth(msg.Width - h)
		m.contentIn.SetHeight(msg.Height - v - 10)

	case tea.KeyMsg:
		if m.state == viewList {
			switch msg.String() {
			case "q", "ctrl+c":
				return m, tea.Quit
			case "n":
				m.state = addNote
				m.titleIn.Focus()
				return m, nil
			case "enter":
				if len(m.notes) > 0 {
					m.selected = m.list.Index()
					m.state = viewNote
				}
				return m, nil
			case "d":
				if len(m.notes) > 0 {
					idx := m.list.Index()
					m.notes = append(m.notes[:idx], m.notes[idx+1:]...)
					m.saveNotes()
					m.list.RemoveItem(idx)
				}
				return m, nil
			}
		} else if m.state == addNote {
			switch msg.String() {
			case "esc":
				m.state = viewList
				return m, nil
			case "tab":
				if m.titleIn.Focused() {
					m.titleIn.Blur()
					m.contentIn.Focus()
				} else {
					m.contentIn.Blur()
					m.titleIn.Focus()
				}
				return m, nil
			case "enter":
				if m.contentIn.Focused() {
					newNote := Note{Title: m.titleIn.Value(), Content: m.contentIn.Value()}
					m.notes = append(m.notes, newNote)
					m.saveNotes()
					m.list.InsertItem(len(m.notes)-1, item{title: newNote.Title, desc: newNote.Content})
					m.titleIn.Reset()
					m.contentIn.Reset()
					m.state = viewList
				} else {
					m.titleIn.Blur()
					m.contentIn.Focus()
				}
				return m, nil
			}
		} else if m.state == viewNote {
			switch msg.String() {
			case "esc", "q":
				m.state = viewList
				return m, nil
			}
		}
	}

	// Update components based on state
	if m.state == viewList {
		m.list, cmd = m.list.Update(msg)
	} else if m.state == addNote {
		if m.titleIn.Focused() {
			m.titleIn, cmd = m.titleIn.Update(msg)
		} else {
			m.contentIn, cmd = m.contentIn.Update(msg)
		}
	}

	return m, cmd
}

func (m model) View() string {
	var s string

	switch m.state {
	case viewList:
		s = docStyle.Render(m.list.View())
	case addNote:
		s = docStyle.Render(
			fmt.Sprintf(
				"Add New Note\n\n%s\n\n%s\n\n%s",
				m.titleIn.View(),
				m.contentIn.View(),
				helpStyle.Render("(tab to switch, enter to save, esc to cancel)"),
			),
		)
	case viewNote:
		note := m.notes[m.selected]
		s = docStyle.Render(
			fmt.Sprintf(
				"%s\n\n%s\n\n%s",
				titleStyle.Render(note.Title),
				note.Content,
				helpStyle.Render("(esc/q to go back)"),
			),
		)
	}

	return s
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}
