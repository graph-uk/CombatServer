#!/bin/bash
curl http://malibu-server:3133/assets/dist/worker/malibu-worker --output malibu-worker
chmod +x malibu-worker
export GOPATH=/worker/gopath
/opt/bin/entry_point.sh&
./malibu-worker http://malibu-server:3133