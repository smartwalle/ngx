package main

import (
	"bytes"
	"context"
	"github.com/smartwalle/ngx"
)

func main() {
	var get = ngx.NewRequest(ngx.Get, "http://192.168.1.99:9090/get?q1=GET URL中的参数1&q2=GET URL中的参数2")
	get.AddParam("k1", "GET 请求参数1")
	get.AddParam("k2", "GET 请求参数2")
	get.AddQuery("q3", "GET URL中的参数3")
	get.Exec(context.Background())

	var getBody = ngx.NewRequest(ngx.Get, "http://192.168.1.99:9090/get_body?q1=GET URL中的参数1&q2=GET URL中的参数2")
	getBody.AddParam("k1", "GET 请求参数1")
	getBody.AddParam("k2", "GET 请求参数2")
	getBody.SetBody(bytes.NewReader([]byte("hello, 这段文字来源于 GET 请求中的 Body。")))
	getBody.SetContentType(ngx.ContentTypeText)
	getBody.Exec(context.Background())

	var post = ngx.NewRequest(ngx.Post, "http://192.168.1.99:9090/post?q1=POST URL中的参数1&q2=POST URL中的参数2")
	post.AddParam("k1", "POST 请求参数1")
	post.AddParam("k2", "POST 请求参数2")
	post.AddQuery("q3", "POST URL中的参数3")
	post.Exec(context.Background())

	var postBody = ngx.NewRequest(ngx.Post, "http://192.168.1.99:9090/post_body?q1=POST URL中的参数1&q2=POST URL中的参数2")
	postBody.AddParam("k1", "POST 请求参数1，由于后面设置了 ContentType 为 Text，服务端无法接收到本参数")
	postBody.AddParam("k2", "POST 请求参数2，由于后面设置了 ContentType 为 Text，服务端无法接收到本参数")
	postBody.SetBody(bytes.NewReader([]byte("hello, 这段文字来源于 POST 请求中的 Body。")))
	postBody.SetContentType(ngx.ContentTypeText)
	postBody.Exec(context.Background())

	var put = ngx.NewRequest(ngx.Put, "http://192.168.1.99:9090/put?q1=PUT URL中的参数1&q2=PUT URL中的参数2")
	put.AddParam("k1", "PUT 请求参数1")
	put.AddParam("k2", "PUT 请求参数2")
	put.Exec(context.Background())

	var putBody = ngx.NewRequest(ngx.Put, "http://192.168.1.99:9090/put_body?q1=PUT URL中的参数1&q2=PUT URL中的参数2")
	putBody.AddParam("k1", "PUT 请求参数1，由于后面设置了 ContentType 为 Text，服务端无法接收到本参数")
	putBody.AddParam("k2", "PUT 请求参数2，由于后面设置了 ContentType 为 Text，服务端无法接收到本参数")
	putBody.SetBody(bytes.NewReader([]byte("hello, 这段文字来源于 PUT 请求中的 Body。")))
	putBody.SetContentType(ngx.ContentTypeText)
	putBody.Exec(context.Background())

	var delete = ngx.NewRequest(ngx.Delete, "http://192.168.1.99:9090/delete?q1=DELETE URL中的参数1&q2=DELETE URL中的参数2")
	delete.AddParam("k1", "DELETE 请求参数1")
	delete.AddParam("k2", "DELETE 请求参数2")
	delete.Exec(context.Background())

	var deleteBody = ngx.NewRequest(ngx.Delete, "http://192.168.1.99:9090/delete_body?q1=DELETE URL中的参数1&q2=DELETE URL中的参数2")
	deleteBody.AddParam("k1", "DELETE 请求参数1")
	deleteBody.AddParam("k2", "DELETE 请求参数2")
	//deleteBody.SetBody(bytes.NewReader([]byte("hello, 这段文字来源于 DELETE 请求中的 Body。")))
	//deleteBody.SetContentType(ngx.ContentTypeText)
	deleteBody.Exec(context.Background())

	//var upload = ngx.NewRequest(ngx.Post, "http://192.168.1.99:9090/upload?q1=上传文件URL中的参数1&q2=上传文件URL中的参数2")
	//upload.AddParam("k1", "上传文件请求参数1")
	//upload.AddParam("k2", "上传文件请求参数2")
	//upload.AddFilePath("file1", "", "1.jpg")
	//upload.AddFilePath("file2", "", "2.jpg")
	//upload.Exec(context.Background())
}
