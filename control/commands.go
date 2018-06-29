package control

import "github.com/gastrid/team-bandit/robots"

type Command interface {
	Action(args []string, p *robots.Payload) (*robots.IncomingWebhook, error)
}

type add struct{}

type del struct{}

type show struct{}

type gen struct{}

type help struct{}

type bulk struct{}
