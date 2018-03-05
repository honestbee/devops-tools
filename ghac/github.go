package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
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
	log.Debugf("SlugPrefix %v - SlugSuffix %v", d.SlugPrefix, d.SlugSuffix)
	return nil
}

func findTeamsYaml(teamsDir string) ([]string, error) {
	fileList := []string{}
	err := filepath.Walk(teamsDir, func(path string, f os.FileInfo, err error) error {
		match, err := regexp.MatchString(".(yaml|yml)", path)
		if err != nil {
			log.Errorf("Error matching file name %q: %v\n", path, err)
			return err
		}
		if match == true {
			log.Debugf("Found match %q\n", path)
			fileList = append(fileList, path)
		}
		return nil
	})
	if err != nil {
		log.Errorf("Error walking path %q: %v\n", teamsDir, err)
		return nil, err
	}
	return fileList, nil
}

func makeTeams(teamsDir string) (*TeamList, error) {
	fileList, err := findTeamsYaml(teamsDir)
	if err != nil {
		return nil, err
	}
	tl := TeamList{}

	for _, file := range fileList {
		teamList, err := readTeams(file)
		if err != nil {
			return nil, err
		}
		tl.Teams = append(tl.Teams, teamList.Teams...)
	}
	return &tl, nil
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
