#!/bin/sh

here=$(dirname $0)

for d in $(ls $here)
do
    [ -d $d ] && (cd $d && go build main.go) &
done
