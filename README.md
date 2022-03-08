# twitter-api #

[![GoDoc](https://img.shields.io/static/v1?label=godoc&message=reference&color=blue)](https://pkg.go.dev/github.com/inpublic-io/twitter-api?tab=doc)

twitter-api is a grpc based API to attend inpublic.io Twitter integration needs

## Usage ##

## Deploy on GCP ##

```shell
# create the secret on GCP secrets manager
gcloud secrets create twitter-bearer-token \
    --replication-policy="automatic"

# create a secret from a string
printf "twitter_bearer_token_here" | gcloud secrets versions add twitter-bearer-token \
            --data-file=-

# deploy to cloud run
gcloud run deploy twitter-api \
    --image gcr.io/<google_project_id>/twitter-api \
    --update-secrets=TWITTER_BEARER_TOKEN=twitter-bearer-token:latest
```

## License ##

This library is distributed under the MIT license found in the [LICENSE](./LICENSE)
file.
