package config

import (
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

type HomeserverConfig struct {
	Name string `yaml:"name"`
	DownloadRequiresAuth bool `yaml:"downloadRequiresAuth"`
	ClientServerApi string `yaml:"csApi"`
}

type MediaRepoConfig struct {
	Homeservers []HomeserverConfig `yaml:"homeservers,flow"`

	Database struct {
		Postgres string `yaml:"postgres"`
	} `yaml:"database"`

	Uploads struct {
		StoragePaths []string `yaml:"storagePaths,flow"`
		MaxSizeBytes int64 `yaml:"maxBytes"`
	} `yaml:"uploads"`

	Downloads struct {
		MaxSizeBytes int64 `yaml:"maxBytes"`
	} `yaml:"downloads"`

	Thumbnails struct {
		MaxSourceBytes int64 `yaml:"maxSourceBytes"`
		Sizes []struct {
			Width int `yaml:"width"`
			Height int `yaml:"height"`
			Method string `yaml:"method"`
		} `yaml:"sizes,flow"`
	} `yaml:"thumbnails"`
}

func ReadConfig() (MediaRepoConfig, error) {
	c := &MediaRepoConfig{}

	f, err := os.Open("media-repo.yaml")
	if err != nil {
		return *c, err
	}

	defer f.Close()

	buffer, err := ioutil.ReadAll(f)
	err = yaml.Unmarshal(buffer, &c)
	if err != nil {
		return *c, err
	}

	return *c, nil
}