# greenhouse-receiver

This is the server-side counterpart to the [https://github.com/ytsworld/greenhouse-client](raspberry device) that sends sensor data.
The implementation is done based on a cloud function provided by Google Cloud Platform and creates custom stackdriver metrics to store metrics.

## Authentication
The cloud function expects a identity token to be part of the request header otherwise it will respond with `403 Forbidden`.

For local tests the token can be created using the gcloud tools:
`curl -H "Authorization: Bearer $(gcloud auth print-identity-token)" https://[region]-[project_id].cloudfunctions.net/greenhouse-receiver`

[https://github.com/ytsworld/greenhouse-client](A client device) use a service account key to create & refresh id tokens.
A new service account key can be created using this script:

```
PROJECT_ID=[your gcp project id]
./scripts/createClientSA.sh
```

## Local testing
To run the integration tests locally create a service account with permissions to write monitoring metrics:
```
PROJECT_ID=[your gcp project id]
./scripts/createLocalReceiverSA.sh
```

Then set these variable and run the test:
```
GOOGLE_APPLICATION_CREDENTIALS=$(pwd)/secrets/greenhouse-receiver-local.sa.json
export GOOGLE_APPLICATION_CREDENTIALS
GCP_PROJECT=[your gcp project id]
export GCP_PROJECT
go test
```

## Deploy
For deployment use these script:
```
./scripts/deploy.sh
```
