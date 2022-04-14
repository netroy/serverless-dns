import { Construct } from 'constructs'
import {
  aws_lambda,
  Stack, StackProps,
} from 'aws-cdk-lib'

import { HttpApi } from '@aws-cdk/aws-apigatewayv2-alpha'
import { HttpLambdaIntegration } from '@aws-cdk/aws-apigatewayv2-integrations-alpha'
import { resolve } from 'path'

export class DOHStack extends Stack {
  constructor(scope:Construct, id: string, props?: StackProps) {
    super(scope, id, props)

    const lambda = new aws_lambda.Function(this, 'DNSHandler', {
      code: aws_lambda.Code.fromAsset(resolve(__dirname, '../lambda/app.zip')),
      handler: 'app',
      runtime: aws_lambda.Runtime.GO_1_X
    })

    new HttpApi(this, 'Api', {
      apiName: 'DOH',
      defaultIntegration: new HttpLambdaIntegration('default', lambda)
    })
  }
}
