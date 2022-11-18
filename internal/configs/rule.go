// Copyright 2022 Liuxiangchao iwind.liu@gmail.com. All rights reserved. Official site: https://goedge.cn .

package configs

import (
	"errors"
	"github.com/iwind/TeaGo/types"
	"net"
	"regexp"
)

var portReg = regexp.MustCompile(`^\d{1,5}$`)

type Rule struct {
	Src  string `yaml:"src" json:"src"`
	Dest string `yaml:"dest" json:"dest"`

	srcPort  int
	destHost string
	destPort int
}

func (this *Rule) Init() error {
	if !portReg.MatchString(this.Src) {
		return errors.New("invalid src port '" + this.Src + "'")
	}
	this.srcPort = types.Int(this.Src)
	if this.srcPort <= 0 || this.srcPort >= 65535 {
		return errors.New("invalid src port '" + this.Src + "'")
	}

	// support local port to port
	if portReg.MatchString(this.Dest) {
		this.destHost = "127.0.0.1"
		this.destPort = types.Int(this.Dest)
	} else { // support address "host:port"
		destHost, destPort, err := net.SplitHostPort(this.Dest)
		if err != nil {
			return errors.New("invalid dest address '" + this.Dest + "'")
		}
		this.destHost = destHost
		this.destPort = types.Int(destPort)
	}

	return nil
}

func (this *Rule) SrcPort() int {
	return this.srcPort
}

func (this *Rule) DestHost() string {
	return this.destHost
}

func (this *Rule) DestPort() int {
	return this.destPort
}

func (this *Rule) String() string {
	return this.Src + " => " + this.Dest
}
