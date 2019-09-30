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

function sanity_check() {
  if [[ ! -z $(git status --porcelain) ]]; then
    echo "There are uncommitted changes. Please make sure branch is clean."
    git status --porcelain
    exit 1
  fi
  local_branch=$(git rev-parse --abbrev-ref HEAD)
  if [[ $local_branch != "master" ]]; then
    echo "This script can only be run from the master branch."
    echo "You are on '$local_branch'. Aborting."
    exit 1
  fi
  # Check if local branch is up-to-date with remote master branch
  git fetch origin master
  git diff origin/master --exit-code > /dev/null
  if [[ $? -ne 0 ]]; then
    echo "Local branch is not up-to-date with remote master. Please pull the latest changes."
    git diff origin/master --name-only
    exit 1
  fi
}

## Verify that the local branch is pristine
sanity_check

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
