package config

import (
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"sync"

	"github.com/hughie21/NovelMaker/lib/logging"
	"github.com/hughie21/NovelMaker/lib/utils"

	"gopkg.in/yaml.v2"
)

var logger = logging.NewLog(logging.FatalLevel, true)

type ConfigManager struct {
	config *Config
	path   string
}

var (
	configManager *ConfigManager
	once          *sync.Once
)

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

func (cm *ConfigManager) SaveConfig() error {
	fp, err := os.Create(filepath.Join(cm.path, "config.yaml"))
	if err != nil {
		logger.Error(err.Error(), logging.RunFuncName())
		return err
	}

	data, err := yaml.Marshal(cm.config)
	if err != nil {
		logger.Error(err.Error(), logging.RunFuncName())
		return err
	}

	_, err = fp.Write(data)
	if err != nil {
		logger.Error(err.Error(), logging.RunFuncName())
		return err
	}

	return nil
}

func (cm *ConfigManager) GetConfig() *Config {
	return cm.config
}

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
