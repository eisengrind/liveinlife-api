#!/bin/sh

# get the github key
ssh-keyscan github.com >> ~/.ssh/known_hosts

# setup git stuff
git config --global user.email "$GITHUB_BOT_EMAIL"
git config --global user.name "$GITHUB_BOT_NAME"

# create temp location for api libs
mkdir -p /tmp/lib

# generate csharp library
git clone git@github.com:51st-state/cs-api-client.git /tmp/lib/cs
java -jar /openapi-generator-cli.jar generate \
    -g csharp-netcore \
    -i $PROJECT_ROOT/api/openapi-spec/openapi.json \
    -o /tmp/lib/cs/src/ \
    --additional-properties=modelPropertyNaming=original \
    --additional-properties=packageName=FF \
    --additional-properties=clientPackage=FF.Client
cd /tmp/lib/cs
git add -A
git commit -a -m "generate csharp library\ncommit reference 51st-state/api@$SHORT_SHA"
git push
