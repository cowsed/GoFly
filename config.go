package main

import (
	"encoding/json"
	"io"
	"log"
	"os"
)

type Config struct {
	CameraFOV       float32
	ModelPath       string
	EnvironmentPath string
}

var DefaultConfig Config = Config{
	CameraFOV:       90,
	ModelPath:       "Assets/Models/Simple/final NP v17.obj",
	EnvironmentPath: "",
}

func LoadConfig() Config {
	fpath := "config.json"
	f, err := os.Open(fpath)
	if err != nil {
		//No File found, save new one
		if err.Error() == "open config.json: no such file or directory" {
			log.Println("Creating new config.json")
			df, err := os.Create("config.json")
			check(err)
			bytes, err := json.Marshal(DefaultConfig)
			check(err)
			df.Write(bytes)
			df.Close()
			return DefaultConfig
		} else {
			panic(err)
		}
	}
	//Load File
	newConf := Config{}
	bytes, err := io.ReadAll(f)
	check(err)
	err = json.Unmarshal(bytes, &newConf)
	check(err)
	log.Println("Loaded config.json")
	return newConf
}
func check(err error) {
	if err != nil {
		panic(err)
	}
}
