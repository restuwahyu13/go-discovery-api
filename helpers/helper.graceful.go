package helpers

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

const (
	TCP        = "tcp"
	UDPConn    = "udp"
	UNIXConn   = "unix"
	UNIXPACKET = "unixpacket"
)

type GracefulOptions struct {
	Server  *grpc.Server
	Address string
	Port    string
}

func Graceful(options *GracefulOptions) {
	nlc := net.ListenConfig{KeepAlive: 0}

	nls, err := nlc.Listen(context.Background(), TCP, fmt.Sprintf("%s:%s", options.Address, options.Port))
	if err != nil {
		logrus.Fatal(err)
	}

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGTERM, syscall.SIGINT, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGALRM, syscall.SIGABRT, syscall.SIGUSR1)

	select {
	case <-signalChan:
		defer nls.Close()
		os.Exit(0)

	default:
		logrus.Print("\n")
		logrus.Infof("Server listening on port %s", options.Port)

		if err := options.Server.Serve(nls); err != nil {
			logrus.Fatal(err)
		}
	}
}
