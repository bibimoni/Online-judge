#!/bin/sh

# getting tests

source="$1" # Polygon package path
dest="$2" # problem path

mkdir -p "$dest/tests/input" "$dest/tests/output" || exit -1

cd "$source/tests" || exit -1

for f in *; do
  [ -f "$f" ] || continue

  name="${f%.a}" || exit -1

  # Safely convert to integer to remove leading zeros
  # (echo "$name" | grep -qE '^[0-9]+$' && name=$(printf "%d" "$name")) || exit -1
  name=$(echo "$name" | sed 's/^0*//') || exit -1
  
  # if [ "$f" = *.a ]; then 
  #   cp -- "$f" "$dest/tests/output/$name" || exit -1
  # else 
  #   cp -- "$f" "$dest/tests/input/$name" || exit -1
  # fi

  case "$f" in
    *.a)
      echo "$f ends with .a"
      cp -- "$f" "$dest/tests/output/$name" || exit -1
      ;;
    *)
      echo "$f does not end with .a"
      cp -- "$f" "$dest/tests/input/$name" || exit -1
      ;;
  esac
done

exit 0