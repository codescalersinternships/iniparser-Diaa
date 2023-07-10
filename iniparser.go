package iniparser

import (
	"github.com/pkg/errors"
	"os"
	"path/filepath"
	reg "regexp"
	"strings"
)

// exported variable so it should start with uppercase
var (
	ErrInvalidFormat    = errors.New("not valid format in line ")
	ErrInvalidExtension = errors.New("File is not in the INI format or does not have a .ini extension")
	ErrKeyNotExist      = errors.New("key doesn't exist")
	ErrSectionNotExist  = errors.New("Section doesn't exist")
)

type Section map[string]string

type Ini map[string]Section

type Parser struct {
	sections Ini
}

func NewParser() Parser {
	return Parser{sections: make(Ini)}
}

func (p *Parser) LoadFromString(content string) error {
	lines := strings.Split(content, "\n")
	var currentSection string = ""

	for idx, line := range lines {

		line = strings.TrimSpace(line)

		// comment
		if len(line) == 0 || string(line[0]) == "#" {
			continue
		}

		// checking if the line in that format [section] and just exist one time in the line
		regex := reg.MustCompile(`^\[[a-zA-Z\s]+\]$`)

		isMatched := regex.MatchString(line)

		// new section
		if isMatched {

			currentSection = line[1 : len(line)-1]

			p.sections[currentSection] = make(map[string]string)
			continue

		}

		if strings.Contains(line, "=") && currentSection != "" {

			// key and value line
			keyAndValue := strings.Split(line, "=")
			key, value := keyAndValue[0], keyAndValue[1]
			key, value = strings.TrimSpace(key), strings.TrimSpace(value)

			p.sections[currentSection][key] = value
		} else {
			return errors.Wrapf(ErrInvalidFormat, "%d", idx)
		}
	}
	return nil
}

func (p *Parser) LoadFromFile(path string) error {

	fileExt := filepath.Ext(path)

	if fileExt != ".ini" {
		return ErrInvalidExtension
	}

	data, err := os.ReadFile(path)

	if err != nil {
		return err
	}

	return p.LoadFromString(string(data))
}

func (p *Parser) GetSectionNames() []string {
	sections := make([]string, 0, len(p.sections))
	for section := range p.sections {
		sections = append(sections, section)
	}
	return sections
}

func (p *Parser) GetSections() Ini {
	return p.sections
}

func (p *Parser) Get(section_name, key string) (string, error) {

	_, ok := p.sections[section_name]
	if !ok {
		return "", ErrSectionNotExist
	}

	value, ok := p.sections[section_name][key]
	if !ok {
		return "", ErrKeyNotExist
	}

	return value, nil
}

func (p *Parser) Set(section_name, key, value string) {

	if p.sections[section_name] == nil {
		p.sections[section_name] = make(map[string]string)
	}
	p.sections[section_name][key] = value
}

func (p *Parser) String() string {

	configText := ""
	for section, configs := range p.sections {
		configText += "[" + section + "]\n"
		for key, value := range configs {
			configText += key + "=" + value + "\n"
		}
	}

	return configText
}

func (p *Parser) SaveToFile(path string) error {

	configString := p.String()
	stringBytes := []byte(configString)

	// 0644 is an octal code for access (admin: read and write, other users :read)
	return os.WriteFile(path, stringBytes, 0644)

}

// adding new case to test Get
// new test with 2 sections in one line []
