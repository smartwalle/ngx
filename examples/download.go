package main

import (
	"context"
	"fmt"
	"github.com/smartwalle/ngx"
)

func main() {
	var req = ngx.NewRequest(ngx.Get,
		"https://t7.baidu.com/it/u=1956604245,3662848045&fm=193&f=GIF",
		ngx.WithReceived(func(total uint64, received uint64) {
			fmt.Println("Receive:", total, received)
		}),
	)

	req.Download(context.Background(), "./xx.jpg")
}
