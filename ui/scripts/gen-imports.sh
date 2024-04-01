#!/usr/bin/bash

echo "Generating index in src/components/index.ts"
echo "// Generated with 'scripts/gen-imports.sh'" > src/components/index.ts
find src/components -name \*.vue \
  |cut -d/ -f3- \
  |xargs -n1 bash -c 'echo export "{ default as $(echo ${0#*/}|cut -d. -f1) }" from "\"./$0"\"' \
  >> src/components/index.ts
