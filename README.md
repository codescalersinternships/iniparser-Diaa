# INI Parser for Go

This package provides a parser for INI files in Go.

# Parsing INI Data
To parse INI data, create a new Parser object using the NewParser() function:
```
parser := iniparser.NewParser()
```

You can then load INI data into the parser from a string, file, or io.Reader using one of the following methods:

```
// Load from a string
err := parser.LoadFromString("[section]\nkey=value\n")

// Load from a file
err := parser.LoadFromFile("/path/to/file.ini")

// Load from an io.Reader
file, _ := os.Open("/path/to/file.ini")
err := parser.LoadFromReader(bufio.NewReader(file))
```

# Retrieving INI Data
You can retrieve the data using the following methods:
```
// Get the names of all sections
sectionNames := parser.GetSectionNames()

// Get the entire INI data as a map of sections to key-value pairs
iniData := parser.GetSections()

// Get the value of a key in a section
value, err := parser.Get("section", "key")

// Get the INI data as a string
iniString := parser.String()
```

# Modifying INI Data
You can modify the INI data using the following methods:

```
// Set the value of a key in a section
parser.Set("section", "key", "value")
```

# Saving INI Data

You can save the INI data to a file using the SaveToFile() method:

```
err := parser.SaveToFile("/path/to/file.ini")
```

# Error Handling
The following error messages can be returned by the parser:

```
iniparser.ErrInvalidFormat
iniparser.ErrInvalidExtension
iniparser.ErrKeyNotExist
iniparser.ErrSectionNotExist
```
