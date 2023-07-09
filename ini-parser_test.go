package main

import (
	"testing"
	"sort"
	"reflect"
)

const ini_Valid_Syntax ="[Simple Values]\nkey=value\nspaces in keys=allowed\nspaces in values=allowed as well\n[Complex Values]\nspaces around the delimiter=obviously\n"

var ini_Invalid_Syntax =[] string{"# Hi\n[Simple Values\n", ";Hi\n[Simple Values]\nkey-value"}


func TestLoadFromString(t * testing.T){
	c := Config{}
	c.InitializeMap()

	t.Run("Valid INI Syntax", func (t * testing.T){
		err :=c.LoadFromString(ini_Valid_Syntax)
		
		assertErrorsMatching(t, err,nil, ini_Valid_Syntax)
	})

	t.Run("Invalid INI Syntax", func (t * testing.T){

		for _, text := range ini_Invalid_Syntax{

			err :=c.LoadFromString(text)
			assertErrorsMatching( t, err, ErrNotValidSyntax, text)
		}
		
	})
}

func assertErrorsMatching(t testing.TB, got, want error, testing_String string){
	t.Helper()

	if got!=want{
		t.Errorf("got %q want %q with text %q", got ,want, testing_String)
	}
}


func TestLoadFromFile(t * testing.T){
	c := Config{}
	c.InitializeMap()

	t.Run("Valid File Path",func (t * testing.T){
		valid_Path := "config.ini"
		err := c.LoadFromFile(valid_Path)

		if err!=nil{
			t.Errorf("Expected no error but got %q", err.Error())
		}
		
	})

	t.Run("Not Valid File Path",func (t * testing.T){
		valid_Path := "configuration.ini"
		err := c.LoadFromFile(valid_Path)

		if err == nil {
			t.Errorf("Expected error but got no error")
		}
		
	})
}

func TestGetSectionNames(t * testing.T){
	c := Config{}
	c.InitializeMap()

	t.Run("Get Sections names from empty map", func (t * testing.T){
		gotSections := c.GetSectionNames()

		if len(gotSections)!=0{
			t.Errorf("got %q want %q", gotSections, []string{})
		}
	})

	t.Run("Get Sections names", func (t * testing.T){
		err := c.LoadFromString(ini_Valid_Syntax)
		if err!=nil{
			t.Errorf("Error in loading the string, Error message: %q", err.Error())
		}

		gotSections := c.GetSectionNames()
		wantedSections := []string{"Simple Values", "Complex Values"}

		// we don't care about the order
		sort.Strings(gotSections)
		sort.Strings(wantedSections)
		
		if !reflect.DeepEqual(gotSections, wantedSections) {
			t.Errorf("Actual %v does not match expected %v", gotSections, wantedSections)
		}
	})
}


func TestGetSections(t * testing.T){
	c := Config{}
	c.InitializeMap()

	t.Run("Get Sections from empty map", func (t * testing.T){

		gotSections := c.GetSections()

		if len(gotSections)!=0{
			t.Errorf("got %q want Empty Map", gotSections)
		}
	})

	t.Run("Get Sections", func (t * testing.T){

		err := c.LoadFromString(ini_Valid_Syntax)
		if err!=nil{
			t.Errorf("Error in loading the string, Error message: %q", err.Error())
		}

		gotSections := c.GetSections()

		wantedSections := map[string]map[string]string{
			"Simple Values": map[string]string{
				"key": "value",
				"spaces in keys": "allowed",
				"spaces in values": "allowed as well",
			},
			"Complex Values": map[string]string{
				"spaces around the delimiter": "obviously",
			},
		}

		if !reflect.DeepEqual(gotSections, wantedSections) {
			t.Errorf("Actual map %v does not match expected map %v", gotSections, wantedSections)
		}
		
	})
}

func TestGet(t * testing.T){
	c := Config{}
	c.InitializeMap()

	t.Run("Get value from empty map", func (t * testing.T){

		got := c.Get("Simple Values", "key")

		want :=""
		assertStrings(t, got, want)
	})

	t.Run("Get value from empty section", func (t * testing.T){

		err := c.LoadFromString(ini_Valid_Syntax)
		if err!=nil{
			t.Errorf("Error in loading the string, Error message: %q", err.Error())
		}

		got := c.Get("Not Existing Section", "key")
		
		want :=""

		assertStrings(t, got, want)
	})

	t.Run("Get existing value", func (t * testing.T){
		got :=c.Get("Simple Values", "key")
		want :="value"
		assertStrings(t, got, want)
    })
}



func assertStrings(t testing.TB, got, want string){
	t.Helper()

	if got!=want{
		t.Errorf("got %q want %q", got ,want)
	}
}



func TestSet(t * testing.T){
	c:= Config{}
	c.InitializeMap()

	t.Run("Set value to map", func (t * testing.T){
		c.Set("Simple Values", "key", "new value")
		got := c.Get("Simple Values", "key")
		want :="new value"
		assertStrings(t, got, want)
	})
}


func TestSaveToFile(t * testing.T){
	c := Config{}
	c.InitializeMap()

	t.Run("Save to file", func (t * testing.T){
		err := c.LoadFromString(ini_Valid_Syntax)
		if err!=nil{
			t.Errorf("Error in loading the string, Error message: %q", err.Error())
		}

		err = c.SaveToFile("config.ini")
		if err!=nil{
			t.Errorf("Error in saving to file, Error message: %q", err.Error())
		}
	})
}

// func TestToString(t * testing.T){
// 	c := Config{}
// 	c.InitializeMap()

// 	t.Run("Convert to String", func (t * testing.T){
// 		err := c.LoadFromString(ini_Valid_Syntax)
// 		if err!=nil{
// 			t.Errorf("Error in loading the string, Error message: %q", err.Error())
// 		}

// 		got := c.ToString()
		
// 		want := ini_Valid_Syntax

// 		assertStrings(t, got, want)
// 	})
// }
