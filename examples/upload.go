package main

import (
	"fmt"
	"github.com/smartwalle/ngx/curl"
	"net/http"
)

func main() {
	//var req = ngx.NewRequest(ngx.Post, "http://127.0.0.1:9091/test")
	//
	////req.FileForm().AddFilePath("file1", "1.jpg", "./1.jpg")
	////req.FileForm().AddFilePath("file2", "2.png", "./2.png")
	//
	////fmt.Println(req.Do(context.Background()))
	//req.Form().Add("k1", "v1")
	//req.Form().Add("k2", "v2")
	//
	//req.Query().Add("q1", "qv1")
	//req.Query().Add("q2", "qv2")
	//
	//req.FileForm().AddFilePath("file1", "go.mod", "./go.mod")
	//
	//var xxx, _ = req.Request(context.Background())
	//fmt.Println(ngx.CURL(xxx))

	var cmd = curl.New(http.MethodPost, "http://127.0.0.1:9091/test").
		Header("h1", "h'v'1").Header("h2", "h\"v2").
		Form("k1", "kv1").Form("k1", "asss\"ss\"ss\"").FormFile("file", "./upload.go")
	fmt.Println(cmd.Encode())
}
