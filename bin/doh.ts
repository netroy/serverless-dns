#!/usr/bin/env node
import { App, CfnParameter } from 'aws-cdk-lib'
import { Construct } from 'constructs'
import { Stack } from 'aws-cdk-lib'
import { Certificate, CertificateValidation } from 'aws-cdk-lib/aws-certificatemanager'
import { GoFunction } from '@aws-cdk/aws-lambda-go-alpha'
import { DomainName, EndpointType, HttpApi, HttpMethod } from '@aws-cdk/aws-apigatewayv2-alpha'
import { HttpLambdaIntegration } from '@aws-cdk/aws-apigatewayv2-integrations-alpha'

const name = 'DoH'

class DoHStack extends Stack {
  constructor(scope:Construct) {
    super(scope, name)

    const domainName = (new CfnParameter(this, 'DomainName')).valueAsString
    console.log(domainName)

    const certificate = new Certificate(this, 'Certificate', {
      domainName,
      validation: CertificateValidation.fromDns()
    })

    const domain = new DomainName(this, 'Domain', {
      domainName,
      certificate,
      endpointType: EndpointType.REGIONAL,
    })

    const lambda = new GoFunction(this, 'Handler', {
      functionName: name,
      entry: 'lambda',
    })

    const api = new HttpApi(this, 'Api', {
      apiName: name,
      defaultDomainMapping: {
        domainName: domain,
      }
    })

    // TODO: add support for JSON API for DoH
    api.addRoutes({
      methods: [
        HttpMethod.GET,
        HttpMethod.POST,
      ],
      path: '/dns-query',
      integration: new HttpLambdaIntegration('dns-query', lambda),
    })
  }
}


const app = new App()
new DoHStack(app)
