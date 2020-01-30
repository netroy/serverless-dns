#!/usr/bin/env node
import { App } from '@aws-cdk/core'
import { DOHStack } from '../lib/doh-stack'

const app = new App()
new DOHStack(app, 'DOHStack')
