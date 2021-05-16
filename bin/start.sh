#!/bin/sh
if [ -e ./tmp ]; then
  :
else
  mkdir ./tmp
fi
./seaweedfs/weed server \
  -dir ./tmp \
  -master.dir ./tmp \
  -volume.dir.idx ./tmp \
  -ip localhost \
  -ip.bind 0.0.0.0 \
  -filer
