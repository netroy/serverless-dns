import { CfnApi, CfnDomainName, CfnApiMapping } from '@aws-cdk/aws-apigatewayv2'
import { Certificate } from '@aws-cdk/aws-certificatemanager'
import { ServicePrincipal } from '@aws-cdk/aws-iam'
import { Code, Function as Lambda, Runtime } from '@aws-cdk/aws-lambda'
import { Construct, Stack, StackProps } from '@aws-cdk/core'
import { resolve } from 'path'

export class DOHStack extends Stack {
  constructor(scope:Construct, id: string, props?: StackProps) {
    super(scope, id, props)

    const lambda = new Lambda(this, 'DNSHandler', {
      code: Code.fromAsset(resolve(__dirname, '../lambda/app.zip')),
      handler: 'app',
      runtime: Runtime.GO_1_X
    })

    lambda.grantInvoke(new ServicePrincipal('apigateway.amazonaws.com'))

    const api = new CfnApi(this, 'Resource', {
      name: 'DOH',
      protocolType: 'HTTP',
      target: lambda.functionArn
    })

    const domainName = 'doh.netroy.in'
    const certificate = new Certificate(this, 'Certificate', { domainName })
    new CfnDomainName(this, 'CustomDomain', {
      domainName,
      domainNameConfigurations: [{
        certificateArn: certificate.certificateArn,
        endpointType : 'REGIONAL'
      }]
    })

    const mapping = new CfnApiMapping(this, 'ApiMapping', {
      apiId: api.ref,
      domainName,
      stage: '$default'
    })
  }
}
