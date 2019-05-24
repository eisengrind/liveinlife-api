#!/bin/sh

# get the github key
ssh-keyscan github.com >> /root/.ssh/known_hosts

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
git commit -a -m "generate csharp library"
git push

# generate typescript angular library
git clone git@github.com:51st-state/ng-api-client.git /tmp/lib/ts-ng
java -jar /openapi-generator-cli.jar generate \
    -g typescript-angular \
    -i $PROJECT_ROOT/api/openapi-spec/openapi.json \
    -o /tmp/lib/ts-ng/src/ \
    --additional-properties=modelPropertyNaming=original \
    --additional-properties=npmName=@51st-state/ng-api-client \
    --additional-properties=npmVersion=1.0.0 \
    --additional-properties=ngVersion=7.0.0
cd /tmp/lib/ts-ng
git add -A
git commit -a -m "generate typescript angular library"
git push
