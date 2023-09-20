package main

import (
	"flag"
	"fmt"
	"strconv"
	"strings"

	"github.com/fantasyczl/hermes/proxy"
)

func main() {
	point, err := parseFlags()
	if err != nil {
		fmt.Println(err)
		return
	}

	if err = point.StartServe(); err != nil {
		fmt.Println(err)
		return
	}
}

func parseFlags() (e proxy.Endpoint, err error) {
	var p string
	flag.StringVar(&p, "p", "", "addr")
	flag.Parse()

	fmt.Printf("p: %s\n", p)

	if p == "" {
		flag.Usage()
		err = fmt.Errorf("p is empty")
		return
	}

	return parseStrToEndpoint(p)
}

func parseStrToEndpoint(s string) (e proxy.Endpoint, err error) {
	list := strings.Split(s, ":")
	if len(list) != 3 {
		err = fmt.Errorf("invalid endpoint: %s", s)
		return
	}

	e.LocalPort, err = strconv.Atoi(list[0])
	if err != nil {
		return
	}

	e.RemoteHost = list[1]
	e.RemotePort, err = strconv.Atoi(list[2])
	if err != nil {
		return
	}

	if err = e.Valid(); err != nil {
		err = fmt.Errorf("invalid endpoint: %s", err)
		return
	}

	return
}
