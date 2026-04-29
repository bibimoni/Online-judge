#!/bin/sh

#compile checker

source="$1" # Polygon package path
filepath="$2"
dest="$3" # problem path

cd "$source"
g++ -o checker "$filepath"
cp ./checker "$dest"

exit 0