#!/usr/bin/env bash
set -euo pipefail

for src in resources/raw/*.json; do
  base=$(basename "$src" .json)
  python3 -c "
import json, sys
data = json.load(open('$src', encoding='utf-8'))
lines = [entry['Contents'].strip() for entry in data if entry['Contents'].strip()]
print('\n'.join(lines))
" > "resources/raw/$base.txt"
  echo "wrote resources/$base.txt"
done
