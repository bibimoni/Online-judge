#!/bin/sh

#compile interactor

source="$1" # Polygon package path
dest="$2" # problem path

cd "$source/files"
g++ -o interactor interactor.cpp
cp ./interactor "$dest"

exit 0