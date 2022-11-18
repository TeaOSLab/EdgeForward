// Copyright 2022 Liuxiangchao iwind.liu@gmail.com. All rights reserved. Official site: https://goedge.cn .

package forwarders

import (
	"github.com/TeaOSLab/EdgeForward/internal/configs"
	"github.com/TeaOSLab/EdgeForward/internal/goman"
	"github.com/iwind/TeaGo/types"
	"io"
	"net"
	"time"
)

type Listener struct {
	rule *configs.Rule

	rawListener net.Listener
}

func NewListener(rule *configs.Rule) *Listener {
	return &Listener{
		rule: rule,
	}
}

func (this *Listener) Start() error {
	rawListener, err := net.Listen("tcp", ":"+types.String(this.rule.SrcPort()))
	if err != nil {
		return err
	}
	this.rawListener = rawListener
	goman.New(func() {
		for {
			conn, err := rawListener.Accept()
			if err != nil {
				break
			}

			go func(conn net.Conn) {
				err := this.handleConn(conn)
				if err != nil {
					_ = conn.Close()
				}
			}(conn)
		}
	})
	return nil
}

func (this *Listener) handleConn(conn net.Conn) error {
	destConn, err := net.DialTimeout("tcp", this.rule.DestHost()+":"+types.String(this.rule.DestPort()), 30*time.Second)
	if err != nil {
		return err
	}

	go func() {
		defer func() {
			_ = destConn.Close()
			_ = conn.Close()
		}()

		_, _ = io.Copy(destConn, conn)
	}()

	go func() {
		defer func() {
			_ = destConn.Close()
			_ = conn.Close()
		}()

		_, _ = io.Copy(conn, destConn)
	}()

	return nil
}

func (this *Listener) Stop() error {
	if this.rawListener == nil {
		return nil
	}
	return this.rawListener.Close()
}
