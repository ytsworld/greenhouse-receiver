#! /bin/bash
set -e

if [ -z "$GCP_PROJECT" ]; then
    echo "GCP project id is required in environment variable GCP_PROJECT"
    exit 1
fi


echo "----------- Local testing account -------------"

sa_name="greenhouse-receiver-local"
sa_id="${sa_name}@${GCP_PROJECT}.iam.gserviceaccount.com"

# Create service account
gcloud iam service-accounts create "${sa_name}" --display-name "Greenhouse receiver SA for local tests" 

# Create private key for account
gcloud iam service-accounts keys create "./secrets/${sa_name}.sa.json" --iam-account "${sa_id}"

# Provide service account write access to metrics
gcloud projects add-iam-policy-binding "$GCP_PROJECT" --member "serviceAccount:${sa_id}" --role "roles/monitoring.metricWriter"
