//+build wireinject

package main

import (
	"github.com/google/wire"
	"gopkg.in/yaml.v2"
	"io/ioutil"

	"github.com/ONSdigital/ras-rm-sample-file-uploader/config"
	"github.com/ONSdigital/ras-rm-sample-file-uploader/file"
)

func Inject() file.FileProcessor{
	wire.Build(NewFileProcessor, ConfigSetup)
	return file.FileProcessor{}
}

func ConfigSetup() config.Config {
	file, err := ioutil.ReadFile("application.yaml")
	if err != nil {
		panic(err)
	}
	config := &config.Config{}
	yaml.Unmarshal(file, &config)
	return *config
}

func NewFileProcessor(config config.Config) file.FileProcessor {
	return file.FileProcessor{Config: config}
}