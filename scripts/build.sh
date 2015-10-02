#!/bin/bash

set -e

echo $EB_CONFIG > /eb_config
base64 -d /eb_config > /deploy/.ebextensions/eb.config
cd /deploy
zip -X -D -r -9 ../$CI_BRANCH-$CI_COMMIT_ID.zip ./
aws s3 cp ../$CI_BRANCH-$CI_COMMIT_ID.zip s3://$AWS_S3_BUCKET/prod/$CI_BRANCH-$CI_COMMIT_ID.zip
