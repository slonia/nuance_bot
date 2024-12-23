#! /bin/bash

CGO_ENABLED=0 go build nuance.go
ssh -t server@ams_server "supervisorctl stop nuance_bot"
scp ./nuance server@ams_server:/home/server/apps/nuance_bot/
scp ./config.json server@ams_server:/home/server/apps/nuance_bot/config.json
ssh -t server@ams_server "supervisorctl start nuance_bot"
rm nuance
