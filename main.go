package main

import (
	"log"
	"net/http"

	_ "net/http/pprof"

	"gopkg.in/natefinch/lumberjack.v2"
)

func init() {
	logger = log.New(&lumberjack.Logger{
		Filename:   "./log/oxchat_relay.log", // 日志文件路径
		MaxSize:    500,                      // 每个日志文件的最大尺寸 (MB)
		MaxBackups: 7,                        // 保留的旧日志文件最大数量
		MaxAge:     1,                        // 保留旧日志文件的最大天数
		Compress:   true,                     // 是否压缩归档旧日志文件
	}, "INFO: ", log.LstdFlags|log.Lshortfile)
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

	relay.OnDisconnect = append(relay.OnDisconnect)

	// 启动 pprof 服务器
	go func() {
		http.ListenAndServe(":6060", nil)
	}()

	logger.Printf("relay-pubkey:%s, running on http://0.0.0.0:%s\n", config.RelayInfo.RelayPubkey, config.RelayInfo.Port)

	if err := http.ListenAndServe(":"+config.RelayInfo.Port, relay); err != nil {
		logger.Fatalf("Failed to start server:%v\n", err)
	}
}
