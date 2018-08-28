#!/usr/bin/env bash
DIR=`dirname $0`

go build -o $DIR/qlix $DIR/../main.go

./qlix --config $DIR/project1/qli.yml reload > $DIR/output/project1-reload.txt
./qlix --config $DIR/project1/qli.yml fields > $DIR/output/project1-fields.txt
./qlix --config $DIR/project1/qli.yml field numbers > $DIR/output/project1-field-numbers.txt
./qlix --config $DIR/project2/qli.yml reload > $DIR/output/project2-reload.txt
./qlix --config $DIR/project2/qli.yml fields > $DIR/output/project2-fields.txt

git --no-pager diff --ignore-all-space --word-diff=porcelain -- $DIR/output

git diff --ignore-all-space --exit-code $DIR/output
