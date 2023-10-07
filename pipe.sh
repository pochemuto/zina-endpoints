#!/bin/bash
# Create pipe: mkfifo /endpoints-pipe
# Put this file somewhere on host

while true; do
	cmd="$(cat /endpoints-pipe)"
	case $cmd in

	shutdown)
		echo "Shutting down..."
		systemctl poweroff
		;;
	
	stop-pipe-reader)
	    echo "Stop pipe reader..."
		exit 0
		;;
	
	*)
		echo "unknown '$cmd'"
		;;
	esac
done