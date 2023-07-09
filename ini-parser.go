package main

import (
	"errors"
	"os"
	reg "regexp"
	"strings"
)

type Dictionary map[string]map[string]string

var ErrNotValidSyntax = errors.New("Not Valid Syntax")


type Config struct {
	iniConfigs Dictionary
}

func (c *Config) InitializeMap() {
	c.iniConfigs = make(Dictionary)
}

func (config *Config) LoadFromString(configs string) error {
	lines := strings.Split(configs, "\n")
	var lastSection string = ""

	for _, line := range lines {

		line = strings.Trim(line, " ")

		// comment
		if len(line) == 0 || string(line[0]) == "#" {
			continue
		}

		regex, err := reg.Compile(`\[[^\[\]]*\]`)

		if err != nil {
			return err
		}

		section := regex.FindStringSubmatch(line)

		// new section
		if len(section) > 0 {

			// getting section name
			lastSection = strings.ReplaceAll(section[0], "[", "")
			lastSection = strings.ReplaceAll(lastSection, "]", "")

			lastSection = strings.Trim(lastSection, " ")

			config.iniConfigs[lastSection] = make(map[string]string)

		} else if strings.Contains(line, "=") && lastSection != "" {

			// key and value line
			keyAndValue := strings.Split(line, "=")
			key, value := keyAndValue[0], keyAndValue[1]
			key, value = strings.Trim(key, " "), strings.Trim(value, " ")

			config.iniConfigs[lastSection][key] = value
		} else {
			return ErrNotValidSyntax
		}
	}
	return nil
}

func (config *Config) LoadFromFile(path string) error {
	data, err := os.ReadFile(path)

	if err != nil {
		return err
	}

	err = config.LoadFromString(string(data))

	if err != nil {
		return err
	}
	return nil
}
func (config *Config) GetSectionNames() [] string {
	sections := make([]string, 0, len(config.iniConfigs))
	for section, _ := range config.iniConfigs {
		sections = append(sections, section)
	}
	return sections
}

func (config *Config) GetSections() map[string]map[string]string {
	return config.iniConfigs
}

func (config *Config) Get(section_name, key string) string {
	return config.iniConfigs[section_name][key]
}

func (config *Config) Set(section_name, key, value string) {

	// checking if the section doesn't exist (new section)
	if config.iniConfigs[section_name] == nil {
		config.iniConfigs[section_name] = make(map[string]string)
	}
	config.iniConfigs[section_name][key] = value
}

func (config *Config) ToString() string {

	configText := ""
	for section, configs := range config.iniConfigs {
		configText += "[" + section + "]\n"
		for key, value := range configs {
			configText += key + "=" + value + "\n"
		}
	}

	return configText
}

func (config *Config) SaveToFile(path string) error {

	configString := config.ToString()
	stringBytes := []byte(configString)

	// 0644 is an octal code for access (admin: read and write, other users :read)
	err := os.WriteFile(path, stringBytes, 0644)
	if err != nil {
		return err
	}

	return nil
}
