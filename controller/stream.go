package controller

import (
	"fmt"

	pkgstream "github.com/lineworld-lab/go-boomerland/pkg/stream"
)

type STREAM_CTL struct {
	TurnServerAddr string
}

func (streamctl *STREAM_CTL) StartPeers() error {

	pkgstream.TurnServerAddr = streamctl.TurnServerAddr

	srv, err := pkgstream.CreateStreamServerForPeers()

	if err != nil {

		return fmt.Errorf("start: %s", err.Error())

	}

	if err := srv.Run(":8080"); err != nil {

		return fmt.Errorf("start: %s", err.Error())

	}

	return nil
}

func (streamctl *STREAM_CTL) StartRoom() error {

	srv, err := pkgstream.CreateStreamServerForRoom()

	if err != nil {

		return fmt.Errorf("start: %s", err.Error())

	}

	if err := srv.Run(":8080"); err != nil {

		return fmt.Errorf("start: %s", err.Error())

	}

	return nil
}

func (streamctl *STREAM_CTL) StartFiles() error {

	srv, err := pkgstream.CreateStreamServerForFiles()

	if err != nil {

		return fmt.Errorf("start: %s", err.Error())

	}

	if err := srv.Run(":8080"); err != nil {

		return fmt.Errorf("start: %s", err.Error())

	}

	return nil
}

func (streamctl *STREAM_CTL) StartCCTV() error {

	pkgstream.TurnServerAddr = streamctl.TurnServerAddr

	srv, err := pkgstream.CreateStreamServerForCCTV()

	if err != nil {

		return fmt.Errorf("start: %s", err.Error())

	}

	if err := srv.Run(":8080"); err != nil {

		return fmt.Errorf("start: %s", err.Error())

	}

	return nil

}
