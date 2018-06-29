package control

import (
	"fmt"
	"strings"

	"github.com/gastrid/team-bandit/robots"
)

func (b bulk) Action(fields []string, p *robots.Payload) (*robots.IncomingWebhook, error) {

	response := &robots.IncomingWebhook{
		Domain:      p.TeamDomain,
		Channel:     p.ChannelID,
		Username:    "team-bandit",
		UnfurlLinks: true,
		Parse:       robots.ParseStyleFull,
	}

	Data.Mx.Lock()
	defer Data.Mx.Unlock()

	if len(fields) < 2 {
		return nil, fmt.Errorf("Not enough arguments here!")
	}

	fields = fields[1:len(fields)]

	for _, field := range fields {
		nameDept := strings.Split(field, ":")
		if len(nameDept) != 2 {
			return nil, fmt.Errorf("Are you sure you've written all entries in the format personName:departmentName ?")
		}
		name := nameDept[0]
		dept := nameDept[1]

		_, err := addDepartment(dept)
		if err != nil && err.Error() != fmt.Sprint("This department exists already.") {
			return nil, err
		}

		_, err = addPerson(name, dept)
		if err != nil {
			return nil, err
		}
	}

	response.Text = "Done! All of these people and departments have been added!"
	return response, nil
}
