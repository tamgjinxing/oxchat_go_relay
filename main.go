package main

import (
	"log"
	"net/http"
	"time"

	_ "net/http/pprof"

	"gopkg.in/natefinch/lumberjack.v2"
)

func init() {
	loggerWrite = &lumberjack.Logger{
		Filename:   "./log/oxchat_relay.log", //log file path
		MaxSize:    500,                      // Maximum size of each log file (MB)
		MaxBackups: 7,                        // Maximum number of old log files to keep
		MaxAge:     1,                        // Maximum number of days to keep old log files
		Compress:   false,                    // Whether to archive old log files
		LocalTime:  true,
	}

	logger = log.New(loggerWrite, "INFO ", log.LstdFlags|log.Lshortfile)
}

func main() {
	// read config
	err := ReadConfig("config.json")
	if err != nil {
		logger.Fatalf("Error reading config file:%v\n", err)
	}

	// load db
	db.Path = config.RelayInfo.DatabasePath
	if err := db.Init(); err != nil {
		logger.Fatalf("failed to initialize database:%v\n", err)
		return
	}

	// init relay
	relay.Info.Name = config.RelayInfo.RelayName
	relay.Info.PubKey = config.RelayInfo.RelayPubkey
	relay.Info.Description = config.RelayInfo.RelayDescription
	relay.Info.Contact = config.RelayInfo.RelayContact
	relay.Info.Icon = config.RelayInfo.RelayIcon
	relay.Info.SupportedNIPs = config.RelayInfo.SupportNips

	relay.StoreEvent = append(relay.StoreEvent, DBSaveEventLogic)
	relay.OnEphemeralEvent = append(relay.OnEphemeralEvent, PrintHeartEvent)

	relay.QueryEvents = append(relay.QueryEvents,
		db.QueryEvents,
	)
	relay.DeleteEvent = append(relay.DeleteEvent, db.DeleteEvent)

	relay.RejectEvent = append(relay.RejectEvent,
		PreventTimestampsFor1059(60),
	)

	relay.OnConnect = append(relay.OnConnect,
		OnConnectLogic,
	)

	relay.OnDisconnect = append(relay.OnDisconnect, OnDisconnectLogic)

	changeFile()

	go func() {
		http.ListenAndServe(":6060", nil)
	}()

	logger.Printf("relay-pubkey:%s, running on http://0.0.0.0:%s\n", config.RelayInfo.RelayPubkey, config.RelayInfo.Port)

	if err := http.ListenAndServe(":"+config.RelayInfo.Port, relay); err != nil {
		logger.Fatalf("Failed to start server:%v\n", err)
	}
}

func changeFile() {
	go func() {
		for {
			nowTime := time.Now()
			nowTimeStr := nowTime.Format("2006-01-02")
			t2, _ := time.ParseInLocation("2006-01-02", nowTimeStr, time.Local)
			next := t2.AddDate(0, 0, 1)
			after := next.UnixNano() - nowTime.UnixNano() - 1
			<-time.After(time.Duration(after) * time.Nanosecond)
			loggerWrite.Rotate()
		}
	}()
}
