package main

import (
	"context"
	"io"
	"log"
	"net/http"

	"github.com/smartwalle/ngx"
)

func main() {
	var req = ngx.NewRequest(http.MethodPost, "http://localhost:8080?qk1=qv1")
	req.JoinPath("h1")
	req.Query.Add("qk2", "qv2")
	req.Form.Add("fk1", "fv1")
	req.Form.Add("fk2", "fv2")
	var resp, err = req.Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	var data, _ = io.ReadAll(resp.Body)
	log.Println(string(data))
}
