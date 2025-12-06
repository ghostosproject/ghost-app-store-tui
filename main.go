package main

import (
	"fmt"
	"ghost-app-store-tui/apps"
	"ghost-app-store-tui/dto"
	"io"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ghostosproject/ghost-app/model"
	gdto "github.com/ghostosproject/ghost-apps/dto"
	"github.com/ghostosproject/ghost-apps/tui/messages"
)

type GhostAppModel struct {
	height      int
	width       int
	apps        list.Model
	selectedApp Item
}

var GhostApps []dto.App

var ItemStyle = lipgloss.NewStyle().PaddingLeft(2)
var selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))

type ItemDelegate struct{}

func (d ItemDelegate) Height() int                             { return 3 }
func (d ItemDelegate) Spacing() int                            { return 1 }
func (d ItemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d ItemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(Item)
	if !ok {
		return
	}

	str := fmt.Sprintf("%s %s", i.icon, i.title)

	fn := ItemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return selectedItemStyle.Render(strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(str))
}

type Item struct {
	icon, title, id string
	versions        []dto.Version
}

func (i Item) Title() string       { return i.title }
func (i Item) Id() string          { return i.id }
func (i Item) FilterValue() string { return i.title }

func NewModel(apps []dto.App, width, height int) GhostAppModel {
	// apps := functions.GetApps()
	// apps := []model.App{}
	items := make([]list.Item, len(apps))
	for i, app := range apps {
		items[i] = Item{icon: app.Icon, id: app.Id, title: app.Name, versions: app.Versions}
	}
	del := ItemDelegate{}
	//list.NewDefaultDelegate()
	lst := list.New(items, del, 30, 30)
	lst.Title = "My Apps"
	lst.SetShowTitle(false)
	lst.SetFilteringEnabled(false)
	lst.SetShowStatusBar(false)
	lst.SetFilteringEnabled(false)
	lst.SetShowHelp(false)
	// apps to list

	// l.Styles.Title = titleStyle
	// lst.Styles.Title = lipgloss.NewStyle().Background(lipgloss.Color("#FFFFFF00")).Foreground(lipgloss.Color("#0a9f9a"))
	// lst.Styles.TitleBar = lipgloss.NewStyle().Background(lipgloss.Color("#FFFFFF00")).Foreground(lipgloss.Color("#0a9f9a"))
	m := GhostAppModel{apps: lst, height: height, width: width}
	return m
}

func (m GhostAppModel) Init() tea.Cmd {
	return nil
}

func (m GhostAppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "d":
			// download the app... add config to app_config
			data := gdto.ConfigApp{
				Id:      m.selectedApp.id,
				Name:    m.selectedApp.title,
				Icon:    m.selectedApp.icon,
				Version: m.selectedApp.versions[0].Version,
			}
			cmd = func() tea.Msg {
				return messages.ChildMessage{Target: "update-config", Data: data}
			}
			cmds = append(cmds, cmd)

		case "enter":
			i := m.apps.SelectedItem()
			m.selectedApp = i.(Item)
		}
	}

	m.apps, cmd = m.apps.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func getConfig() {
	bts, err := os.ReadFile("usr-config.json")
	if err != nil {
		panic(err)
	}

	err = os.WriteFile("testing.json", bts, 0644)
	if err != nil {
		panic(err)
	}
}

func (m GhostAppModel) View() string {
	sidePanel := lipgloss.NewStyle().
		Align(lipgloss.Left).
		Width((m.width*3/4*1/3)-5).
		Height(m.height-4-2).
		Border(lipgloss.NormalBorder(), false, true, false, false).
		BorderForeground(lipgloss.Color("#0a1")).
		Render(m.apps.View())

	appInfoData := fmt.Sprintf("%s %s\nThis is the app description\nGithub: ghostosproject/ghost-app-example\n\nVersions\n\nv0.2.6\tOctober 3, 2025 11:55pm\tlatest\t*\nv0.2.1\tOctober 1, 2025 10:42pm", "👻", "Ghost App")

	if m.selectedApp.id != "" {
		appInfoData = fmt.Sprintf("%s %s %s\n\nVersions\n\n", m.selectedApp.icon, m.selectedApp.title, m.selectedApp.id)
		for _, v := range m.selectedApp.versions {
			appInfoData += fmt.Sprintf("%s\n", v.Version)
		}
	}

	bottomAppInfoData := "App Controls\n\nd - download app\narrows - navigate list\nv - select app version"

	bottomAppInfoView := lipgloss.NewStyle().
		Align(lipgloss.Left).
		Width((m.width*3/4*2/3)-10).
		Height(m.height*1/3-5).
		Border(lipgloss.NormalBorder(), true, true, true, true).
		Background(lipgloss.Color("#ff19")).
		BorderForeground(lipgloss.Color("#ff19")).
		Render(bottomAppInfoData)

	appInfoView := lipgloss.NewStyle().
		Align(lipgloss.Left).
		Width((m.width*3/4*2/3)-10).
		Height(m.height*2/3-5).
		Border(lipgloss.NormalBorder(), true, true, true, true).
		Background(lipgloss.Color("#ff19")).
		BorderForeground(lipgloss.Color("#ff19")).
		Render(appInfoData)
	// homePanel := lipgloss.NewStyle().
	// 	Align(lipgloss.Left).
	// 	Width((m.width*3/4)-5).
	// 	Height(15).
	// 	Border(lipgloss.NormalBorder(), false, true, false, false).
	// 	BorderForeground(lipgloss.Color("#0a1")).
	// 	Render(m.apps.View())
	content := lipgloss.NewStyle().
		Width((m.width * 3 / 4 * 2 / 3)).
		// accommodate header and footer
		Height(m.height-5).
		Align(lipgloss.Center, lipgloss.Center).
		BorderForeground(lipgloss.Color("#0a1")).
		Background(lipgloss.Color("#ff19aa")).
		Render(lipgloss.JoinVertical(lipgloss.Top, appInfoView, bottomAppInfoView))
	return lipgloss.JoinHorizontal(lipgloss.Left, sidePanel, content)
	//return "This is the Ghost App Model"
	// return fmt.Sprintf("Width: %d, Height: %d", m.width, m.height)
}

func MakeApp(width, height int) (model.GhostApp, error) {
	GhostApps = apps.GetApps()
	GhostApps = append(GhostApps, dto.App{
		Id:   "testId",
		Name: "Ghost Messenger",
		Icon: "💬",
	})
	m := NewModel(GhostApps, height, width)
	return m, nil
}

// list of apps
// filter the apps (search for app)

// app details list
// link to github
// options...
// 	- bookmark
//	- download

// app details
// 	- name
// 	- recent version
// 	- link to release on github (may be our own platform )

// updates are made in the config file for the user

// list of apps
// enter takes you to the apps page
// description... download options... versions
