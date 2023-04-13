package main

import (
	"context"
	"fmt"
	"github.com/smartwalle/ngx"
	"time"
)

func main() {
	var req = ngx.NewRequest(ngx.Post,
		"http://192.168.1.99:9090/upload",
		ngx.WithReceive(func(total uint64, finished uint64) {
			fmt.Println("已接收:", total, finished)
		}),
		ngx.WithSend(func(total uint64, finished uint64) {
			fmt.Println("已发送:", time.Now().Unix(), total, finished)
		}),
	)

	req.AddFile("file1", "", "./1.jpg")
	req.AddFile("file2", "", "./2.jpg")

	req.Exec(context.Background())
}
