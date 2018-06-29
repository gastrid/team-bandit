package control

import (
	"encoding/json"
	"log"
)

type ButtonResponse struct {
	CallbackId string   `json:"callback_id"`
	Token      string   `json:"token"`
	Actions    []Action `json"actions"`
}

type Action struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

func addToMatches(person *Person) {
	for i, _ := range Data.People {
		p := Data.People[i]
		match := &Match{
			Match: [2]*Person{p, person},
		}
		if p.Department == person.Department {
			match.Score = 5
		} else {
			match.Score = 0
		}
		Data.People[i].Score = p.Score + match.Score
		person.Score = person.Score + match.Score

		Data.Matches = append(Data.Matches, match)

	}
}

func delFromMatches(person *Person) {
	if len(Data.Matches) != 0 {
		var delMatch Matches
		for _, m := range Data.Matches {
			// If the first or the second person in the match is the given person.
			if m.Match[0].Name == person.Name || m.Match[1].Name == person.Name {
				delMatch = append(delMatch, m)
			}
		}

		for i := 0; i < len(Data.Matches); i++ {
			match := Data.Matches[i]
			for _, d := range delMatch {
				if match == d {
					Data.Matches = append(Data.Matches[:i], Data.Matches[i+1:]...)
					i-- // Important: decrease index
					break
				}
			}
		}

	}
}

func persistLoad() error {
	deptJson, err := json.Marshal(Data.Departments)
	if err != nil {
		return err
	}

	// Saving everything in redis
	_, err = Data.DB.Exec(`INSERT INTO stores AS s (storetype, value) VALUES ('departments', $1) ON CONFLICT (storetype) DO UPDATE SET value = $1 WHERE s.storetype = 'departments'`, deptJson)
	if err != nil {
		return err
	}

	peopleJson, err := json.Marshal(Data.People)

	_, err = Data.DB.Exec(`INSERT INTO stores AS s (storetype, value) VALUES ('people', $1) ON CONFLICT (storetype) DO UPDATE SET value = $1 WHERE s.storetype = 'people'`, peopleJson)
	if err != nil {
		return err
	}

	matchesJson, err := json.Marshal(Data.Matches)

	_, err = Data.DB.Exec(`INSERT INTO stores AS s (storetype, value) VALUES ('matches', $1) ON CONFLICT (storetype) DO UPDATE SET value = $1 WHERE s.storetype = 'matches'`, matchesJson)
	if err != nil {
		return err
	}

	return nil
}

func (d *Datastruct) PersistTeams() {
	for _, t := range d.Teams {
		l := len(t.Members)
		for i := 0; i < l-1; i++ {
			for j := i + 1; j < l; j++ {
				firstPers := t.Members[i]
				secondPers := t.Members[j]
				match := getMatch(firstPers.Name, secondPers.Name)
				if firstPers.Department == secondPers.Department {
					firstPers.Score += 2
					secondPers.Score += 2
					match.Score += 2
				} else {
					firstPers.Score++
					secondPers.Score++
					match.Score++
				}

			}
		}
	}

	persistLoad()

	log.Printf("[slackdebug] teams persisted")
}

func getMatch(firstName string, secondName string) *Match {
	for i := 0; i < len(Data.Matches); i++ {
		if (Data.Matches[i].Match[0].Name == firstName && Data.Matches[i].Match[1].Name == secondName) ||
			(Data.Matches[i].Match[1].Name == firstName && Data.Matches[i].Match[0].Name == secondName) {
			return Data.Matches[i]
		}
	}
	return &Match{}
}
