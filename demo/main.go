package main

import (
	"github.com/smartwalle/ngx"
	"fmt"
	"context"
	"time"
)

func main() {
	var ctx, cancel = context.WithCancel(context.Background())
	var r = ngx.NewRequest("GET", "http://www.google.com")
	go func() {
		fmt.Println("2 秒后取消")
		time.Sleep(time.Second * 2)
		cancel()
	}()
	rep := r.ExecWithContext(ctx)
	fmt.Println("结果：")
	fmt.Println(rep.String())
}
