package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dihedron/kraft/log"
	"github.com/jessevdk/go-flags"
	"go.uber.org/zap"
)

// also check out:
//  - https://github.com/yongman/leto
//    example with join, leave and snapshot (emulating Redis)
//  - https://github.com/yusufsyaifudin/raft-sample
//    example with HTTP endpoint to join/leave, emulating Badger

// Options is the set of command line options.
type Options struct {
	Listen   int `short:"l" long:"listen" description:"Port to listen on"`
	Desc     string
	RaftDir  string
	RaftBind string
	Join     string
	NodeID   string
}

func main() {
	defer log.L.Sync()

	// install interrupt handler
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Kill, os.Interrupt, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	// run main business login in a separate goroutine
	go func() {
		options := Options{}
		_, err := flags.Parse(&options)
		if err != nil {
			log.L.Error("error parsing command line", zap.Error(err))
			os.Exit(1)
		}

		fmt.Printf("listen: %d\n", options.Listen)
		time.Sleep(30 * time.Second)
	}()

	<-interrupt
	fmt.Println("cleaning up...")
}

/*
var (
	listen   string
	raftdir  string
	raftbind string
	nodeID   string
	join     string
)

func init() {
	flag.StringVar(&listen, "listen", ":5379", "server listen address")
	flag.StringVar(&raftdir, "raftdir", "./", "raft data directory")
	flag.StringVar(&raftbind, "raftbind", ":15379", "raft bus transport bind address")
	flag.StringVar(&nodeID, "id", "", "node id")
	flag.StringVar(&join, "join", "", "join to already exist cluster")
}

func main() {
	flag.Parse()

	var (
		c *config.Config
	)

	c = config.NewConfig(listen, raftdir, raftbind, nodeID, join)

	app := server.NewApp(c)

	quitCh := make(chan os.Signal, 1)
	signal.Notify(quitCh, os.Kill, os.Interrupt, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go app.Run()

	<-quitCh
}
*/
