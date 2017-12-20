package main

import (
	"io/ioutil"
	"strings"

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
		SlugPrefix  string
		SlugSuffix  string
	}
	// TeamListFromYaml is a TeamList wrapper with custom UnmarshalYaml to fill computed fields
	TeamListFromYaml TeamList
	// TeamList struct holds a list of GitHub teams
	TeamList struct {
		Teams []*Team `yaml:"teams"`
	}
)

// UnmarshalYAML implements how to unmarshal while adding computed fields
func (d *TeamListFromYaml) UnmarshalYAML(unmarshal func(interface{}) error) error {
	if err := unmarshal((*TeamList)(d)); err != nil {
		return err
	}
	// compute additional fields for all teams
	for _, t := range d.Teams {
		sl := strings.Split(t.Slug, "-")
		t.SlugPrefix = strings.Join(sl[:len(sl)-1], "-")
		t.SlugSuffix = sl[len(sl)-1]
		log.Debugf("SlugPrefix %v - SlugSuffix %v", t.SlugPrefix, t.SlugSuffix)
	}
	return nil
}

func readTeams(yamlfile string) (*TeamList, error) {
	tl := TeamListFromYaml{}
	data, err := ioutil.ReadFile(yamlfile)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal([]byte(data), &tl)
	if err != nil {
		return nil, err
	}
	log.Debugf("--- tl:\n%v\n\n", tl)

	return (*TeamList)(&tl), nil
}
