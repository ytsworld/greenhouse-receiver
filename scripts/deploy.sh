#! /bin/bash
set -e

go test

gcloud beta functions deploy greenhouse-receiver --entry-point EntryPoint --runtime go111 --trigger-http --region europe-west1 --max-instances 2 --memory 128M --no-allow-unauthenticated
