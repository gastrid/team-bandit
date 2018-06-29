package teamalloc

import (
	"log"
	"strings"

	"fmt"

	"github.com/gastrid/team-bandit/control"
	"github.com/gastrid/team-bandit/robots"
)

type bot struct{}

func init() {
	s := &bot{}
	robots.RegisterRobot("team-bandit", s)
}

func (r bot) Run(p *robots.Payload) (slashCommandImmediateReturn string) {
	text := strings.TrimSpace(p.Text)

	resp, err := control.Data.Execute(text, p)
	if err != nil {
		log.Printf("[slackerror] %v", err)
		resp = &robots.IncomingWebhook{
			Domain:      p.TeamDomain,
			Channel:     p.ChannelID,
			Username:    "team-bandit",
			UnfurlLinks: true,
			Parse:       robots.ParseStyleFull,
			Text:        fmt.Sprintf("Oops, something went wrong: %v", err),
		}
	}

	resp.Send()

	return ""
}

func (r bot) Description() (description string) {
	return "This command will generate teams for you."
}
