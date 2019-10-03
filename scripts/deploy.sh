#! /bin/bash
set -e

# Service account to be used by function
sa_name="greenhouse-receiver-function"
sa_id="${sa_name}@${PROJECT_ID}.iam.gserviceaccount.com"

# Check if service account exists and create if necessary
sa_exists=$(gcloud iam service-accounts list | grep -c "$sa_id" || echo -n "")

if [ "${sa_exists}" == "0" ]; then
    echo "Creating service account for function"
    
    gcloud iam service-accounts create "${sa_name}" --display-name "Greenhouse receiver SA for cloud function" 

    gcloud projects add-iam-policy-binding "$PROJECT_ID" --member "serviceAccount:${sa_id}" --role "roles/monitoring.metricWriter"

fi

echo "Deploying function ..."

gcloud beta functions deploy greenhouse-receiver \
    --entry-point EntryPoint \
    --runtime go111 \
    --trigger-http \
    --region europe-west1 \
    --max-instances 2 \
    --memory 128M \
    --no-allow-unauthenticated \
    --service-account "${sa_id}"

echo "Deploying finished"
