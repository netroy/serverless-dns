# Serverless DNS over HTTPS

DNS-over-HTTPS is the future of every Privacy-Conscious person on the internet.
Unfortunately, there aren't many public DOH servers that guarantee privacy.

This application uses serverless technologies from AWS to run an extremely low-cost & zero-maintaince DNS server.

## Deployment

Currenty this repo isn't completely automated. Please follow these steps:

### Requirements

- gvm with go 1.11+
- nodejs 10+
- yarn
- aws cli installed & configured with correct credentials

### Setup

- run `make setup build` in the lambda folder, to prepare an uploadable zip file for AWS Lambda
- run `yarn && cdk deploy` in the root folder
