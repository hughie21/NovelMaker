// Description: Program-configured controller, which is used to manage the configuration of the program.
// Author: Hughie21
// Date: 2024-11-29
// license that can be found in the LICENSE file.
package config

import (
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"sync"

	"github.com/hughie21/NovelMaker/lib/utils"

	"gopkg.in/yaml.v2"
)

type ConfigManager struct {
	config *Config
	path   string
}

var (
	configManager *ConfigManager
	once          *sync.Once
)

// constructor for ConfigManager, which is used to initialize the ConfigManager instance
func NewConfigManager(path string) *ConfigManager {
	if once == nil {
		once = &sync.Once{}
	}
	once.Do(func() {
		configManager = &ConfigManager{
			path:   path,
			config: &Config{},
		}
	})
	return configManager
}

// LoadConfig is used to load the configuration file
func (cm *ConfigManager) LoadConfig() error {
	fp, err := os.ReadFile(filepath.Join(cm.path, "config.yaml"))
	if err != nil {
		utils.ShowMessage("Error", "Failed to load config file: "+err.Error(), "error")
		return err
	}

	var config Config
	err = yaml.Unmarshal(fp, &config)
	if err != nil {
		utils.ShowMessage("Error", "Failed to unmarshal config file: "+err.Error(), "error")
		return err
	}
	cm.config = &config
	return nil
}

// save the conifguration file
func (cm *ConfigManager) SaveConfig() error {
	fp, err := os.Create(filepath.Join(cm.path, "config.yaml"))
	if err != nil {
		utils.ShowMessage("Error", "Failed to save config file: "+err.Error(), "error")
		return err
	}

	data, err := yaml.Marshal(cm.config)
	if err != nil {
		utils.ShowMessage("Error", "Failed to marshal config file: "+err.Error(), "error")
		return err
	}

	_, err = fp.Write(data)
	if err != nil {
		utils.ShowMessage("Error", "Failed to write config file: "+err.Error(), "error")
		return err
	}

	return nil
}

// get the information of the about page
func (cm *ConfigManager) GetInfo() map[string]string {
	info := map[string]string{
		"version":   cm.config.Version,
		"buildTime": cm.config.BuildTime,
		"commit":    cm.config.Commit,
	}
	return info
}

// get the configuration of the program
func (cm *ConfigManager) GetConfig() *Config {
	return cm.config
}

// This method is the one that interacts with the front-end,
// which can't read Go's structs directly, so it's queried by
// passing in the key name
func (cm *ConfigManager) GetConfigByKey(section string, key string) (string, error) {
	v := reflect.ValueOf(cm.config).Elem().FieldByName(section)
	if !v.IsValid() {
		return "", errors.New("invalid sector")
	}

	field := v.FieldByName(key)
	if !field.IsValid() {
		return "", errors.New("invalid key")
	}

	switch field.Kind() {
	case reflect.String:
		return field.String(), nil
	case reflect.Bool:
		return strconv.FormatBool(field.Bool()), nil
	case reflect.Int:
		return strconv.Itoa(int(field.Int())), nil
	default:
		return "", errors.New("unsupported field type")
	}
}

// Set the configuration of the program by the key, value pair
func (cm *ConfigManager) SetConfig(sector string, key string, value string) error {
	v := reflect.ValueOf(cm.config).Elem().FieldByName(sector)
	if !v.IsValid() {
		return errors.New("invalid sector")
	}

	field := v.FieldByName(key)
	if !field.IsValid() {
		return errors.New("invalid key")
	}

	if !field.CanSet() {
		return errors.New("cannot set field")
	}

	switch field.Kind() {
	case reflect.String:
		field.SetString(value)
	case reflect.Bool:
		if value == "true" {
			field.SetBool(true)
		} else if value == "false" {
			field.SetBool(false)
		} else {
			return errors.New("invalid value for bool")
		}
	case reflect.Int:
		intValue, err := strconv.Atoi(value)
		if err != nil {
			return errors.New("invalid value for int")
		}
		field.SetInt(int64(intValue))
	default:
		return errors.New("unsupported field type")
	}

	return nil
}
