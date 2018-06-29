package control

import (
	"fmt"
	"sort"

	"github.com/gastrid/team-bandit/robots"
)

func (s show) Action(fields []string, p *robots.Payload) (*robots.IncomingWebhook, error) {

	response := &robots.IncomingWebhook{
		Domain:      p.TeamDomain,
		Channel:     p.ChannelID,
		Username:    "team-bandit",
		UnfurlLinks: true,
		Parse:       robots.ParseStyleFull,
		Text:        "Done",
	}

	if len(fields) == 2 {
		switch fields[1] {
		case "people":
			showPeople(response)
		case "departments":
			showDepartment(response)
		case "matches":
			showMatches(response)
		default:
			return nil, fmt.Errorf("That's not a listable type..")
		}

		return response, nil
	}

	return nil, fmt.Errorf("Wrong number of arguments!")

}

func showMatches(i *robots.IncomingWebhook) {
	sort.Sort(sort.Reverse(Data.Matches))
	showString := "Here is a list of this group's matches: \n"
	for _, m := range Data.Matches {
		showString += fmt.Sprintf("%s - %s : %v \n", m.Match[0].Name, m.Match[1].Name, m.Score)
	}

	i.Text = showString

}

func showPeople(i *robots.IncomingWebhook) {
	showString := "Here is a list of this group's people: \n"
	for _, p := range Data.People {
		showString += fmt.Sprintf(" - %s (%s)\n", p.Name, p.Department.Name)
	}

	i.Text = showString

}

func showDepartment(i *robots.IncomingWebhook) {
	showString := "Here is a list of this group's departments: \n"
	for _, d := range Data.Departments {
		showString += fmt.Sprintf("- %s \n", d.Name)
	}

	i.Text = showString
}
