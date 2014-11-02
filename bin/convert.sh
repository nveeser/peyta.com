#!/bin/sh

IN=$1
OUT=$2
FFMPEG=/home/nicholas/bin/ffmpeg
BITRATE=1800k


find $IN -type f | egrep '(mov|wmv|flv|mpv)$' | while read infile; do
	outfile=${infile/$IN/$OUT}	
	outfile=${outfile%.*}.mp4
	echo "IN:  $infile"
	echo "OUT: $outfile"
done

#${FFMPEG} -y -i "${IN}" -c:v libx264 -preset veryslow -b:v ${BITRATE} -pass 1 -an -f mp4 /dev/null || exit 1
#${FFMPEG} -i "${IN}" -c:v libx264 -preset veryslow -b:v ${BITRATE} -pass 2 -c:a libfdk_aac -b:a 128k "${OUT}"
