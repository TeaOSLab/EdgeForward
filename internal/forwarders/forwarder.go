// Copyright 2022 Liuxiangchao iwind.liu@gmail.com. All rights reserved. Official site: https://goedge.cn .

package forwarders

import (
	"errors"
	"github.com/TeaOSLab/EdgeForward/internal/configs"
	teaconst "github.com/TeaOSLab/EdgeForward/internal/const"
	"github.com/TeaOSLab/EdgeForward/internal/goman"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/maps"
	"github.com/iwind/TeaGo/types"
	"github.com/iwind/gosock/pkg/gosock"
	"os"
	"runtime"
	"runtime/debug"
	"sort"
)

type Forwarder struct {
	sock *gosock.Sock
}

func NewForwarder() *Forwarder {
	return &Forwarder{
		sock: gosock.NewTmpSock(teaconst.ProcessName),
	}
}

func (this *Forwarder) Listen() error {
	// 本地Sock
	err := this.listenSock()
	if err != nil {
		return errors.New("starting sock failed: " + err.Error())
	}

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

func (this *Forwarder) Test() error {
	logs.Println("loading config ...")
	config, err := configs.LoadForwardConfig()
	if err != nil {
		return err
	}

	logs.Println("testing rules ...")
	for _, rule := range config.Rules {
		logs.Println("testing rule '" + rule.String() + "' ...")
		var listener = NewListener(rule)
		err = listener.Test()
		if err != nil {
			return errors.New("test rule '" + rule.String() + "' failed: " + err.Error())
		}
	}
	return nil
}

// 监听本地sock
func (this *Forwarder) listenSock() error {
	// 检查是否在运行
	if this.sock.IsListening() {
		reply, err := this.sock.Send(&gosock.Command{Code: "pid"})
		if err == nil {
			return errors.New("error: the process is already running, pid: " + types.String(maps.NewMap(reply.Params).GetInt("pid")))
		} else {
			return errors.New("error: the process is already running")
		}
	}

	// 启动监听
	goman.New(func() {
		this.sock.OnCommand(func(cmd *gosock.Command) {
			switch cmd.Code {
			case "pid":
				_ = cmd.Reply(&gosock.Command{
					Code: "pid",
					Params: map[string]interface{}{
						"pid": os.Getpid(),
					},
				})
			case "info":
				exePath, _ := os.Executable()
				_ = cmd.Reply(&gosock.Command{
					Code: "info",
					Params: map[string]interface{}{
						"pid":     os.Getpid(),
						"version": teaconst.Version,
						"path":    exePath,
					},
				})
			case "stop":
				_ = cmd.ReplyOk()
				os.Exit(0)
			case "quit":
				_ = cmd.ReplyOk()
				_ = this.sock.Close()
			case "goman":
				var posMap = map[string]maps.Map{} // file#line => Map
				for _, instance := range goman.List() {
					var pos = instance.File + "#" + types.String(instance.Line)
					m, ok := posMap[pos]
					if ok {
						m["count"] = m["count"].(int) + 1
					} else {
						m = maps.Map{
							"pos":   pos,
							"count": 1,
						}
						posMap[pos] = m
					}
				}

				var result = []maps.Map{}
				for _, m := range posMap {
					result = append(result, m)
				}

				sort.Slice(result, func(i, j int) bool {
					return result[i]["count"].(int) > result[j]["count"].(int)
				})

				_ = cmd.Reply(&gosock.Command{
					Params: map[string]interface{}{
						"total":  runtime.NumGoroutine(),
						"result": result,
					},
				})

			case "gc":
				runtime.GC()
				debug.FreeOSMemory()
				_ = cmd.ReplyOk()
			case "reload":
				// TODO
				_ = cmd.ReplyOk()
			}
		})

		err := this.sock.Listen()
		if err != nil {
			logs.Println(err.Error())
		}
	})

	return nil
}
