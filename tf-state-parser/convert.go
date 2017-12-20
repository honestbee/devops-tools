package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/google/go-github/github"
	log "github.com/sirupsen/logrus"
)

// currently all teams are at root module
// retrieve a map of id -> attributes
// but Tf state is missing slug, so map that in too!

// return map which wraps Resource ID to Resource Attributes for all data in a module
func getResourceMap(m map[string]resourceState) map[string]map[string]string {
	rMap := make(map[string]map[string]string)

	for k, r := range m {
		rMap[r.Primary.ID] = r.Primary.Attributes
		rMap[r.Primary.ID]["type"] = r.Type
		rMap[r.Primary.ID]["key"] = k
	}

	return rMap
}

// index github teams by ID
func mapTeams(teams []*github.Team) map[string]*github.Team {
	teamMap := make(map[string]*github.Team)

	for _, t := range teams {
		i := strconv.FormatInt(int64(*t.ID), 10)
		teamMap[i] = t
	}
	return teamMap
}

func writeImports(s *state, gitHub *GitHub) error {
	//get all teams from github (to get attributes missing in state)
	teams, err := gitHub.ListTeams()
	if err != nil {
		return err
	}

	log.Debugf("GitHub teams retrieved: %v\n", len(teams))

	//index github teams by team_id
	teamMap := mapTeams(teams)

	for i, m := range s.Modules {
		resourceMap := getResourceMap(m.Resources)

		for k, rm := range resourceMap {
			// root (0) only has the github_team resource type
			// everything else only has github_team_membership resource type
			if i == 0 {
				fmt.Printf("terraform state mv %v module.ghac.github_team.%v\n",
					rm["key"],
					*teamMap[k].Slug,
				)
				//import module.ghac.github_team.<slug> <id>
			} else {
				user := rm["username"]
				slug := *teamMap[rm["team_id"]].Slug

				//log.Debugf("key:%v", rm["key"])
				s := strings.Split(rm["key"], ".")
				var path string
				if len(s) == 3 {
					path = strings.Join(s[:2], ".") + "[" + s[2] + "]"
				} else {
					path = rm["key"]
				}

				fmt.Printf("terraform state mv module.%v.%v module.ghac.github_team_membership.%v-%v\n",
					m.Path[1],
					path,
					slug,
					user,
				)

				//import module.ghac.github_team_membership.<slug>-<user> <id>
				//fmt.Printf("terraform state import module.ghac.github_team_membership.%v-%v %v\n", slug, user, k)
			}
		}
	}
	return nil
}
