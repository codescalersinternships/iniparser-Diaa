package main

import (
	"strings"
	reg "regexp"
	"os"
)


func main(){
	c := Config{
		iniConfigs: make(map[string]map[string]string),
	}

	c.LoadFromString(" # Hi\n[Simple Values]\nkey=value\nspaces in keys=allowed\nspaces in values=allowed as well\n[Complex Values]\nspaces around the delimiter = obviously\n")

	// config.GetSectionNames()
	// config.LoadFromFile("/tmp/dat")

	c.Set("section","key2","value2")
	c.SaveToFile()


}




type Methods interface {
	LoadFromString(configs string)
	LoadFromFile(path string)
}

type Config struct{
	iniConfigs map[string]map[string]string
}

func (config* Config) LoadFromString(configs string){
	lines := strings.Split(configs,"\n")
	var lastSection string=""

	for _,line := range lines {
		
		line = strings.Trim(line," ")

		// comment
		if  len(line)==0 || string(line[0])=="#" {
			continue
		}

		// it will panic the function if any error happens
		regex := reg.MustCompile(`\[[^\[\]]*\]`)

		section := regex.FindStringSubmatch(line)

		// new section
		if(len(section)>0){

			// getting section name
			lastSection = strings.ReplaceAll(section[0],"[","")
			lastSection=strings.ReplaceAll(lastSection,"]","")
			
			config.iniConfigs[lastSection]=make(map[string]string)

		} else if strings.Contains(line,"=")&& lastSection!=""{
			// key and value line
			keyAndValue :=strings.Split(line,"=")
			key,value :=keyAndValue[0],keyAndValue[1]

			config.iniConfigs[lastSection][key]=value
		}
	}
}



func check(err error){
	if err!=nil{
		panic(err)
	}
}

func (config * Config) LoadFromFile(path string){
	data, err := os.ReadFile(path)
	check(err)

	config.LoadFromString(string(data))
}


func (config * Config) GetSectionNames ()[] string{

	sections :=make([]string,0, len(config.iniConfigs))
	for section,_:=range config.iniConfigs {
		sections = append(sections,section)
	}
	return sections
}

func (config * Config) GetSections () map[string]map[string]string {
	return config.iniConfigs
}

func (config *Config) Get(section_name, key string) string {
	return config.iniConfigs[section_name][key]
}

func (config *Config) Set(section_name, key, value string) {
	
	// checking if the section doesn't exist (new section)
	if config.iniConfigs[section_name]==nil{
		config.iniConfigs[section_name] = make(map[string]string)
	}
	config.iniConfigs[section_name][key]=value
}


func (config * Config) ToString() string {

	configText :=""
	for section,configs:=range config.iniConfigs {
		configText+= "["+section+"]\n"
		for key,value := range configs {
			configText+=key +"=" + value +"\n"
		}
	}

	return configText
}


func (config * Config) SaveToFile() {

	configString := config.ToString()
	stringBytes := []byte(configString)

	// 0644 is an octal code for access (admin: read and write, other users :read)
	err := os.WriteFile("config.ini",stringBytes, 0644)
	check(err)
}

