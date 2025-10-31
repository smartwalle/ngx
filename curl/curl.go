package curl

import (
	"bytes"
	"strings"
)

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

type Command struct {
	options []*Option
}

func New(method, url string) *Command {
	var cmd = &Command{
		options: make([]*Option, 0, 6),
	}
	cmd.Option("curl")
	cmd.Option("--request")
	cmd.Option(method)
	cmd.Option("", url)
	return cmd
}

func (c *Command) Header(key, value string) *Command {
	return c.Option("--header", key+":"+value)
}

func (c *Command) Data(data string) *Command {
	return c.Option("--data", data)
}

func (c *Command) Form(key, value string) *Command {
	return c.Option("--form", key+"="+value)
}

func (c *Command) FormFile(key, filepath string) *Command {
	return c.Option("--form", key+"=@"+filepath)
}

func (c *Command) User(username, password string) *Command {
	return c.Option("--user", username+":"+password)
}

func (c *Command) UserAgent(value string) *Command {
	return c.Option("--user-agent", value)
}

func (c *Command) Option(key string, values ...string) *Command {
	c.options = append(c.options, &Option{key: key, values: values})
	return c
}

func (c *Command) Encode() string {
	var buffer = &bytes.Buffer{}
	for _, option := range c.options {
		option.write(buffer)
	}
	return buffer.String()
}

var replacer = strings.NewReplacer("\\", "\\\\", "\"", "\\\"")

func escape(buffer *bytes.Buffer, values ...string) {
	buffer.WriteByte('"')
	for _, value := range values {
		buffer.WriteString(replacer.Replace(value))
	}
	buffer.WriteByte('"')
}
