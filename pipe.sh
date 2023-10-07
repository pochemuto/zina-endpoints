#!/bin/bash
# Create pipe: mkfifo /endpoints-pipe
# Put this file somewhere on host

while true; do
	cmd="$(cat /endpoints-pipe)"
	case $cmd in

	shutdown)
		echo "Shutdown..."
		;;

	*)
		echo "unknown '$cmd'"
		;;
	esac
done