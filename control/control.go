package control

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	"encoding/json"

	"github.com/gastrid/team-bandit/robots"
	_ "github.com/lib/pq"
	"github.com/nlopes/slack"
)

// TODO:

var cmds = map[string]Command{
	"add":    add{},
	"delete": del{},
	"show":   show{},
	"gen":    gen{},
	"help":   help{},
	"bulk":   bulk{},
}

type Datastruct struct {
	Mx sync.Mutex
	People
	Departments []*Department
	Teams       []Team
	Matches     Matches
	DB          *sql.DB
	Me          string
	Locked      bool
	RTM         *slack.RTM
	CallbackId  string
	Fields      []string
}

var Data = Datastruct{}

var ResponseRobot = "response-robot"

func init() {

	connectClient()

	// Retrieve all the "tables"
	err := Data.getStore()
	if err != nil {
		fmt.Println(err)
	}

	return
}

func (d *Datastruct) Execute(cmd string, p *robots.Payload) (*robots.IncomingWebhook, error) {
	// Turn command string to slice
	Data.Fields = strings.Fields(cmd)

	if len(Data.Fields) > 0 {
		firstarg := Data.Fields[0]

		if firstarg == "hello" {
			return &robots.IncomingWebhook{
				Domain:      p.TeamDomain,
				Channel:     p.ChannelID,
				Username:    "team-bandit",
				UnfurlLinks: true,
				Parse:       robots.ParseStyleFull,
				Text:        "Hello! What can I do for you?",
			}, nil
		}

		if _, in := cmds[firstarg]; !in {
			return nil, fmt.Errorf("Oops, command doesn't exist")
		}

		// retrieve command
		command := cmds[firstarg]

		// send to the right function
		resp, err := command.Action(Data.Fields, p)
		if err != nil {
			return nil, err
		}

		return resp, nil
	}

	return nil, fmt.Errorf("Oops, not enough arguments here!")

}

type Scanner interface {
	Scan(...interface{}) error
}

// Opens the right file and reads the bytes into a struct
func (d *Datastruct) getStore() error {
	rows, err := Data.DB.Query(`SELECT storeType, value FROM stores`)
	if err != nil {
		return err
	}

	for rows.Next() {
		err := d.populateStores(rows)
		if err != nil {
			return err
		}
	}
	return nil

}

func (d *Datastruct) populateStores(sc Scanner) error {
	var storeType string
	var jsonString []byte

	sc.Scan(
		&storeType,
		&jsonString,
	)

	if storeType == "departments" {
		err := json.Unmarshal(jsonString, &d.Departments)
		if err != nil {
			return err
		}
	} else if storeType == "people" {
		err := json.Unmarshal(jsonString, &d.People)
		if err != nil {
			return err
		}
	} else if storeType == "matches" {
		err := json.Unmarshal(jsonString, &d.Matches)
		if err != nil {
			return err
		}
	} else {
		return fmt.Errorf("wrong row name: %s", storeType)
	}

	return nil

}

func connectClient() {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Panicf("Error opening database: %q", err)
	}

	Data.DB = db

	if _, err := Data.DB.Exec("CREATE TABLE IF NOT EXISTS stores (storetype varchar(40) NOT NULL UNIQUE, value bytea NOT NULL)"); err != nil {
		log.Panic(
			fmt.Sprintf("Error creating database table: %q", err))
		return
	}

}
