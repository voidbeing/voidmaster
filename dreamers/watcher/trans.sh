#!/bin/bash

if [ $# != 1 ]; then
	echo please specify the video filename in videos dir.
	exit 0
fi

rm -f frames/*.png

NAME=videos/$1
if [ $1 == "new" ]; then
	NAME=videos/$(ls -lt videos/ | awk '{print $NF}' | head -2 | tail -1)
fi

ffmpeg -i $NAME frames/%05d.png -hide_banner -vf fps=6
