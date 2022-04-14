#!/usr/bin/env node
import { App } from 'aws-cdk-lib'
import { Construct } from 'constructs'
import { Stack } from 'aws-cdk-lib'
import { GoFunction } from '@aws-cdk/aws-lambda-go-alpha'
import { HttpApi, HttpMethod } from '@aws-cdk/aws-apigatewayv2-alpha'
import { HttpLambdaIntegration } from '@aws-cdk/aws-apigatewayv2-integrations-alpha'

const name = 'DoH'

class DoHStack extends Stack {
  constructor(scope:Construct) {
    super(scope, name)

    const lambda = new GoFunction(this, 'Handler', {
      functionName: name,
      entry: 'lambda',
    })

    const api = new HttpApi(this, 'Api', {
      apiName: name,
    })

    // TODO: add support for JSON API for DoH
    api.addRoutes({
      methods: [
        HttpMethod.GET,
        // HttpMethod.POST, // TODO: Add support for POST
      ],
      path: '/dns-query',
      integration: new HttpLambdaIntegration('dns-query', lambda),
    })
  }
}


const app = new App()
new DoHStack(app)
