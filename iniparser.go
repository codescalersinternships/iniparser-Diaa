package iniparser

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

var (
	// ErrInvalidFormat is returned when the INI file has an invalid format.
	ErrInvalidFormat = errors.New("invalid format ")

	// ErrInvalidExtension is returned when the file is not in the INI format or does not have a .ini extension.
	ErrInvalidExtension = errors.New("file is not in the ini format or does not have a .ini extension")

	// ErrKeyNotExist is returned when the given key does not exist.
	ErrKeyNotExist = errors.New("key doesn't exist")

	// ErrSectionNotExist is returned when the given section does not exist.
	ErrSectionNotExist = errors.New("section doesn't exist")

	//ErrSectionsNotEmpty is returned if a load operation is attempted while the sections is not empty
	ErrSectionsNotEmpty = errors.New("sections is not empty")
)

// Section is an alias for a map of string key-value pairs representing a section in INI data.
type Section map[string]string

// IniData is an alias for a map of string keys to Section values representing the entire INI data.
type IniData map[string]Section

// Parser represents an INI parser Object contains a map of sections representing the INI data.
type Parser struct {
	sections IniData
}

// NewParser returns a new Parser object.
func NewParser() Parser {
	return Parser{sections: make(IniData)}
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

	p.sections = make(IniData)

	var currentSection string = ""
	scanner := bufio.NewScanner(reader)
	idx := 0
	for scanner.Scan() {
		idx++
		line := scanner.Text()

		line = strings.TrimSpace(line)

		// Skip empty lines and comments
		if len(line) == 0 || line[0] == ';' || line[0] == '#' {
			continue
		}

		// Check for section header
		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {

			// Extract section name
			sectionName := line[1 : len(line)-1]
			sectionName = strings.TrimSpace(sectionName)

			if len(sectionName) == 0 {
				return fmt.Errorf("%w: invalid section at line %d", ErrInvalidFormat, idx)
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
				return fmt.Errorf("%w: invalid key at line %d", ErrInvalidFormat, idx)
			}

			// Add key-value pair to current section
			p.sections[currentSection][key] = value
		} else {
			return fmt.Errorf("%w at line %d", ErrInvalidFormat, idx)
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
func (p *Parser) GetSections() IniData {
	return p.sections
}

// Get returns the value of a key in a section.
func (p *Parser) Get(sectionName, key string) (string, error) {

	// Check if section exists
	_, ok := p.sections[sectionName]
	if !ok {
		return "", ErrSectionNotExist
	}

	// Check if key exists
	value, ok := p.sections[sectionName][key]
	if !ok {
		return "", ErrKeyNotExist
	}

	return value, nil
}

// Set sets the value of a key in a section.
func (p *Parser) Set(sectionName, key, value string) {

	// Create section if it doesn't exist
	if p.sections[sectionName] == nil {
		p.sections[sectionName] = make(Section)
	}

	// Set key-value pair
	p.sections[sectionName][key] = value
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
	fileExt := filepath.Ext(path)
	if fileExt != ".ini" {
		return ErrInvalidExtension
	}

	// Get INI data as string
	configString := p.String()
	stringBytes := []byte(configString)

	// 0644 is an octal code for access (Owner: read and write, Members of the file's group and other users : read)
	return os.WriteFile(path, stringBytes, 0644)
}
