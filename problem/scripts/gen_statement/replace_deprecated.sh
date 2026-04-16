#!/bin/sh

# Process all .tex files in the current directory (recursively)
find . -type f -name "*.tex" | while read -r file; do
  echo "Fixing $file"

  # Replace deprecated font commands with proper \text... equivalents
  sed -i.bak -E \
    -e 's/\\bf\s*\{([^}]*)\}/\\textbf{\1}/g' \
    -e 's/\\it\s*\{([^}]*)\}/\\textit{\1}/g' \
    -e 's/\\rm\s*\{([^}]*)\}/\\textrm{\1}/g' \
    -e 's/\\sf\s*\{([^}]*)\}/\\textsf{\1}/g' \
    -e 's/\\tt\s*\{([^}]*)\}/\\texttt{\1}/g' \
    "$file"
done