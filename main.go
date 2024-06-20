package main

import (
	"fmt"

	solictl "github.com/seantywork/sorrylinus-again/ctl"
)

func main() {

	err := solictl.LoadConfig()

	if err != nil {

		fmt.Println(err.Error())

		return
	}

	server := solictl.CreateServer()

	server.Run(solictl.CONF.ServeAddr + ":" + fmt.Sprintf("%d", solictl.CONF.ServePort))

}
