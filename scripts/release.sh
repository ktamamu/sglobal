#!/bin/bash

# Release script for sglobal
set -e

VERSION=${1:-"0.1.0"}
REPO="ktamamu/sglobal"

echo "Creating release for version v$VERSION"

# Create git tag
git tag -a "v$VERSION" -m "Release v$VERSION"
git push origin "v$VERSION"

echo "Tag v$VERSION created and pushed. GitHub Actions will build and release automatically."
