#!/bin/sh

echo "Start network"
minifab netup -s couchdb -e true -i 2.4.8 -o manufacturer.auto.com
sleep 5

echo "Create channel"
minifab create -c autochannel
sleep 2

echo "Join channel"
minifab join -c autochannel
sleep 2

echo "Anchor update"
minifab anchorupdate
sleep 2