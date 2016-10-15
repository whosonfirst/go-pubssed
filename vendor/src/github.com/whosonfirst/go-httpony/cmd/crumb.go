package main

import (
	"fmt"
	"github.com/whosonfirst/go-httpony/crumb"
)

func main() {

	ctx, _ := crumb.NewCommandLineContext()

	key := "foo"
	target := "bar"
	length := 10
	ttl := 600

	c, _ := crumb.NewCrumb(ctx, key, target, length, ttl)
	cr := c.Generate()

	fmt.Printf("crumb is %s\n", cr)

	ok, err := c.Validate(cr)

	fmt.Printf("crumb validates %t\n", ok)

	if err != nil {
		fmt.Println(err)
	}
}
