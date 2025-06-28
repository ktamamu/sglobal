#!/bin/bash

# Release script for sglobal
set -e

VERSION=${1:-"0.1.0"}
REPO="ktamamu/sglobal"

echo "Creating release for version v$VERSION"

# Create git tag
git tag -a "v$VERSION" -m "Release v$VERSION"
git push origin "v$VERSION"

# Build binaries
make clean
make release

echo "Binaries built in dist/ directory:"
ls -la dist/

echo ""
echo "To complete the release:"
echo "1. Upload binaries to GitHub release: https://github.com/$REPO/releases/tag/v$VERSION"
echo "2. Update homebrew-tap/Formula/sglobal.rb with new version and SHA256"
echo "3. Test the formula: brew install --build-from-source homebrew-tap/Formula/sglobal.rb"
echo ""
echo "SHA256 checksums:"
for file in dist/*; do
    if [[ -f "$file" ]]; then
        echo "$(basename $file): $(shasum -a 256 $file | cut -d' ' -f1)"
    fi
done
