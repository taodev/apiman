#!/bin/zsh

file=$1

if grep -q "cases:" "$file"; then
# if [ `grep -c "url:" $file` -ne '0' ];then
    apiman --config=$file --work=./example run
else
    apiman --config=$file --work=./example api
fi