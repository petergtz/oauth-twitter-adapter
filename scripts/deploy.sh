#!/bin/bash -ex

cd $(dirname $0)/..

cf login -a api.eu-de.bluemix.net --sso -o $(lpass show Personal\\api_keys/Alexa-Wikipedia-Skill --notes) -s alexa

cf push --no-start

for line in $(lpass show Personal\\api_keys/oauth-twitter-adapter --notes); do
    cf set-env alexa-skill-linked-account-twitter-middleware $line
done

cf restart alexa-skill-linked-account-twitter-middleware
