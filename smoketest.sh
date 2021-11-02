#!/bin/bash
IP=$(docker inspect -f '{{range.NetworkSettings.Networks}}{{.IPAddress}}{{end}}' dta)
PORT=$(docker port dta | grep -o [0-9]*$)
docker run --env IP=$IP --env PORT=$PORT benpeters/dns-threat-analyser:latest ./SmokeTestClient