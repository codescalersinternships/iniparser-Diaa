package iniparser

import (
	"errors"
	"reflect"
	"sort"
	"strings"
	"testing"
	"os"
)

// no underscore in the name
const iniValidFormat = `
[Simple Values]
key=value
spaces in keys=allowed
spaces in values=allowed as well
[Complex Values]
spaces around the delimiter=obviously
`

var iniInvalidFormat = []string{
	`# Hi
[Simple Values`,
	`Hi
[Simple Values]
key-value`,
	`[Simple Values][section]`,
}

func TestLoadFromString(t *testing.T) {
	p := Parser{}
	p.sections = NewINI()

	t.Parallel()

	t.Run("Valid INI Syntax", func(t *testing.T) {
		err := p.LoadFromString(iniValidFormat)

		if !errors.Is(err, nil) {
			t.Errorf("got %q want %q with text %q", err, "nil", iniValidFormat)
		}
	})

	t.Run("Invalid INI Syntax", func(t *testing.T) {

		for _, text := range iniInvalidFormat {

			err := p.LoadFromString(text)

			if !errors.Is(err, ErrInvalidFormat) {
				t.Errorf("got %q want %q with text %q", err, ErrInvalidFormat, iniInvalidFormat)
			}
		}

	})
}

func TestLoadFromFile(t *testing.T) {
	p := Parser{}
	p.sections = NewINI()

	t.Parallel()

	t.Run("Valid File Path", func(t *testing.T) {

		file, err := os.CreateTemp("", "config.ini")
		if err != nil {
			t.Errorf("Error creating temporary file: %q", err)
		}
		defer os.Remove(file.Name())

		validFile := "config.ini"

		err = p.LoadFromFile(validFile)

		if err != nil {
			t.Errorf("expected no error but got %q", err.Error())
		}

	})

	t.Run("Not Valid File Path", func(t *testing.T) {
		invalidFile := "configuration.ini"
		err := p.LoadFromFile(invalidFile)

		if err == nil {
			t.Errorf("expected error but got no error")
		}
	})

	t.Run("Invalid Extension", func(t *testing.T) {
		invalidExtension := "package.json"

		err := p.LoadFromFile(invalidExtension)

		if !errors.Is(err, ErrInvalidExtension) {
			t.Errorf("want %q but got %q", ErrInvalidExtension, err)
		}
	})
}

func TestGetSectionNames(t *testing.T) {
	p := Parser{}
	p.sections = NewINI()

	t.Parallel()

	t.Run("Get Sections names from empty map", func(t *testing.T) {
		gotSections := p.GetSectionNames()

		if len(gotSections) != 0 {
			t.Errorf("got %q want %q", gotSections, []string{})
		}
	})

	t.Run("Get Sections names", func(t *testing.T) {
		err := p.LoadFromString(iniValidFormat)
		if err != nil {
			t.Errorf("error in loading the string, Error message: %q", err.Error())
		}

		gotSections := p.GetSectionNames()
		wantedSections := []string{"Simple Values", "Complex Values"}

		// we don't care about the order
		sort.Strings(gotSections)
		sort.Strings(wantedSections)

		if !reflect.DeepEqual(gotSections, wantedSections) {
			t.Errorf("actual %v does not match expected %v", gotSections, wantedSections)
		}
	})
}

func TestGetSections(t *testing.T) {
	p := Parser{}
	p.sections = NewINI()

	t.Parallel()

	t.Run("Get Sections from empty map", func(t *testing.T) {

		gotSections := p.GetSections()

		if len(gotSections) != 0 {
			t.Errorf("got %q want Empty Map", gotSections)
		}
	})

	t.Run("Get Sections from non Empty Map", func(t *testing.T) {

		err := p.LoadFromString(iniValidFormat)
		if err != nil {
			t.Errorf("error in loading the string, Error message: %q", err.Error())
		}

		got := p.GetSections()

		wanted := Ini{
			"Simple Values": Section{
				"key":              "value",
				"spaces in keys":   "allowed",
				"spaces in values": "allowed as well",
			},
			"Complex Values": Section{
				"spaces around the delimiter": "obviously",
			},
		}

		for section, gotSectionData := range got {
			wantedSectionData, ok := wanted[section]
			if !ok || !assertSectionData(gotSectionData, wantedSectionData) {
				t.Errorf("actual map %v does not match expected map %v", got, wanted)
			}

		}
	})
}

func assertSectionData(gotSection, wantedSection Section) bool {
	for key, value := range gotSection {
		wantedValue, ok := wantedSection[key]
		if !ok || value != wantedValue {
			return false
		}
	}
	return true
}

func TestGet(t *testing.T) {
	p := Parser{}
	p.sections = NewINI()

	t.Parallel()

	t.Run("Get value from empty map", func(t *testing.T) {

		_, err := p.Get("Simple Values", "key")

		if err != ErrKeyNotExist {
			t.Errorf("want %q but got %q", ErrKeyNotExist, err)
		}
	})

	t.Run("Get value from empty section", func(t *testing.T) {

		err := p.LoadFromString(iniValidFormat)
		if err != nil {
			t.Errorf("error in loading the string, Error message: %q", err.Error())
		}

		_, err = p.Get("Not Existing Section", "key")

		if !errors.Is(ErrKeyNotExist, err) {
			t.Errorf("want %q but got %q", ErrKeyNotExist, err)
		}
	})

	t.Run("Get existing value", func(t *testing.T) {
		got, _ := p.Get("Simple Values", "key")
		want := "value"

		assertStrings(t, got, want)
	})
}

func TestSet(t *testing.T) {
	p := Parser{}
	p.sections = NewINI()

	t.Run("Set value to map", func(t *testing.T) {
		p.Set("Simple Values", "key", "new value")
		got, _ := p.Get("Simple Values", "key")
		want := "new value"
		assertStrings(t, got, want)
	})
}

func assertStrings(t testing.TB, got, want string) {
	if got != want {
		t.Errorf("got %q want %q", got, want)
	}
}

func TestSaveToFile(t *testing.T) {
	p := Parser{}
	p.sections = NewINI()

	t.Parallel()

	t.Run("Save to file with correct path", func(t *testing.T) {
		err := p.LoadFromString(iniValidFormat)
		if err != nil {
			t.Errorf("error in loading the string, Error message: %q", err.Error())
		}

		err = p.SaveToFile("config.ini")
		if err != nil {
			t.Errorf("wanted nil got %q", err.Error())
		}
	})

	t.Run("Save to file with wrong path", func(t *testing.T) {
		err := p.LoadFromString(iniValidFormat)
		if err != nil {
			t.Errorf("Error in loading the string, Error message: %q", err.Error())
		}

		err = p.SaveToFile("/folder/config.ini")
		if err == nil {
			t.Errorf("wanted error on saving but got nil")
		}
	})
}

func TestString(t *testing.T) {
	p := Parser{}
	p.sections = NewINI()

	t.Parallel()
	
	t.Run("Testing String Function", func(t *testing.T) {
		err := p.LoadFromString(iniValidFormat)

		if err != nil {
			t.Errorf("Error in loading the string, Error message: %q", err.Error())
		}
		out := p.String()

		inputNoSpaces := strings.ReplaceAll(iniValidFormat, " ", "")
		outNoSpaces := strings.ReplaceAll(out, " ", "")

		for section, sectionData := range p.GetSections() {
			sectionNoSpaces := strings.ReplaceAll(section, " ", "")

			if !assertContainsSubString(inputNoSpaces, outNoSpaces, "["+sectionNoSpaces+"]") {
				t.Errorf("Expected section [%s] not found in output: %s", sectionNoSpaces, out)
			}

			for key, value := range sectionData {
				keyNoSpaces := strings.ReplaceAll(key, " ", "")
				valueNoSpaces := strings.ReplaceAll(value, " ", "")

				if !assertContainsSubString(inputNoSpaces, outNoSpaces, keyNoSpaces+"="+valueNoSpaces) {
					t.Errorf("Expected section [%s] not found in output: %s", sectionNoSpaces, out)
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