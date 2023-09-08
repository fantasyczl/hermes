package main

import (
	"flag"
	"fmt"
	"net"
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

	if err = start(point); err != nil {
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

func start(ep proxy.Endpoint) error {
	fmt.Printf("start: %s\n", ep.String())

	l, err := net.Listen("tcp", ep.LocalAddr())
	if err != nil {
		return err
	}

	defer func() {
		if e := l.Close(); e != nil {
			fmt.Printf("close error: %s\n", e)
		}
	}()

	for {
		conn, lErr := l.Accept()
		if lErr != nil {
			fmt.Printf("accept error: %s\n", lErr)
			continue
		}

		go proxy.Handle(conn, ep)
	}
}
