package main

import (
	"io/ioutil"

	log "github.com/sirupsen/logrus"
	yaml "gopkg.in/yaml.v2"
)

type (
	// Team struct represents a GitHub team
	Team struct {
		ID          int               `yaml:"id,omitempty"`
		Name        string            `yaml:"name,omitempty"`
		Description string            `yaml:"description,omitempty"`
		Slug        string            `yaml:"slug,omitempty"`
		Privacy     string            `yaml:"privacy,omitempty"`
		UserRoles   map[string]string `yaml:"user_roles"`
		Parent      *Team             `yaml:"parent,omitempty"`
	}
	// TeamList struct holds a list of GitHub teams
	TeamList struct {
		Teams []Team `yaml:"teams"`
	}
)

func readTeams(yamlfile string) (*TeamList, error) {
	// tl := TeamList{
	// 	Teams: []Team{
	// 		Team{
	// 			ID:          1,
	// 			Name:        "Test team1",
	// 			Slug:        "test-team1",
	// 			Description: "test team1",
	// 			UserRoles: map[string]string{
	// 				"admin1": "maintainer",
	// 				"admin2": "maintainer",
	// 				"user1":  "member",
	// 			},
	// 		},
	// 		Team{
	// 			ID:          2,
	// 			Name:        "Test team2",
	// 			Slug:        "test-team2",
	// 			Description: "test team2",
	// 			UserRoles: map[string]string{
	// 				"admin1": "maintainer",
	// 				"admin2": "maintainer",
	// 				"user2":  "member",
	// 			},
	// 		},
	// 	},
	// }
	// fmt.Printf("--- tl:\n%v\n\n", tl)

	// d, err := yaml.Marshal(&tl)
	// if err != nil {
	// 	log.Fatalf("error: %v", err)
	// }
	// fmt.Printf("--- tl dump:\n%s\n\n", string(d))

	tl := TeamList{}
	data, err := ioutil.ReadFile(yamlfile)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal([]byte(data), &tl)
	if err != nil {
		return nil, err
	}
	log.Debugf("--- tl:\n%v\n\n", tl)

	return &tl, nil
}
