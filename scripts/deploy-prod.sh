#!/bin/bash

set -e

aws elasticbeanstalk create-application-version --application-name "$AWS_APPLICATION" --version-label "$CI_COMMIT_ID" --source-bundle S3Bucket="$AWS_S3_BUCKET",S3Key="prod/$CI_BRANCH-$CI_COMMIT_ID.zip"
aws elasticbeanstalk update-environment --environment-name "$AWS_PROD_ENVIRONMENT" --version-label "$CI_COMMIT_ID"
