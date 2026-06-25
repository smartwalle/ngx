package curl

import (
	"bytes"
	"strconv"
	"strings"
)

// Option 表示一个 curl 命令参数或选项。
type Option struct {
	key    string
	values []string
}

func (opt *Option) write(buffer *bytes.Buffer) {
	if opt.key != "" {
		buffer.WriteString(opt.key)
		buffer.WriteByte(' ')
	}
	if len(opt.values) > 0 {
		escape(buffer, opt.values...)
		buffer.WriteByte(' ')
	}
}

// Command 用于构建经过 shell 转义的 curl 命令。
type Command struct {
	options []*Option
}

// New 基于指定的 HTTP 方法和 URL 创建 curl 命令。
func New(method, url string) *Command {
	var cmd = &Command{
		options: make([]*Option, 0, 6),
	}
	cmd.option("curl")
	cmd.option("--request", method)
	cmd.option("", url)
	return cmd
}

// Header 追加 --header 选项。
func (c *Command) Header(key, value string) *Command {
	return c.option("--header", key+":"+value)
}

// Data 追加 --data 选项。
func (c *Command) Data(data string) *Command {
	return c.option("--data", data)
}

// DataRaw 追加 --data-raw 选项。
func (c *Command) DataRaw(data string) *Command {
	return c.option("--data-raw", data)
}

// DataBinary 追加 --data-binary 选项。
func (c *Command) DataBinary(data string) *Command {
	return c.option("--data-binary", data)
}

// Form 追加带键值字段的 --form 选项。
func (c *Command) Form(key, value string) *Command {
	return c.option("--form", key+"="+value)
}

// File 追加文件上传形式的 --form 选项；filename 不为空时作为上传文件名。
func (c *Command) File(key, filepath, filename string) *Command {
	var value = key + "=@" + filepath
	if filename != "" {
		value += ";filename=" + filename
	}
	return c.option("--form", value)
}

// User 追加 --user 选项。
func (c *Command) User(username, password string) *Command {
	return c.option("--user", username+":"+password)
}

// UserAgent 追加 --user-agent 选项。
func (c *Command) UserAgent(value string) *Command {
	return c.option("--user-agent", value)
}

// Cookie 根据 cookie 名和值追加 --cookie 选项。
func (c *Command) Cookie(key, value string) *Command {
	return c.option("--cookie", key+"="+value)
}

// CookieRaw 使用指定的 cookie 字符串追加 --cookie 选项。
func (c *Command) CookieRaw(value string) *Command {
	return c.option("--cookie", value)
}

// Referer 追加 --referer 选项。
func (c *Command) Referer(value string) *Command {
	return c.option("--referer", value)
}

// Output 追加 --output 选项。
func (c *Command) Output(filepath string) *Command {
	return c.option("--output", filepath)
}

// ConnectTimeout 追加以秒为单位的 --connect-timeout 选项。
func (c *Command) ConnectTimeout(seconds int) *Command {
	return c.option("--connect-timeout", strconv.Itoa(seconds))
}

// MaxTime 追加以秒为单位的 --max-time 选项。
func (c *Command) MaxTime(seconds int) *Command {
	return c.option("--max-time", strconv.Itoa(seconds))
}

// Location 追加 --location 选项。
func (c *Command) Location() *Command {
	return c.option("--location")
}

// Compressed 追加 --compressed 选项。
func (c *Command) Compressed() *Command {
	return c.option("--compressed")
}

// Insecure 追加 --insecure 选项。
func (c *Command) Insecure() *Command {
	return c.option("--insecure")
}

// Include 追加 --include 选项。
func (c *Command) Include() *Command {
	return c.option("--include")
}

// Verbose 追加 --verbose 选项。
func (c *Command) Verbose() *Command {
	return c.option("--verbose")
}

// Silent 追加 --silent 选项。
func (c *Command) Silent() *Command {
	return c.option("--silent")
}

func (c *Command) option(key string, values ...string) *Command {
	c.options = append(c.options, &Option{key: key, values: values})
	return c
}

// Encode 将命令渲染为经过 shell 转义的字符串。
func (c *Command) Encode() string {
	var buffer = &bytes.Buffer{}
	for _, option := range c.options {
		option.write(buffer)
	}
	return buffer.String()
}

var replacer = strings.NewReplacer("'", "'\\''")

func escape(buffer *bytes.Buffer, values ...string) {
	buffer.WriteByte('\'')
	for _, value := range values {
		buffer.WriteString(replacer.Replace(value))
	}
	buffer.WriteByte('\'')
}
