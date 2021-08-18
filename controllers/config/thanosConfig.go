package config

import (
	"sync"

	"github.com/banzaicloud/thanos-operator/pkg/sdk/api/v1alpha1"
)

type ControllerConfig struct {
	*sync.Mutex
	Values     map[string]interface{}
	//ThanosImage string
	//ThanosImageTag string
}

var instance *ControllerConfig
var once sync.Once

func GetControllerConfig() *ControllerConfig {
	once.Do(func() {
		instance = &ControllerConfig{
			Mutex:      &sync.Mutex{},
			Values:     map[string]interface{}{},
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

func (c *ControllerConfig) GetConfigItem(key string, defaultValue interface{}) interface{} {
	if c.HasConfigItem(key) {
		return c.Values[key]
	}
	return defaultValue
}

func (c *ControllerConfig) HasConfigItem(key string) bool {
	c.Lock()
	defer c.Unlock()
	_, ok := c.Values[key]
	return ok
}

func (c *ControllerConfig) init(){
	c.Values["ThanosImage"] = v1alpha1.ThanosImageRepository
	c.Values["ThanosImageTag"] = v1alpha1.ThanosImageTag
}