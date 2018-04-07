#! /bin/bash

go build nuance.go
ssh -t berlin@95.85.49.5 "sudo service nuance stop"
scp ./nuance berlin@95.85.49.5:/home/berlin/nuance_bot/
scp ./config.json berlin@95.85.49.5:/home/berlin/nuance_bot/config.json
ssh -t berlin@95.85.49.5 "sudo service nuance start"
rm nuance
