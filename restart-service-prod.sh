#!/bin/bash
cp helloworld-api.service /etc/systemd/system/. && systemctl enable helloworld-api.service && systemctl daemon-reload && systemctl restart helloworld-api.service && systemctl status helloworld-api.service
