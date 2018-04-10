package context

import (
	"io/ioutil"
	"log"
	"gopkg.in/yaml.v2"
	"encoding/json"
	"encoding/xml"
	"strings"
)

type Config struct {
	Jewel struct {
		Log            string `json:"log"`
		Max_Open_conns int    `json:"max-Open-conns"`
		Max_idle_conns int    `json:"max-idle-conns"`
		Mysql          string `json:"mysql"`
		Name           string `json:"name"`
		Port           int    `json:"port"`
		Postgres       string `json:"postgres"`
		Profiles struct {
			Active string `json:"active"`
		} `json:"profiles"`
		Redis struct {
			Db       int    `json:"db"`
			Host     string `json:"host"`
			Password string `json:"password"`
		} `json:"redis"`
		Sqlite3 string `json:"sqlite3"`
	} `json:"jewel"`
}

func (config *Config) Load(fileName string) {
	if strings.HasSuffix(fileName, ".yaml") || strings.HasSuffix(fileName, ".yml") {
		loadYml(fileName, config)
	} else if strings.HasSuffix(fileName, ".xml") {
		loadXml(fileName, config)
	} else {
		loadJson(fileName, config)
	}
}

func loadYml(fileName string, config *Config) {
	buff, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatalln(err)
	}
	err = yaml.Unmarshal(buff, config)
	if err != nil {
		log.Fatalln(err)
	}

}
func loadXml(fileName string, config *Config) {
	buff, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatalln(err)
	}
	xml.Unmarshal(buff, config)
	if err != nil {
		log.Fatalln(err)
	}
}
func loadJson(fileName string, config *Config) {
	buff, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatalln(err)
	}
	json.Unmarshal(buff, config)
	if err != nil {
		log.Fatalln(err)
	}
}
