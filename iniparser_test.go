package iniparser

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

const iniValidFormat = `
[Simple Values]
key=value
spaces in keys=allowed
spaces in values=allowed as well
[Complex Values]
spaces around the delimiter=obviously
key=
[new[section]]
new key = new value
`

var iniInvalidFormats = []string{
	`# Hi
[Simple Values`,
	`Hi
[Simple Values]
key-value`,
	`[Simple Values][section`,
	`[Simple]
	=value`,
	`[Simple]
	=`,
	`[]
	key=value`,
}

func TestLoadFromString(t *testing.T) {
	p := NewParser()

	t.Parallel()

	t.Run("Valid INI Syntax", func(t *testing.T) {
		err := p.LoadFromString(iniValidFormat)

		assert.Nil(t, err, "got %q want nil", err)
		p.EmptySections()
	})

	t.Run("Invalid INI Syntax", func(t *testing.T) {

		for _, text := range iniInvalidFormats {

			err := p.LoadFromString(text)

			assert.NotEqual(t, ErrInvalidFormat, err, "got %q want %q", err, ErrInvalidFormat)
			p.EmptySections()
		}

	})

	t.Run("Loading data while sections not empty", func(t *testing.T) {

		err := p.LoadFromString(iniValidFormat)

		assert.Nil(t, err, "got %q want %s", err, "nil")

		// trying to load data again
		err = p.LoadFromString(iniValidFormat)

		assert.Equal(t, ErrSectionsNotEmpty, err, "got %q want %q", err, ErrInvalidFormat)

	})
}

func TestLoadFromFile(t *testing.T) {
	p := NewParser()

	t.Parallel()

	t.Run("Valid File Path", func(t *testing.T) {

		dir := t.TempDir()
		filePath := filepath.Join(dir, "config.ini")

		err := os.WriteFile(filePath, make([]byte, 0), 0644)

		assert.Nil(t, err, "error creating temp file: %q", err)

		err = p.LoadFromFile(filePath)

		assert.Nil(t, err, "expected no error but got %q", err)

	})

	t.Run("Not Valid File Path", func(t *testing.T) {
		invalidFile := "configuration.ini"
		err := p.LoadFromFile(invalidFile)

		assert.NotNil(t, err, "expected error but got no error")
	})

	t.Run("Invalid Extension", func(t *testing.T) {
		invalidExtension := "package.json"

		err := p.LoadFromFile(invalidExtension)

		assert.Equal(t, ErrInvalidExtension, err, "want %q but got %q", ErrInvalidExtension, err)
	})
}

func TestGetSectionNames(t *testing.T) {
	p := NewParser()

	t.Parallel()

	t.Run("Get Sections names from empty map", func(t *testing.T) {
		gotSections := p.GetSectionNames()

		assert.Equal(t, 0, len(gotSections), "got %q want %q", gotSections, []string{})

		p.EmptySections()
	})

	t.Run("Get Sections names", func(t *testing.T) {
		err := p.LoadFromString(iniValidFormat)
		assert.Nil(t, err, "error in loading the string, Error message: %q", err)

		gotSections := p.GetSectionNames()
		wantedSections := []string{"Simple Values", "Complex Values", "new[section]"}

		// we don't care about the order
		sort.Strings(gotSections)
		sort.Strings(wantedSections)
		if !reflect.DeepEqual(gotSections, wantedSections) {
			t.Errorf("actual %v does not match expected %v", gotSections, wantedSections)
		}
	})
}

func TestGetSections(t *testing.T) {
	p := NewParser()

	t.Parallel()

	t.Run("Get Sections from empty map", func(t *testing.T) {

		gotSections := p.GetSections()

		assert.Equal(t, 0, len(gotSections), "got %q want Empty Map", gotSections)
	})

	t.Run("Get Sections from non Empty Map", func(t *testing.T) {

		err := p.LoadFromString(iniValidFormat)
		assert.Nil(t, err, "error in loading the string, error message: %q", err)

		got := p.GetSections()

		wanted := IniData{
			"Simple Values": Section{
				"key":              "value",
				"spaces in keys":   "allowed",
				"spaces in values": "allowed as well",
			},
			"Complex Values": Section{
				"spaces around the delimiter": "obviously",
				"key":                         "",
			}, "new[section]": Section{
				"new key": "new value",
			},
		}

		if !reflect.DeepEqual(got, wanted) {
			t.Errorf("actual map %v does not match expected map %v", got, wanted)

		}
	})
}

func TestGet(t *testing.T) {
	p := NewParser()

	t.Parallel()

	t.Run("Get value from not exist section", func(t *testing.T) {

		_, err := p.Get("Not Exist", "key")

		assert.Equal(t, ErrSectionNotExist, err, "want %q but got %q", ErrKeyNotExist, err)
	})

	t.Run("Get value from Key not exist", func(t *testing.T) {

		err := p.LoadFromString(iniValidFormat)
		assert.Nil(t, err, "error in loading the string, error message: %q", err)

		_, err = p.Get("Simple Values", "not exist")

		assert.Equal(t, ErrKeyNotExist, err, "want %q but got %q", ErrKeyNotExist, err)
	})

	t.Run("Get existing value", func(t *testing.T) {
		got, _ := p.Get("Simple Values", "key")
		want := "value"

		assert.Equal(t, want, got, "got %q want %q", got, want)
	})
}

func TestSet(t *testing.T) {
	p := NewParser()

	t.Parallel()
	t.Run("Set value to map", func(t *testing.T) {
		p.Set("Simple Values", "key", "new value")
		got, _ := p.Get("Simple Values", "key")
		want := "new value"
		assert.Equal(t, want, got, "got %q want %q", got, want)
	})
}

func TestSaveToFile(t *testing.T) {
	p := NewParser()

	t.Parallel()

	t.Run("Save to file with wrong extension", func(t *testing.T) {
		err := p.LoadFromString(iniValidFormat)
		assert.Nil(t, err, "error in loading the string, error message: %q", err)

		tempDir := t.TempDir()
		filePath := filepath.Join(tempDir, "config.json")

		err = p.SaveToFile(filePath)
		assert.NotNil(t, err, "wanted %q got nil", ErrInvalidExtension)
	})
	t.Run("Save to file with correct path and ext", func(t *testing.T) {

		tempDir := t.TempDir()

		filePath := filepath.Join(tempDir, "config.ini")

		err := p.SaveToFile(filePath)
		assert.Nil(t, err, "wanted nil got %q", err)
	})

	t.Run("Save to file with wrong path", func(t *testing.T) {

		err := p.SaveToFile("/folder/config.ini")
		assert.NotNil(t, err, "wanted error on saving but got nil")
	})
}

func TestString(t *testing.T) {
	p := NewParser()

	t.Parallel()

	t.Run("Testing String Function", func(t *testing.T) {
		err := p.LoadFromString(iniValidFormat)

		assert.Nil(t, err, "Error in loading the string, error message: %q", err)
		out := p.String()

		inputNoSpaces := strings.ReplaceAll(iniValidFormat, " ", "")
		outNoSpaces := strings.ReplaceAll(out, " ", "")

		for section, sectionData := range p.GetSections() {
			sectionNoSpaces := strings.ReplaceAll(section, " ", "")

			if !assertContainsSubString(inputNoSpaces, outNoSpaces, fmt.Sprintf("[%s]", sectionNoSpaces)) {
				t.Errorf("expected section [%s] not found in output: %s", sectionNoSpaces, out)
			}

			for key, value := range sectionData {
				keyNoSpaces := strings.ReplaceAll(key, " ", "")
				valueNoSpaces := strings.ReplaceAll(value, " ", "")

				if !assertContainsSubString(inputNoSpaces, outNoSpaces, fmt.Sprintf("%s=%s", keyNoSpaces, valueNoSpaces)) {
					t.Errorf("expected section [%s] not found in output: %s", sectionNoSpaces, out)
				}

			}

		}

	})
}

func assertContainsSubString(input, out, target string) bool {

	if !strings.Contains(input, target) || !strings.Contains(out, target) {
		return false
	}
	return true

}
