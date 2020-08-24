package airgo

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	airtable "github.com/mottaquikarim/go-airtable"
	log "github.com/sirupsen/logrus"
)

type RefreshOpts struct {
	ApiKey string
	BaseId string
	Conf   string
	Dest   string
}

type Config struct {
	Nav []ConfigItem
}

type ConfigItem struct {
	Name     string `json:"name"`
	Key      string `json:"key"`
	Filename string `json:"filenameRoot"`
	Viewname string `json:"viewname"`
}

func RefreshData(cfg RefreshOpts) error {
	var site []byte
	var err error

	if len(cfg.Conf) == 0 {
		return fmt.Errorf("Path to conf cannot be empty")
	} else if len(cfg.ApiKey) == 0 {
		return fmt.Errorf("ApiKey must be provided")
	} else if len(cfg.BaseId) == 0 {
		return fmt.Errorf("BaseId must be provided")
	} else if len(cfg.Dest) == 0 {
		return fmt.Errorf("Path to dest cannot be empty")
	}

	// read file
	site, err = ioutil.ReadFile(cfg.Conf)
	if err != nil {
		return fmt.Errorf("Failed to read config file %v", err)
	}

	// parse json content
	var configData Config
	if err = json.Unmarshal(site, &configData); err != nil {
		return fmt.Errorf("Config is invalid %v", err)
	}

	// if here, we have a config file to process
	// we create an account struct
	// that is passed in to airtable
	// module
	account := airtable.Account{
		ApiKey: cfg.ApiKey,
		BaseId: cfg.BaseId,
	}

	for _, configItem := range configData.Nav {
		var content []airtable.Record
		var message []byte
		var opts = airtable.Options{
			MaxRecords: 100,
			View:       "Grid view",
		}
		if len(configItem.Viewname) > 0 {
			opts.View = configItem.Viewname
		}

		mainSite := airtable.NewTable(configItem.Key, account)
		if content, err = mainSite.List(opts); err != nil {
			log.Printf("Error! %v", err)
		}

		// write to file now
		if message, err = json.Marshal(content); err != nil {
			return fmt.Errorf("Failed to marshal content: %v", err)
		}

		// check if dest exists
		if _, err := os.Stat(cfg.Dest); os.IsNotExist(err) {
			os.Mkdir(cfg.Dest, 0700)
		}
		filename := fmt.Sprintf("%s/%s.json", cfg.Dest, configItem.Filename)
		if err := ioutil.WriteFile(filename, message, 0644); err != nil {
			return fmt.Errorf("Failed to write file: %v", err)
		}
	}

	return nil
}
