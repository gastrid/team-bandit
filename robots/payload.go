package robots

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type SlashCommand struct {
	Payload
	Command string `schema:"command"`
}

type Payload struct {
	Token           string  `schema:"token"json:"token"`
	TeamID          string  `schema:"team_id"json:"team_id"`
	TeamDomain      string  `schema:"team_domain,omitempty"json:"team_domain,omitempty"`
	ChannelID       string  `schema:"channel_id"json:"channel_id"`
	ChannelName     string  `schema:"channel_name"json:"channel_name"`
	Timestamp       float64 `schema:"timestamp,omitempty"json:"timestamp,omitempty"`
	UserID          string  `schema:"user_id"json:"user_id"`
	UserName        string  `schema:"user_name"json:"user_name"`
	Text            string  `schema:"text,omitempty"json:"text,omitempty"`
	Action          string  `schema:"action,omitempty"json:"action,omitempty"`
	ServiceID       string  `schema:"service_id,omitempty"json:"service_id,omitempty"`
	ResponseUrl     string  `schema:"response_url,omitempty"json:"response_url,omitempty"`
	BotID           string  `schema:"bot_id,omitempty"json:"bot_id,omitempty"`
	BotName         string  `schema:"bot_name,omitempty"json:"bot_name,omitempty"`
	Robot           string
	Actions         []Action        `json:"actions,omitempty"`
	CallbackId      string          `json:"callback_id,omitempty"`
	ActionTS        string          `json:"action_ts",omitempty`
	MessageTS       string          `json:"message_ts",omitempty`
	AttachmentId    string          `json:"attachment_id",omitempty`
	OriginalMessage IncomingWebhook `json:"original_message",omitempty`
	Team            Entity          `json:"team",omitempty`
	Channel         Entity          `json:"channel",omitempty`
	User            Entity          `json:"user",omitempty`
}

type Entity struct {
	Id     string
	Name   string
	Domain string
}

type EventResponse struct {
	Payload Payload `json:"payload"`
}

type OutgoingWebHookResponse struct {
	Text      string     `json:"text"`
	Parse     ParseStyle `json:"parse,omitempty"`
	LinkNames bool       `json:"link_names,omitempty"`
	Markdown  bool       `json:"mrkdwn,omitempty"`
}

type ParseStyle string

var (
	ParseStyleFull = ParseStyle("full")
	ParseStyleNone = ParseStyle("none")
)

type Message struct {
	Domain      string       `json:"domain"`
	Channel     string       `json:"channel"`
	Username    string       `json:"username"`
	Text        string       `json:"text"`
	IconEmoji   string       `json:"icon_emoji,omitempty"`
	IconURL     string       `json:"icon_url,omitempty"`
	Attachments []Attachment `json:"attachments,omitempty"`
	UnfurlLinks bool         `json:"unfurl_links,omitempty"`
	Parse       ParseStyle   `json:"parse,omitempty"`
	LinkNames   bool         `json:"link_names,omitempty"`
	Markdown    bool         `json:"mrkdwn,omitempty"`
}

type IncomingWebhook Message
type SlashCommandResponse Message

type Attachment struct {
	Fallback   string            `json:"fallback"`
	Pretext    string            `json:"pretext,omitempty"`
	Text       string            `json:"text,omitempty"`
	Color      string            `json:"color,omitempty"`
	Fields     []AttachmentField `json:"fields,omitempty"`
	MarkdownIn []MarkdownField   `json:"mrkdown_in,omitempty"`
	CallbackId string            `json:"callback_id,omitempty"`
	Actions    []Action          `json:"actions,omitempty"`
}

type Action struct {
	Name  string `json:"name"`
	Text  string `json:"text"`
	Type  string `json:"type"`
	Value string `json:"value"`
}

type MarkdownField string

var (
	MarkdownFieldPretext  = MarkdownField("pretext")
	MarkdownFieldText     = MarkdownField("text")
	MarkdownFieldTitle    = MarkdownField("title")
	MarkdownFieldFields   = MarkdownField("fields")
	MarkdownFieldFallback = MarkdownField("fallback")
)

type AttachmentField struct {
	Title string `json:"title"`
	Value string `json:"value"`
	Short bool   `json:"short,omitempty"`
}

// Send uses the IncomingWebhook API to post a message to a slack channel
func (i IncomingWebhook) Send() error {
	u := os.Getenv(fmt.Sprintf("%s_IN_URL", strings.ToUpper(i.Domain)))
	if u == "" {
		return fmt.Errorf("Slack Incoming Webhook URL not found for domain %s (check %s)", i.Domain, fmt.Sprintf("%s_IN_URL", strings.ToUpper(i.Domain)))
	}
	return Message(i).sendToUrl(u)
}

// Send a response to the ResponseUrl in the Payload
func (r SlashCommandResponse) Send(p *Payload) error {
	if p.ResponseUrl == "" {
		return fmt.Errorf("Empty ResponseUrl in Payload: %v", p)
	}
	return Message(r).sendToUrl(p.ResponseUrl)
}

func (i Message) sendToUrl(u string) error {
	if u == "" {
		return fmt.Errorf("Empty URL")
	}
	webhook, err := url.Parse(u)
	if err != nil {
		log.Printf("Error parsing URL \"%s\": %v", u, err)
		return err
	}

	p, err := json.Marshal(i)
	if err != nil {
		return err
	}

	data := url.Values{}
	data.Set("payload", string(p))

	webhook.RawQuery = data.Encode()
	resp, err := http.PostForm(webhook.String(), data)
	if resp.StatusCode != 200 {
		message := fmt.Sprintf("ERROR: Non-200 Response from Slack URL \"%s\": %s", u, resp.Status)
		log.Println(message)
	}
	return err
}
