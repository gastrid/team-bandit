package responserobot

import (
	"fmt"

	"github.com/gastrid/team-bandit/control"
	"github.com/gastrid/team-bandit/robots"
)

type bot struct{}

func init() {
	s := &bot{}
	robots.RegisterRobot("response-robot", s)
}

func (r bot) Run(p *robots.Payload) (slashCommandImmediateReturn string) {

	control.Data.PersistTeams()
	resp := &robots.IncomingWebhook{}
	if p.Actions[0].Value == "yes" {
		resp = &robots.IncomingWebhook{
			Domain:      p.Team.Domain,
			Channel:     p.Channel.Id,
			Username:    "team-bandit",
			UnfurlLinks: true,
			Text:        control.MakeTeamConfirmationString(),
		}

	} else {
		resp = &robots.IncomingWebhook{
			Domain:      p.Team.Domain,
			Channel:     p.Channel.Id,
			Username:    "team-bandit",
			UnfurlLinks: true,
			Text:        fmt.Sprint("Ok, try again maybe?"),
		}
	}
	resp.Send()

	return ""
}

func (r bot) Description() (description string) {
	return "This command will confirm teams for you."
}
