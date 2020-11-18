#!/bin/bash
#use update-worker.sh <host>
#host should contain protocol and port. Like http://127.0.0.1:3133
echo $1/assets/dist/worker/malibu-worker
curl $1/assets/dist/worker/malibu-worker --output malibu-worker
chmod +x malibu-worker
export GOPATH=/worker/gopath
/opt/bin/entry_point.sh& #run selenium-server
./malibu-worker $1