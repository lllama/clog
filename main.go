package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	errorStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("9"))
	titleStyle        = lipgloss.NewStyle().MarginLeft(2)
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
	quitTextStyle     = lipgloss.NewStyle().Margin(1, 0, 2, 4)
)

const listHeight = 15

type itemDelegate struct{}
type item string

func (i item) FilterValue() string { return string(i) }

func (d itemDelegate) Height() int                             { return 1 }
func (d itemDelegate) Spacing() int                            { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
	if !ok {
		return
	}

	str := fmt.Sprintf("%d. %s", index+1, i)

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return selectedItemStyle.Render("> " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(str))
}

type model struct {
	logGroups []string
	// cursor    int
	// logGroups list.Model
	// ready     bool
}

// func initialModel() model {
// 	defaultWidth := 20
// 	l := list.New([]list.Item{}, itemDelegate{}, defaultWidth, listHeight)
// 	l.Title = "What do you want for dinner?"
// 	l.SetShowStatusBar(false)
// 	l.SetFilteringEnabled(true)
// 	l.Styles.Title = titleStyle
// 	l.Styles.PaginationStyle = paginationStyle
// 	l.Styles.HelpStyle = helpStyle
//
// 	return model{cursor: 0, logGroups: l}
// }

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	//  	var (
	//  		cmd tea.Cmd
	//  	)
	//
	//  	switch msg := msg.(type) {
	//  	case tea.KeyMsg:
	//  		if m.logGroups.FilterState() == list.Filtering {
	//  			break
	//  		}
	//  		if msg.String() == "q" {
		// 			return m, tea.Quit
		// 		}
		// 	case tea.WindowSizeMsg:
		// 		m.logGroups.SetWidth(msg.Width)
		// 		m.logGroups.SetHeight(msg.Height)
		// 		return m, nil
		// 	}
		// 	m.logGroups, cmd = m.logGroups.Update(msg)
		//
		// 	return m, cmd
		//
	return m, nil
}

func (m model) View() string {
	return strings.Join(m.logGroups, "\n")
	// return m.logGroups.View()
}

func main() {

	model := model{}
	config, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("eu-west-1"))
	if err != nil {
		log.Fatalf("failed to load configuration, %v", err)
	}

	sts_client := sts.NewFromConfig(config)

	_, err = sts_client.GetCallerIdentity(context.Background(), &sts.GetCallerIdentityInput{})
	if err != nil {
		fmt.Println(errorStyle.Render(fmt.Sprintf("Bad AWS Credentials: %v ", err)))
		return
	}

	client := cloudwatchlogs.NewFromConfig(config)

	paginator := cloudwatchlogs.NewDescribeLogGroupsPaginator(client, &cloudwatchlogs.DescribeLogGroupsInput{})

	// groups := []list.Item{}

	for paginator.HasMorePages() {
		output, err := paginator.NextPage(context.Background())
		if err != nil {
			log.Fatalf("failed to list log groups, %v", err)
		}
		for _, logGroup := range output.LogGroups {
			model.logGroups = append(model.logGroups, string(*logGroup.LogGroupName))
			// groups = append(groups, item(string(*logGroup.LogGroupName)))
		}
	}

	// model.logGroups.SetItems(groups)

	btProg := tea.NewProgram(model)

	if _, err := btProg.Run(); err != nil {
		log.Fatalf("failed to start program, %v", err)
	}

}
