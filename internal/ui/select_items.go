package ui

import (
	"fmt"
	"io"
	"os/exec"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/pirpedro/dev-gadgets/internal/catalog"
)

var draculaBg = lipgloss.Color("#282a36")
var draculaFg = lipgloss.Color("#f8f8f2")
var draculaPurple = lipgloss.Color("#bd93f9")

// additional colors can be added later if needed

var welcomeArt = `
  _____          _____   _____ ______ _______ _____
 / ____|   /\   |  __ \ / ____|  ____|__   __/ ____|
| |  __   /  \  | |  | | |  __| |__     | | | (___
| | |_ | / /\ \ | |  | | | |_ |  __|    | |  \___ \
| |__| |/ ____ \| |__| | |__| | |____   | |  ____) |
 \_____/_/    \_|_____/ \_____|______|  |_| |_____/

Welcome to dev gadgets! A bunch of dev extensions to make your development more organized and productive.
`

// Item para seleção
type SelectItem struct {
	ID        string
	Name      string
	Desc      string
	Installed bool
}

func (i SelectItem) Title() string       { return i.Name }
func (i SelectItem) Description() string { return i.Desc }
func (i SelectItem) FilterValue() string { return i.Name }

type extraKeys struct {
	Select      key.Binding
	SelectAll   key.Binding
	DeselectAll key.Binding
}

var ek = extraKeys{
	Select:      key.NewBinding(key.WithKeys("space"), key.WithHelp("space", "select")),
	SelectAll:   key.NewBinding(key.WithKeys("a"), key.WithHelp("a", "select all")),
	DeselectAll: key.NewBinding(key.WithKeys("d"), key.WithHelp("d", "deselect all")),
}

// Modelo principal
type SelectItemsModel struct {
	list       list.Model
	progress   progress.Model
	keys       extraKeys
	installing bool
	selected   map[string]bool
	items      []SelectItem
	step       int
	steps      int
	msg        string
}

// Retorna os IDs dos itens selecionados
func (m SelectItemsModel) SelectedIDs() []string {
	var ids []string
	for id, sel := range m.selected {
		if sel {
			ids = append(ids, id)
		}
	}
	return ids
}

// Delegate customizado para múltipla seleção
type selectDelegate struct {
	selected map[string]bool
}

func (d selectDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	if item == nil {
		return
	}
	it, ok := item.(SelectItem)
	if !ok {
		return
	}
	mark := "[ ]"
	if it.Installed {
		mark = "[✓ installed]"
	} else if d.selected[it.ID] {
		mark = "[x]"
	}
	style := lipgloss.NewStyle().Foreground(draculaFg)
	if m.Index() == index {
		style = style.Background(draculaPurple).Bold(true)
	}
	if it.Installed {
		style = style.Faint(true)
	}
	fmt.Fprintln(w, style.Render(fmt.Sprintf("%s %s — %s", mark, it.Name, it.Desc)))
}

func (d selectDelegate) Height() int                               { return 1 }
func (d selectDelegate) Spacing() int                              { return 0 }
func (d selectDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }

func NewSelectItemsModel(items []SelectItem) SelectItemsModel {
	selected := map[string]bool{}
	delegate := selectDelegate{selected: selected}
	height := len(items) + 5
	if height > 20 {
		height = 20
	}
	// populate list with items to avoid nil entries
	listItems := make([]list.Item, len(items))
	for i := range items {
		listItems[i] = items[i]
	}
	l := list.New(listItems, delegate, 0, height)
	l.Title = "Select tools to install"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = lipgloss.NewStyle().Foreground(draculaPurple).Background(draculaBg).Bold(true)
	l.SetShowHelp(true)

	l.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{
			ek.Select,
			ek.SelectAll,
			ek.DeselectAll}
	}

	return SelectItemsModel{
		list:     l,
		progress: progress.New(progress.WithDefaultGradient()),
		keys:     ek,
		selected: selected,
		items:    items,
		steps:    len(items),
	}
}

func (m SelectItemsModel) Init() tea.Cmd {
	return nil
}

func (m SelectItemsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc":
			return m, tea.Quit
		case "up", "k":
			m.list.CursorUp()
		case "down", "j":
			m.list.CursorDown()
		case " ":
			item, ok := m.list.SelectedItem().(SelectItem)
			if ok && !item.Installed {
				m.selected[item.ID] = !m.selected[item.ID]
			}
		case "a":
			for _, it := range m.items {
				if !it.Installed {
					m.selected[it.ID] = true
				}
			}
		case "d":
			for _, it := range m.items {
				if !it.Installed {
					m.selected[it.ID] = false
				}
			}
		case "enter":
			// Start installation with progress bar
			m.installing = true
			m.step = 0
			m.steps = len(m.selected)
			m.msg = "Installing..."
			return m, m.nextStep()
		}
	case progress.FrameMsg:
		if m.installing && m.step < m.steps {
			m.step++
			m.msg = fmt.Sprintf("Installing %d/%d...", m.step, m.steps)
			return m, m.nextStep()
		} else if m.installing {
			m.msg = "Installation complete!"
			return m, tea.Quit
		}
	}
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

// Simula próximo passo da instalação (pode integrar com lógica real)
func (m *SelectItemsModel) nextStep() tea.Cmd {
	return func() tea.Msg {
		time.Sleep(500 * time.Millisecond)
		return progress.FrameMsg{}
	}
}

func (m SelectItemsModel) View() string {
	style := lipgloss.NewStyle().Foreground(draculaFg).Background(draculaBg)
	out := style.Render(welcomeArt) + "\n"
	out += m.list.View() + "\n"
	out += "\n"
	if m.installing {
		percent := 0.0
		if m.steps > 0 {
			percent = float64(m.step) / float64(m.steps)
		}
		out += m.progress.ViewAs(percent) + "\n"
		out += m.msg + "\n"
	}
	// Dynamic dependency check
	out += "\nDetected dependencies:\n"
	out += checkDepsView(m.items)
	return out
}

// Exibe dependências detectadas
func checkDepsView(items []SelectItem) string {
	var out string
	python := hasBin("python3") || hasBin("python")
	uv := hasBin("uv")
	pipx := hasBin("pipx")
	node := hasBin("node")
	npm := hasBin("npm")
	volta := hasBin("volta")
	if python {
		out += "- Python detected\n"
	} else {
		out += "- Python NOT detected\n"
	}
	if pipx {
		out += "- pipx detected\n"
	} else {
		out += "- pipx NOT detected\n"
	}
	if uv {
		out += "- uv detected\n"
	} else {
		out += "- uv NOT detected\n"
	}
	if node {
		out += "- Node detected\n"
	} else {
		out += "- Node NOT detected\n"
	}
	if npm {
		out += "- npm detected\n"
	} else {
		out += "- npm NOT detected\n"
	}
	if volta {
		out += "- volta detected\n"
	} else {
		out += "- volta NOT detected\n"
	}
	return out
}

func hasBin(bin string) bool {
	_, err := exec.LookPath(bin)
	return err == nil
}

// Função utilitária para checar se item está instalado
func IsInstalled(it catalog.Item) bool {
	if it.Verify == "" {
		return false
	}
	parts := strings.Fields(it.Verify)
	if len(parts) == 0 {
		return false
	}
	cmd := exec.Command(parts[0], parts[1:]...)
	return cmd.Run() == nil
}
