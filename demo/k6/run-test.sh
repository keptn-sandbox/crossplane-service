#!/usr/bin/env sh

cd /keptn/k6

k6 run --duration 60s --vus 30 script.js