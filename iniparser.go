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

// NewParser returns a new Parser object.
func NewParser() Parser {
	return Parser{sections: make(Ini)}
}

// LoadFromString loads INI data from a string.
func (p *Parser) LoadFromString(content string) error {
	return p.LoadFromReader(bufio.NewReader(strings.NewReader(content)))
}

// LoadFromFile loads INI data from a file.
func (p *Parser) LoadFromFile(path string) error {

	// Check file extension
	fileExt := filepath.Ext(path)
	if fileExt != ".ini" {
		return ErrInvalidExtension
	}

	// Open file
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	return p.LoadFromReader(bufio.NewReader(file))
}

// LoadFromReader loads INI data from an io.Reader object.
func (p *Parser) LoadFromReader(reader io.Reader) error {
	var currentSection string = ""
	scanner := bufio.NewScanner(reader)
	idx := 0
	for scanner.Scan() {
		idx++
		line := scanner.Text()

		line = strings.TrimSpace(line)

		// Skip empty lines and comments
		if len(line) == 0 || string(line[0]) == "#" {
			continue
		}

		// Check for section header
		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {

			// Extract section name
			sectionName := line[1 : len(line)-1]
			sectionName = strings.TrimSpace(sectionName)

			if len(sectionName) == 0 {
				return errors.Wrapf(ErrInvalidFormat, "invalid section at line %d", idx)
			}

			// Start new section
			currentSection = sectionName
			p.sections[currentSection] = make(Section)
			continue

		}

		// Check for key-value pairs
		if strings.Contains(line, "=") && currentSection != "" {

			// Parse key and value
			keyAndValue := strings.Split(line, "=")
			key, value := keyAndValue[0], keyAndValue[1]
			key, value = strings.TrimSpace(key), strings.TrimSpace(value)

			if len(key) == 0 {
				return errors.Wrapf(ErrInvalidFormat, "invalid key at line %d", idx)
			}

			// Add key-value pair to current section
			p.sections[currentSection][key] = value
		} else {
			return errors.Wrapf(ErrInvalidFormat, "invalid format at line %d", idx)
		}
	}
	return scanner.Err()
}
// GetSectionNames returns the names of all sections in the INI data.
func (p *Parser) GetSectionNames() []string {
	sections := make([]string, 0, len(p.sections))
	for section := range p.sections {
		sections = append(sections, section)
	}
	return sections
}

// GetSections returns the entire INI data as a map of sections to key-value pairs.
func (p *Parser) GetSections() Ini {
	return p.sections
}

// Get returns the value of a key in a section.
func (p *Parser) Get(section_name, key string) (string, error) {

	// Check if section exists
	_, ok := p.sections[section_name]
	if !ok {
		return "", ErrSectionNotExist
	}

	// Check if key exists
	value, ok := p.sections[section_name][key]
	if !ok {
		return "", ErrKeyNotExist
	}

	return value, nil
}

// Set sets the value of a key in a section.
func (p *Parser) Set(section_name, key, value string) {

	// Create section if it doesn't exist
	if p.sections[section_name] == nil {
		p.sections[section_name] = make(Section)
	}

	// Set key-value pair
	p.sections[section_name][key] = value
}

// String returns the INI data in string format.
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

// SaveToFile saves the INI data to a file.
func (p *Parser) SaveToFile(path string) error {

	// Get INI data as string
	configString := p.String()
	stringBytes := []byte(configString)


	// 0644 is an octal code for access (Owner: read and write, Members of the file's group and other users : read)
	return os.WriteFile(path, stringBytes, 0644)
}