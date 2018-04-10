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
		config.loadYml(fileName)
	} else if strings.HasSuffix(fileName, ".xml") {
		config.loadXml(fileName)
	} else {
		config.loadJson(fileName)
	}
}

func (config *Config) loadYml(fileName string) {
	buff, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatalln(err)
	}
	err = yaml.Unmarshal(buff, config)
	if err != nil {
		log.Fatalln(err)
	}

}
func (config *Config) loadXml(fileName string) {
	buff, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatalln(err)
	}
	xml.Unmarshal(buff, config)
	if err != nil {
		log.Fatalln(err)
	}
}
func (config *Config) loadJson(fileName string) {
	buff, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatalln(err)
	}
	json.Unmarshal(buff, config)
	if err != nil {
		log.Fatalln(err)
	}
}

type ConfigMap map[interface{}]interface{}
func (config *ConfigMap) Load(fileName string) {
	if strings.HasSuffix(fileName, ".yaml") || strings.HasSuffix(fileName, ".yml") {
		config.loadYml(fileName)
	} else if strings.HasSuffix(fileName, ".xml") {
		config.loadXml(fileName)
	} else {
		config.loadJson(fileName)
	}
}

func (config *ConfigMap) loadYml(fileName string) {
	buff, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatalln(err)
	}
	err = yaml.Unmarshal(buff, config)
	if err != nil {
		log.Fatalln(err)
	}

}
func (config *ConfigMap) loadXml(fileName string) {
	buff, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatalln(err)
	}
	xml.Unmarshal(buff, config)
	if err != nil {
		log.Fatalln(err)
	}
}
func (config *ConfigMap) loadJson(fileName string) {
	buff, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatalln(err)
	}
	json.Unmarshal(buff, config)
	if err != nil {
		log.Fatalln(err)
	}
}

