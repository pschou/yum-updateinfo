package main

import (
	"io/ioutil"
	"log"
	"regexp"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Packages []struct {
		Match    string `yaml:"Match"`
		From     string `yaml:"From"`
		RefTitle string `yaml:"RefTitle"`
		RefURL   string `yaml:"RefURL"`
		reg      *regexp.Regexp
	} `yaml:"packages"`
}

func (c *Config) getConf(confFile string) {
	yamlFile, err := ioutil.ReadFile(confFile)
	if err != nil {
		log.Printf("config file read error: %v ", err)
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}
	for i, m := range c.Packages {
		c.Packages[i].reg = regexp.MustCompile(m.Match)
	}
}
