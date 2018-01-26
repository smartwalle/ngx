package main

import (
	"github.com/smartwalle/ngx"
	"fmt"
)

func main() {
	var r = ngx.NewRequest("GET", "http://www.baidu.com")
	rep := r.Exec()
	fmt.Println(rep.String())
}
