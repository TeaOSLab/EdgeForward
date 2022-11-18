// Copyright 2022 Liuxiangchao iwind.liu@gmail.com. All rights reserved. Official site: https://goedge.cn .

package configs_test

import (
	"github.com/TeaOSLab/EdgeForward/internal/configs"
	"github.com/iwind/TeaGo/logs"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	config, err := configs.LoadForwardConfig()
	if err != nil {
		t.Fatal(err)
	}
	err = config.Init()
	if err != nil {
		t.Fatal(err)
	}
	logs.PrintAsJSON(config, t)
}
