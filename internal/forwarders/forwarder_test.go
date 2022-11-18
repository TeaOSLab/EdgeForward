// Copyright 2022 Liuxiangchao iwind.liu@gmail.com. All rights reserved. Official site: https://goedge.cn .

package forwarders_test

import (
	"github.com/TeaOSLab/EdgeForward/internal/forwarders"
	"testing"
)

func TestNewForwarder(t *testing.T) {
	var forwarder = forwarders.NewForwarder()
	err := forwarder.Listen()
	if err != nil {
		t.Fatal(err)
	}
}
