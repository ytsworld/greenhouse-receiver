#! /bin/bash
set -e

if [ -z "$GCP_PROJECT" ]; then
    echo "GCP project id is required in environment variable GCP_PROJECT"
    exit 1
fi

echo "--------------- Client account ----------------"

# TODO every device should have its unique identity
sa_name="greenhouse-client"
sa_id="${sa_name}@${GCP_PROJECT}.iam.gserviceaccount.com"

# Create service account
gcloud iam service-accounts create "${sa_name}" --display-name "Greenhouse Client" 

# Create private key for account
gcloud iam service-accounts keys create "./secrets/${sa_name}.sa.json" --iam-account "${sa_id}"

# Provide service account access to cloud function
gcloud beta functions add-iam-policy-binding greenhouse-receiver \
  --member="serviceAccount:${sa_id}" \
  --role='roles/cloudfunctions.invoker' \
  --region europe-west1
