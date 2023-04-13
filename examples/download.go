package main

import (
	"context"
	"fmt"
	"github.com/smartwalle/ngx"
)

func main() {
	var req = ngx.NewRequest(ngx.Get,
		"https://t7.baidu.com/it/u=1956604245,3662848045&fm=193&f=GIF",
		ngx.WithReceive(func(total uint64, finished uint64) {
			fmt.Println("Receive:", total, finished)
		}),
	)

	req.Download(context.Background(), "./1.jpg")
}
