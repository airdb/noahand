#!/bin/sh

echo "install noah"
WORKDIR="/opt/noah/bin"

_myos="$(uname)"
case $_myos in
    Darwin)
	mkdir -p $WORKDIR
	cp release/noah-darwin $WORKDIR
	cp service/io.airdb.noah.plist /Library/LaunchDaemons/
	launchctl load  /Library/LaunchDaemons/io.airdb.noah.plist
	launchctl start io.airdb.noah
        ;;
    Linux)
	mkdir -p $WORKDIR
	cp release/noah-linux $WORKDIR/noah
	cp service/noah.service /etc/systemd/system/
	systemctl daemon-reload
	systemctl start noah
        ;;
	*) 
	echo "not support"
	;;
esac

echo "finish"
