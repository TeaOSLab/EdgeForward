// Copyright 2022 Liuxiangchao iwind.liu@gmail.com. All rights reserved. Official site: https://goedge.cn .

package configs

import (
	"errors"
	"github.com/iwind/TeaGo/Tea"
	_ "github.com/iwind/TeaGo/bootstrap"
	"gopkg.in/yaml.v3"
	"os"
)

type ForwardConfig struct {
	Rules []*Rule `yaml:"rules" json:"rules"`
}

func LoadForwardConfig() (*ForwardConfig, error) {
	data, err := os.ReadFile(Tea.ConfigFile("forward.yaml"))
	if err != nil {
		return nil, err
	}
	var config = &ForwardConfig{}
	err = yaml.Unmarshal(data, config)
	if err != nil {
		return nil, err
	}

	err = config.Init()
	if err != nil {
		return nil, err
	}

	return config, nil
}

func (this *ForwardConfig) Init() error {
	for _, rule := range this.Rules {
		err := rule.Init()
		if err != nil {
			return errors.New("init rule '" + rule.String() + "' failed: " + err.Error())
		}
	}

	return nil
}
