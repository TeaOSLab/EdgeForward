// Copyright 2022 Liuxiangchao iwind.liu@gmail.com. All rights reserved. Official site: https://goedge.cn .

package forwarders

import (
	"github.com/TeaOSLab/EdgeForward/internal/configs"
	"github.com/iwind/TeaGo/logs"
)

type Forwarder struct {
}

func NewForwarder() *Forwarder {
	return &Forwarder{}
}

func (this *Forwarder) Listen() error {
	logs.Println("loading config ...")
	config, err := configs.LoadForwardConfig()
	if err != nil {
		return err
	}

	logs.Println("starting rules ...")
	for _, rule := range config.Rules {
		logs.Println("starting '" + rule.String() + "' ...")
		var listener = NewListener(rule)
		go func(rule *configs.Rule) {
			err = listener.Start()
			if err != nil {
				logs.Println("start '" + rule.String() + "' failed: " + err.Error())
			}
		}(rule)
	}

	// hold the process
	select {}
}
