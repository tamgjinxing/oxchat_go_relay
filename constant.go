package main

import (
	"log"

	"github.com/fiatjaf/eventstore/lmdb"
	"github.com/fiatjaf/khatru"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	db = &lmdb.LMDBBackend{}

	relay = khatru.NewRelay()

	logger *log.Logger

	loggerWrite *lumberjack.Logger
)
