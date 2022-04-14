#!/usr/bin/env node
import { App } from 'aws-cdk-lib'
import { DOHStack } from '../lib/doh-stack'

const app = new App()
new DOHStack(app, 'DOHStack')
