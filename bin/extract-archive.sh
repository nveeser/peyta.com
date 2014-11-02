#!/bin/bash

infile=${1%/}
outfile=${2%/}

outdir=${outfile%.*}
  
#mkdir -p "${outdir}.TEMP"
case ${infile##*.} in 
  zip|ZIP)
    echo unzip -o -j -d "${outdir}.TEMP" "$filename" || exit 1
	;;
  rar)
    echo unrar e "$filename" "${outdir}.TEMP" || exit 1
	;;
esac
#mv "${dir}.TEMP" "${dir}"

