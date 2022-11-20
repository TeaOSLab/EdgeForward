// Copyright 2022 Liuxiangchao iwind.liu@gmail.com. All rights reserved. Official site: https://goedge.cn .

package main

import (
	"github.com/TeaOSLab/EdgeForward/internal/apps"
	teaconst "github.com/TeaOSLab/EdgeForward/internal/const"
	"github.com/TeaOSLab/EdgeForward/internal/forwarders"
	_ "github.com/iwind/TeaGo/bootstrap"
	"github.com/iwind/TeaGo/logs"
)

func main() {
	var app = apps.NewAppCmd().
		Version(teaconst.Version).
		Product(teaconst.ProductName).
		Usage(teaconst.ProcessName + " [-v|start|stop|restart|status|test]")
	app.On("test", func() {
		err := forwarders.NewForwarder().Test()
		if err != nil {
			logs.Println("[ERROR]" + err.Error())
		} else {
			logs.Println("everything is ok!")
		}
	})
	app.Run(func() {
		err := forwarders.NewForwarder().Listen()
		if err != nil {
			logs.Println("start failed: " + err.Error())
		}
	})
}
