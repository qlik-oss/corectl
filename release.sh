#!/bin/bash

set -e

function print_usage {
  echo "Usage:"
  echo "  RELEASE_VERSION=v1.0.0 release.sh"
  echo "  release.sh (-?, -h, --help)"
  echo
  echo "Options:"
  echo "  -?, -h, --help  - Prints this usage information."
  echo
  echo "Environment variables:"
  echo "  - RELEASE_VERSION - Required. Set to a semantic version"
  echo
}

if [[ $1 == "-?" || $1 == "-h" || $1 == "--help" ]]; then
  print_usage
  exit 0
fi

function pre_flight_checks() {
  if [[ ! -z $(git status --porcelain) ]]; then
    echo "There are uncommitted changes. Please make sure branch is clean."
    exit 1
  fi
}

## Check that there are no uncommitted changes locally
pre_flight_checks

## Build corectl with the version number and generate an API specification
echo "Generating new API specification for version ${RELEASE_VERSION}"
go build -ldflags "-X main.version=${RELEASE_VERSION}"
./corectl generate-spec > docs/spec.json

## Create a commit and add a tag
echo "Creating commit and tag for version ${RELEASE_VERSION}"
git commit -a -m "Release ${RELEASE_VERSION}"
git tag -a $RELEASE_VERSION -m "Release ${RELEASE_VERSION}"

## Push the commit and tag to origin
echo "Pushing release version ${RELEASE_VERSION} to origin"
git push --follow-tags

echo "Done."
