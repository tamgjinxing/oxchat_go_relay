package main

import (
	"log"

	"github.com/fiatjaf/eventstore/lmdb"
	"github.com/fiatjaf/khatru"
)

var (
	db = &lmdb.LMDBBackend{}

	relay = khatru.NewRelay()

	logger *log.Logger
)
