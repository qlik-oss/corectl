#!/bin/bash
PACKAGES=$(go list ./... | grep -v "test/")

echo "running unit tests"
for pkg in $PACKAGES; do
  go test -count=1 -race -coverprofile=$(echo $pkg | tr / -).cover $pkg
done

echo "mode: atomic" > c.out
grep -h -v "^mode:" *.cover >> c.out
rm *.cover
