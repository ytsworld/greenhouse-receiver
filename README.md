# greenhouse-receiver

This is the server-side counterpart to the [raspberry device](https://github.com/ytsworld/greenhouse-client) that sends sensor data.
The implementation is done with a cloud function deployed at Google Cloud Platform and creates custom stackdriver metrics to store metrics.

## Authentication
The cloud function expects an identity token to be part of the request header otherwise it will respond with `403 Forbidden`.

To trigger the deployed function the token can be created using the gcloud tools:
`curl -H "Authorization: Bearer $(gcloud auth print-identity-token)" https://[region]-[project_id].cloudfunctions.net/greenhouse-receiver`

[Client devices](https://github.com/ytsworld/greenhouse-client) use a service account key to create & refresh id tokens.
A new service account key can be created using this script:

```
GCP_PROJECT=[your gcp project id]
export GCP_PROJECT
./scripts/createClientSA.sh
```

## Local testing
To run the tests create a service account with permissions to write monitoring metrics:
```
GCP_PROJECT=[your gcp project id]
export GCP_PROJECT
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

```
GCP_PROJECT=[your gcp project id]
export GCP_PROJECT
./scripts/deploy.sh
```
