# Serverless DNS over HTTPS

DNS-over-HTTPS is the future of every Privacy-Conscious person on the internet.
Unfortunately, there aren't many public DOH servers that guarantee privacy.

This application uses serverless technologies from AWS to run an extremely low-cost & zero-maintaince DNS server.

### Requirements

- gvm with go 1.18+
- nodejs 14+
- yarn
- aws cli installed & configured with correct credentials

## Deployment

- run `yarn && yarn cdk deploy` in the root folder
