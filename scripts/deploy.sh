#!/bin/bash -ex

cd $(dirname $0)/..

cf push --no-start

for line in $(lpass show Personal\\api_keys/oauth-twitter-adapter --notes); do
    echo cf set-env alexa-skill-linked-account-twitter-middleware $line
done

cf restart alexa-skill-linked-account-twitter-middleware
