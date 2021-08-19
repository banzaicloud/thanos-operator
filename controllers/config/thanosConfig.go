package config

import (
	"sync"
)

type ControllerConfig struct {
	*sync.Mutex
	Values map[string]interface{}
}

var instance *ControllerConfig
var once sync.Once

func GetControllerConfig() *ControllerConfig {
	once.Do(func() {
		instance = &ControllerConfig{
			Mutex:  &sync.Mutex{},
			Values: map[string]interface{}{},
		}
	})
	return instance
}

func (c *ControllerConfig) AddConfigItem(key string, value interface{}) {
	c.Lock()
	defer c.Unlock()
	if key != "" && value != nil && value != "" {
		c.Values[key] = value
	}
}

func (c *ControllerConfig) GetConfigString(key string, defaultValue string) string {
	if c.HasConfigItem(key) {
		return c.Values[key].(string)
	}
	return defaultValue
}

func (c *ControllerConfig) HasConfigItem(key string) bool {
	c.Lock()
	defer c.Unlock()
	_, ok := c.Values[key]
	return ok
}
