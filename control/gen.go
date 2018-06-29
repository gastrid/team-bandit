package control

import (
	"fmt"
	"math"
	"sort"
	"strconv"
	"time"

	"github.com/gastrid/team-bandit/robots"
	uuid "github.com/nu7hatch/gouuid"
)

func (g gen) Action(fields []string, p *robots.Payload) (*robots.IncomingWebhook, error) {

	// If there are fewer than 2 argments, there are not enough arguments.
	if len(fields) >= 2 {
		teamsize, err := strconv.Atoi(fields[1])
		if err != nil {
			return nil, err
		}

		var flagNames []string
		var without bool
		if len(fields) > 2 {
			if fields[2] == "-without" {
				flagNames = fields[3:]
				without = true
			} else if fields[2] == "-with" {
				flagNames = fields[3:]
				without = false
			} else {
				return nil, fmt.Errorf("Are you sure you entered the right arguments?")
			}
		}

		Data.Teams = getTeams(teamsize, flagNames, without)

		callbackId, err := uuid.NewV4()
		if err != nil {
			return nil, err
		}

		Data.CallbackId = callbackId.String()

		teamString := makeTeamSuggestionString()

		response := &robots.IncomingWebhook{
			Domain:      p.TeamDomain,
			Channel:     p.ChannelID,
			Username:    "team-bandit",
			UnfurlLinks: true,
			Parse:       robots.ParseStyleFull,
			Text:        teamString,
			Attachments: []robots.Attachment{
				robots.Attachment{
					Text:     "Shall I persist them?",
					Fallback: "Looks like you don't have a choice",
					Actions: []robots.Action{
						robots.Action{
							Name:  "team",
							Text:  "Yes",
							Value: "yes",
							Type:  "button",
						},
						robots.Action{
							Name:  "team",
							Text:  "No",
							Value: "no",
							Type:  "button",
						},
					},
					CallbackId: Data.CallbackId,
				},
			},
		}

		go unlockData()

		return response, nil

	}

	return nil, fmt.Errorf("Wrong number of arguments!")

}

func makeTeamSuggestionString() string {
	teamString := "Here are the teams as generated, do you like them? "

	for _, v := range Data.Teams {
		teamString += fmt.Sprintf(" \n [ ")
		for _, m := range v.Members {
			teamString += fmt.Sprintf("%s, ", m.Name)
		}
		teamString = teamString[:len(teamString)-2]
		teamString += fmt.Sprintf(" ] \n")
	}

	return teamString
}

func MakeTeamConfirmationString() string {
	teamString := "You have confirmed these teams:"

	for _, v := range Data.Teams {
		teamString += fmt.Sprintf(" \n [ ")
		for _, m := range v.Members {
			teamString += fmt.Sprintf("%s, ", m.Name)
		}
		teamString = teamString[:len(teamString)-2]
		teamString += fmt.Sprintf(" ] \n")
	}

	return teamString
}

func unlockData() {
	time.Sleep(time.Minute * 5)
	Data.Mx.Lock()
	Data.Locked = false
	Data.Mx.Unlock()
}

// with is 'false', without is true
func getTeams(teamSize int, absentees []string, without bool) []Team {
	Data.Mx.Lock()
	defer Data.Mx.Unlock()
	// Get a slice of all the people in the order of the person with the highest score first
	orderedPeople := orderPeople(absentees, without)
	sort.Sort(sort.Reverse(Data.Matches))

	// Get the number of teams by dividing the number of people by team size.
	teamNumber := int(math.Ceil(float64(len(orderedPeople)) / float64(teamSize)))

	teams := make([]Team, teamNumber)

	// We're iterating per team then per row, such that all first lines of teams are filled
	// before the second lines are.
	for i := 0; i < teamSize; i++ {
		for j := 0; j < teamNumber; j++ {
			if i == 0 {
				// If we're in the first row, we need to create the team first, and we put in the person with the higest
				// score.
				teams[j] = Team{
					Members: []*Person{orderedPeople[0]},
				}
				orderedPeople = orderedPeople[1:]
			} else {
				// In subsequent rows, we implement the logic to get matching teammates.
				next, index, nextScore, err := getMatchingPerson(teams[j].Members, orderedPeople)
				if err != nil && err.Error() != "No more leaders here!" {
					fmt.Println(err)
				}
				if err == nil {
					teams[j].Members = append(teams[j].Members, next)
					teams[j].Score += nextScore
					orderedPeople = append(orderedPeople[:index], orderedPeople[index+1:]...)
				}
			}
		}
	}

	Data.Locked = true

	return teams
}

func getMatchingPerson(array []*Person, orderedPeople []*Person) (*Person, int, int, error) {
	var leaderboard Leaderboard
	for k, p := range orderedPeople {
		if personNotSelected(array, p) {
			var personTotal int
			for _, m := range Data.Matches {
				if doesMatch(m, array, p) {
					personTotal = m.Score
					if m.Match[0].Department == m.Match[0].Department {
						personTotal += 2
					} else {
						personTotal++
					}
				}
			}
			leader := &Leader{
				Person:     p,
				TotalScore: personTotal,
				Index:      k,
			}
			leaderboard = append(leaderboard, leader)
		}
	}
	sort.Sort(leaderboard)
	if len(leaderboard) != 0 {
		return leaderboard[0].Person, leaderboard[0].Index, leaderboard[0].TotalScore, nil
	} else {
		nothing := &Person{}
		return nothing, 0, 0, fmt.Errorf("No more leaders here!")
	}

}

func orderPeople(absentees []string, without bool) []*Person {
	var slice People
	for _, k := range Data.People {
		absent := !without
		for _, a := range absentees {
			if k.Name == a {
				absent = without
			}
		}
		if absent == false {
			slice = append(slice, k)
		}
	}

	sort.Sort(sort.Reverse(slice))

	return slice
}

func personNotSelected(array []*Person, p *Person) bool {
	for _, person := range array {
		if p == person {
			return false
		}
	}
	return true
}

func doesMatch(match *Match, array []*Person, p *Person) bool {
	if match.Match[0].Name == p.Name {
		for _, person := range array {
			if match.Match[1].Name == person.Name {
				return true
			}
		}
	}

	if match.Match[1].Name == p.Name {
		for _, person := range array {
			if match.Match[0].Name == person.Name {
				return true
			}
		}
	}
	return false
}
