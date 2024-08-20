#!/bin/bash
APP_NAME="oxchat_go_relay"

BIN_DIR=`pwd`
DEPLOY_DIR=`pwd`

SHELL_NAME_START=./start.sh
SHELL_NAME_BUILD=./build.sh

LOGS_FILE=$BIN_DIR/log/oxchat_relay.log

start() {
	echo -e "logs file $LOGS_FILE"
	echo -e "Starting the $APP_NAME ..."
	echo "SHELL_NAME_START:" $SHELL_NAME_START

	nohup sh $SHELL_NAME_START >> $LOGS_FILE 2>&1 &

	sleep 10s

	PIDS=`ps -ef | grep "$APP_NAME" |grep -v grep |awk '{print $2}'`
	echo "Pid:" $PIDS
	echo "Start successful!!!"
}

build(){
	echo -e "logs file $LOGS_FILE"
	echo -e "Building the $APP_NAME ..."
	echo "SHELL_NAME_BUILD:" $SHELL_NAME_BUILD

	nohup sh $SHELL_NAME_BUILD >> $LOGS_FILE 2>&1 &

	sleep 10s

	PIDS=`ps -ef | grep "$APP_NAME" |grep -v grep |awk '{print $2}'`
	echo "Pid:" $PIDS
	echo "build and start successful!!!"
}

restart () {
    killPIDS=`ps -ef | grep "$APP_NAME" |grep -v grep |awk '{print $2}'`

	if [ -n "$killPIDS" ]; then
		echo -e "Stopping the $APP_NAME ..."
        for PID in $killPIDS ; do
    		kill $PID
        	echo "kill pid $PID"
        done
	fi

	sleep 1
	start
}

rebuild () {
    killPIDS=`ps -ef | grep "$APP_NAME" |grep -v grep |awk '{print $2}'`

	if [ -n "$killPIDS" ]; then
		echo -e "Stopping the $APP_NAME ..."
        for PID in $killPIDS ; do
    		kill $PID
        	echo "kill pid $PID"
        done
	fi

	sleep 1
	build
}

stop() {
	killPIDS=`ps -ef | grep "$APP_NAME" |grep -v grep |awk '{print $2}'`

	if [ -n "$killPIDS" ]; then
		echo -e "Stopping the $APP_NAME ..."
        for PID in $killPIDS ; do
    		kill $PID
        	echo "kill pid $PID"
        done
	fi

	echo "stop successful!!!"
}

case "$1" in

    stop|-stop|--stop)
        stop
        ;;

	restart|-restart|--restart)
        restart
		;;

	build|-build|--build)
        rebuild
		;;

	rebuild|-rebuild|--rebuild)
        rebuild
		;;

	dump|-dump|--dump)
        dump
		;;

	status|-status|--status)
        status
		;;

	help|-help|--help)
        echo ""
        echo "This script used for {start|stop|restart|dump} !"
        echo ""
        echo ""
        ;;

	start|-start|--start|*)
        start
        ;;

esac
