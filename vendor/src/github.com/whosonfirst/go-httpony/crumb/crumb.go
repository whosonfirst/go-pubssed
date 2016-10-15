package crumb

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type Context interface {
	Foo() []string
}

type CommandLineContext struct {
	Context
}

func NewCommandLineContext() (*CommandLineContext, error) {

	ctx := CommandLineContext{}
	return &ctx, nil
}

func (ctx *CommandLineContext) Foo() []string {

	pid := os.Getpid()
	str_pid := strconv.Itoa(pid)

	stuff := []string{
		str_pid,
	}

	return stuff
}

type WebContext struct {
	Context
	req *http.Request
}

func NewWebContext(req *http.Request) (*WebContext, error) {

	ctx := WebContext{
		req: req,
	}

	return &ctx, nil
}

func (ctx *WebContext) Foo() []string {

	stuff := []string{
		ctx.req.RemoteAddr,
	}

	return stuff
}

type Crumb struct {
	ctx    Context
	key    string
	target string
	length int
	ttl    int
}

func NewCrumb(ctx Context, key string, target string, length int, ttl int) (*Crumb, error) {

	c := Crumb{
		ctx:    ctx,
		key:    key,
		target: target,
		length: length,
		ttl:    ttl,
	}

	return &c, nil
}

func (c *Crumb) Generate() string {

	base := c.Base()
	now := time.Now().Unix()

	hash := fmt.Sprintf("%s%d", base, now)
	hash = c.Hash(hash)

	str_now := fmt.Sprintf("%d", now)

	parts := []string{
		str_now,
		hash,
		"\xE2\x98\x83",
	}

	return strings.Join(parts, "-")
}

func (c *Crumb) Validate(crumb string) (bool, error) {

	parts := strings.Split(crumb, "-")

	if len(parts) != 3 {
		return false, errors.New("invalid crumb")
	}

	t, err := strconv.Atoi(parts[0])

	if err != nil {
		return false, errors.New("failed to parse crumb timestamp")
	}

	hash := parts[1]

	if c.ttl > 0 {
		then := t + c.ttl
		now := time.Now().Unix()

		if now > int64(then) {
			return false, errors.New("crumb has expired")
		}
	}

	base := c.Base()

	test := fmt.Sprintf("%s%d", base, t)
	test = c.Hash(test)

	// to do - test one character at a time...

	if test != hash {

		// fmt.Printf("test %s != hash %s\n", test, hash)
		return false, errors.New("crumb does not match")
	}

	return true, nil
}

func (c *Crumb) Base() string {

	parts := make([]string, 0)
	parts = append(parts, c.key)
	parts = append(parts, c.target)

	for _, v := range c.ctx.Foo() {
		parts = append(parts, v)
	}

	// fmt.Println(parts)
	return strings.Join(parts, ":")
}

func (c *Crumb) Hash(raw string) string {

	key := []byte(c.key)
	h := hmac.New(sha1.New, key)
	h.Write([]byte(raw))
	enc := hex.EncodeToString(h.Sum(nil))

	return enc[0:c.length]
}
