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

# Format
### When using the INIParser library, it's important to follow these rules to ensure proper usage.
#### These rules include:

- Comments just at the beginning of a line: Comments in INI files or strings are only valid when they appear at the beginning of a line and are preceded by a semicolon (;).

- Ensuring trimmed keys, values, and section headers: Leading and trailing spaces are trimmed by the library, and keys, and section headers cannot be empty.

- Using the equals sign (=) as the key-value separator: INI files or strings use the equals sign to denote the assignment of a value to a key.


- No global keys: The library assumes that all keys must belong to a section, and global keys are not permitted in INI files or strings.