package main

import (
	"fmt"
	vhl "github.com/vrherog/VirtualHereLibrary"
)

func main() {
	var client = vhl.NewClient()

	for i:=0 ;i < 100; i++ {
		go test(client)
	}

	var cmd string
	_, _ = fmt.Scanf(`Press any key...%s`, &cmd)
}

func test(client *vhl.Client) {
	var result, err = client.List()
	if err == nil {
		fmt.Printf("%+v\n", result)
	} else {
		fmt.Println(err)
	}
}
