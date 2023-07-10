package iniparser

import (
	"bufio"
	"github.com/pkg/errors"
	"os"
	"path/filepath"
	"io"
	"fmt"
	"strings"
)

// exported variable so it should start with uppercase
var (
	ErrInvalidFormat    = errors.New("invalid format ")
	ErrInvalidExtension = errors.New("file is not in the ini format or does not have a .ini extension")
	ErrKeyNotExist      = errors.New("key doesn't exist")
	ErrSectionNotExist  = errors.New("section doesn't exist")
)

type Section map[string]string

type Ini map[string]Section

type Parser struct {
	sections Ini
}

// NewParser returns new Parser
func NewParser() Parser {
	return Parser{sections: make(Ini)}
}

func (p *Parser) LoadFromString(content string) error {
	return p.LoadFromReader(bufio.NewReader(strings.NewReader(content)))
}

func (p *Parser) LoadFromFile(path string) error {

	fileExt := filepath.Ext(path)

	if fileExt != ".ini" {
		return ErrInvalidExtension
	}

	file, err := os.Open(path)

	if err != nil {
		return err
	}

	defer file.Close()

	return p.LoadFromReader(bufio.NewReader(file))
}

func (p *Parser) LoadFromReader(reader io.Reader) error {
	var currentSection string = ""
	scanner := bufio.NewScanner(reader)
	idx := 0
	for scanner.Scan() {
		idx++
		line := scanner.Text()

		line = strings.TrimSpace(line)
		// comment
		if len(line) == 0 || string(line[0]) == "#" {
			continue
		}

		// checking if the line in that format [section] and just exist one time in the line

		// new section
		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {

			sectionName := line[1 : len(line)-1]
			sectionName = strings.TrimSpace(sectionName)

			if len(sectionName) == 0 {
				return errors.Wrapf(ErrInvalidFormat, "invalid section at line %d", idx)
			}

			currentSection = sectionName

			p.sections[currentSection] = make(Section)
			continue

		}

		if strings.Contains(line, "=") && currentSection != "" {

			// key and value line
			keyAndValue := strings.Split(line, "=")
			key, value := keyAndValue[0], keyAndValue[1]
			key, value = strings.TrimSpace(key), strings.TrimSpace(value)

			if len(key) == 0 {
				return errors.Wrapf(ErrInvalidFormat, "invalid key at line %d", idx)
			}

			p.sections[currentSection][key] = value
		} else {
			return errors.Wrapf(ErrInvalidFormat, "invalid format at line %d", idx)
		}
	}
	return scanner.Err()
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
		p.sections[section_name] = make(Section)
	}
	p.sections[section_name][key] = value
}

func (p *Parser) String() string {

	configText := ""
	for section, configs := range p.sections {
		configText += fmt.Sprintf("[%s]\n", section)
		for key, value := range configs {
			configText += fmt.Sprintf("%s=%s\n", key, value)
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
