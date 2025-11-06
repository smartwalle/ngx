package main

import (
	"fmt"
	"github.com/smartwalle/ngx/curl"
	"net/http"
	"net/url"
)

func main() {
	var cmd = curl.New(http.MethodPost, "http://127.0.0.1:8080/h1").
		Header("h1", "h'v'1").
		Header("h2", "h\"v2").
		Data(url.Values{"k1": []string{"v\"1"}, "k2": []string{"v'2"}}.Encode())
	fmt.Println(cmd.Encode())
}
