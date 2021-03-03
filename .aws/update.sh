#!/bin/bash

PROJECT=crane

aws cloudformation update-stack --stack-name crane-gh-actions-user \
  --template-body file://main.yml \
  --parameters ParameterKey=AssetsBucket,ParameterValue=$ASSETS_BUCKET \
  ParameterKey=ExternalId,ParameterValue=$EXTERNAL_ID \
  ParameterKey=ProjectName,ParameterValue=$PROJECT \
  --tags Key=project,Value=$PROJECT \
  --region ap-southeast-2 \
  --capabilities CAPABILITY_NAMED_IAM \
  --output yaml \
  --profile $PROFILE
