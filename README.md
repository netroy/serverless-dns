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

- run `yarn` to get all the node dependencies
- run `yarn cdk deploy --parameters DomainName=[[DOMAIN_NAME]]` in the root folder to deploy
- The deployment creates a ACM certificate, that needs to be validated by adding DNS records on your domain's hosted zone

### Configure Firefox
* `network.trr.bootstrapAddress` = `4.2.2.2`
* `network.trr.mode` = `3`
* `network.trr.uri` = `https://[[DOMAIN_NAME]]/dns-query`
* `network.trr.useGET` = `true`
* `network.security.esni.enabled` = `true`
