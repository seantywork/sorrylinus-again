package controller

import (
	"fmt"

	pkgstream "github.com/lineworld-lab/go-tv/pkg/stream"
)

type STREAM_CTL struct {
}

func (streamctl *STREAM_CTL) StartPeers() error {

	srv, err := pkgstream.CreateStreamServerForPeers()

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

	srv, err := pkgstream.CreateStreamServerForCCTV()

	if err != nil {

		return fmt.Errorf("start: %s", err.Error())

	}

	if err := srv.Run(":8080"); err != nil {

		return fmt.Errorf("start: %s", err.Error())

	}

	return nil

}
