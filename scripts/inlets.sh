#!/bin/bash

export REMOTE="91.211.152.115:8080"    # for testing inlets on your laptop, replace with the public IPv4
export TOKEN="4bd0c17ad13a5ba0a7065f4addf648f3dad09eea"  # the client token is found on your VPS or on start-up of "inlets server"
inlets client \
 --remote=$REMOTE \
 --upstream=http://127.0.0.1:8000 \
 --token $TOKEN