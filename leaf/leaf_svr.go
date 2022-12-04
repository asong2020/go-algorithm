package main

import (
	"log"
	"sync"

	"asong.cloud/go-algorithm/leaf/config"
	"asong.cloud/go-algorithm/leaf/wire"
)

type LeafSvr struct {
	conf *config.Server
	once sync.Once
}

func (s *LeafSvr) init() {
	s.once.Do(func() {
		conf := &config.Server{}
		err := conf.Load("./conf/config.yaml")
		if err != nil {
			log.Panic("load conf file failed", err)
		}
		s.conf = conf
	})
}

func (u *LeafSvr) Run() {
	handler := wire.InitializeHandler(u.conf)
	handler.Run()
}
