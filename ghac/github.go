package main

import (
	"io/ioutil"
	"strings"

	yaml "gopkg.in/yaml.v2"
)

type (
	// Team struct represents a GitHub team
	Team struct {
		ID          int               `yaml:"id,omitempty"`
		Name        string            `yaml:"name,omitempty"`
		Description string            `yaml:"description,omitempty"`
		Slug        string            `yaml:"slug,omitempty"`
		Namespace   string            `yaml:"namespace,omitempty"`
		Privacy     string            `yaml:"privacy,omitempty"`
		UserRoles   map[string]string `yaml:"user_roles"`
		Parent      *Team             `yaml:"parent,omitempty"`
		SlugPrefix  string
		SlugSuffix  string
	}
	// TeamFromYaml is a Team wrapper with custom UnmarshalYaml to fill computed fields
	TeamFromYaml Team
	// TeamList struct holds a list of GitHub teams
	TeamList struct {
		Teams []*TeamFromYaml `yaml:"teams"`
	}
)

// UnmarshalYAML implements custom unmarshal which adds computed fields
func (d *TeamFromYaml) UnmarshalYAML(unmarshal func(interface{}) error) error {
	if err := unmarshal((*Team)(d)); err != nil {
		return err
	}
	// compute additional fields for all teams
	sl := strings.Split(d.Slug, "-")
	d.SlugPrefix = strings.Join(sl[:len(sl)-1], "-")
	d.SlugSuffix = sl[len(sl)-1]

	return nil
}

func readTeams(yamlfile string) (*TeamList, error) {
	tl := TeamList{}
	data, err := ioutil.ReadFile(yamlfile)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal([]byte(data), &tl)
	if err != nil {
		return nil, err
	}

	return &tl, nil
}
