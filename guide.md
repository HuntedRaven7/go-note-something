# Building a Terminal Note-Taking App in Go

This guide will walk you through creating a professional-grade terminal application using Go. We'll use the **Charm Bracelet** ecosystem (Bubble Tea, Lip Gloss) to create a beautiful, interactive interface.

## Prerequisites
- Go installed on your system.
- Basic knowledge of Go syntax.

---

## Step 1: Project Initialization

First, we need to set up our Go module and install the necessary libraries.

```bash
mkdir go-note-something
cd go-note-something
go mod init go-note-something
go get github.com/charmbracelet/bubbletea github.com/charmbracelet/bubbles github.com/charmbracelet/lipgloss
```

- **Bubble Tea**: The TUI (Terminal User Interface) framework based on the Elm Architecture.
- **Lip Gloss**: Used for styling the terminal output (colors, borders, etc.).
- **Bubbles**: A library of common TUI components (text inputs, lists).

---

## Step 2: Defining the Data Structure

We need a way to represent a note. In `main.go`, we'll start by defining a `Note` struct and a simple way to store them.

```go
type Note struct {
    Title   string `json:"title"`
    Content string `json:"content"`
}
```

---

## Step 3: Persistence (Storage)

To make our notes stick around after we close the app, we'll save them to a `notes.json` file. We use the standard `encoding/json` package for this.

```go
func saveNotes(notes []Note) error {
    data, _ := json.Marshal(notes)
    return os.WriteFile("notes.json", data, 0644)
}

func loadNotes() ([]Note, error) {
    data, err := os.ReadFile("notes.json")
    if err != nil {
        return []Note{}, nil // Return empty if file doesn't exist
    }
    var notes []Note
    err = json.Unmarshal(data, &notes)
    return notes, err
}
```

---

## Step 4: The Bubble Tea Model

Bubble Tea apps have three main parts:
1. **Model**: The state of your application.
2. **Update**: A function that handles events (key presses, etc.) and updates the state.
3. **View**: A function that renders the UI based on the state.

### The Model
Our model will track the list of notes and which one is currently selected.

```go
type model struct {
    notes    []Note
    cursor   int
    adding   bool
    input    textinput.Model
}
```

---

## Step 5: Handling Interaction (Update)

In the `Update` function, we listen for keys like `up`, `down`, `enter`, and `q`.

```go
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.String() {
        case "ctrl+c", "q":
            return m, tea.Quit
        case "up":
            if m.cursor > 0 { m.cursor-- }
        case "down":
            if m.cursor < len(m.notes)-1 { m.cursor++ }
        }
    }
    return m, nil
}
```

---

## Step 6: Rendering the UI (View)

This is where we use **Lip Gloss** to make things look pretty.

```go
func (m model) View() string {
    s := "Current Notes:\n\n"
    for i, note := range m.notes {
        cursor := " "
        if m.cursor == i {
            cursor = ">"
        }
        s += fmt.Sprintf("%s %s\n", cursor, note.Title)
    }
    s += "\n(q to quit)\n"
    return s
}
```

---

## Step 7: Putting it All Together

In the `main` function, we initialize the model and start the Bubble Tea program.

```go
func main() {
    p := tea.NewProgram(initialModel())
    if _, err := p.Run(); err != nil {
        fmt.Printf("Error: %v", err)
        os.Exit(1)
    }
}
```

---

## Next Steps
In the actual implementation, we'll add:
- **Rich Styling**: Using gradients and bold text.
- **Input Fields**: For typing new notes.
- **Deletions**: Using the `d` key.
