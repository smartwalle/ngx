package main

import (
	"context"
	"fmt"
	"github.com/smartwalle/ngx"
	"time"
)

func main() {
	var req = ngx.NewRequest(ngx.Post,
		"http://127.0.0.1:9090/upload",
		ngx.WithReceive(func(total, chunk, finished uint64) {
			fmt.Println("已接收:", total, finished)
		}),
		ngx.WithSend(func(total, chunk, finished uint64) {
			fmt.Println("已发送:", time.Now().Unix(), total, finished)
		}),
	)

	req.AddFilePath("file1", "1.jpg", "./1.jpg")
	req.AddFilePath("file2", "2.png", "./2.png")

	fmt.Println(req.Exec(context.Background()).Error())
}
