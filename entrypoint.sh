#!/bin/sh

set -eo pipefail

if [ -z $GITHUB_OUTPUT ];
    then /gha-private-auth "$1" $2 $3
    else echo token=$(/gha-private-auth "$1" $2 $3) >> $GITHUB_OUTPUT;
fi
