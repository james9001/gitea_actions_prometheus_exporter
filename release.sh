#!/bin/bash

set -e

if [[ "$1" =~ ^[1-9][0-9]*\.[0-9]+\.[0-9]+$ ]]; then
  echo "Valid version format confirmed: $1"
else
  echo "Invalid version format: $1"
  echo "Usage: $0 X.Y.Y"
  exit 1
fi

TAG_VERSION=$1

if [ -z "$TAG_VERSION" ]; then
    echo "Error: No tag name provided" 1>&2
    exit 1
fi

TAG_MESSAGE="release ${TAG_VERSION}"

echo "git tag -a v${TAG_VERSION} -m ${TAG_MESSAGE}"
git tag -a v${TAG_VERSION} -m ${TAG_MESSAGE}
echo "git push origin v${TAG_VERSION}"
git push origin v${TAG_VERSION}
