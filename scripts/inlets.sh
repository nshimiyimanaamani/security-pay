#!/bin/bash

export REMOTE="wss://inlets.quarksgroup.com"    # for testing inlets on your laptop, replace with the public IPv4
export TOKEN="4bd0c17ad13a5ba0a7065f4addf648f3dad09eea"  # the client token is found on your VPS or on start-up of "inlets server"
export UPSTREAM="inlets.quarksgroup.com=http://127.0.0.1:8000"
inlets client \
 --remote=$REMOTE \
 --upstream=$UPSTREAM \
 --token $TOKEN