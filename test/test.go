package main

import (
	"fmt"
	"github.com/zdnscloud/node-agent/command"
)

func main() {
	infos, err := command.GetDisksInfo("")
	fmt.Println(err)
	for k, v := range infos {
		fmt.Println(k, v)
	}
}
