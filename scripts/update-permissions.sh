#!/bin/bash
set -euo pipefail

if [ "$#" -ne 1 ]; then
    echo "Usage: $0 <path-to-openapi-json>"
    exit 1
fi

INPUT_JSON="$1"
ACTION_YML="action.yml"
NEW_YML="action.yml.new"

if [ ! -f "$INPUT_JSON" ]; then
    echo "Error: File $INPUT_JSON not found."
    exit 1
fi

jq -r '.components.schemas["app-permissions"].properties | to_entries | .[] |
  "  permission-" + (.key | gsub("_"; "-")) + ":\n    description: \"" + .value.description + " (" + (.value.enum | join("/")) + ")\"\n    required: false"' \
  "$INPUT_JSON" > inputs_fragment.txt

# args用
jq -r '.components.schemas["app-permissions"].properties | keys_unsorted[] |
  "    - ${{ inputs.permission-" + (. | gsub("_"; "-")) + " }}"' \
  "$INPUT_JSON" > args_fragment.txt

# 2. awkでマージ
awk '
  { print }
  /repositories:/ { in_repos = 1 }
  in_repos && /required: false/ { system("cat inputs_fragment.txt"); in_repos = 0 }
  / - \${{ inputs.repositories }}/ { system("cat args_fragment.txt") }
' "$ACTION_YML" > "$NEW_YML"

# 後片付け
mv "$NEW_YML" "$ACTION_YML"
rm inputs_fragment.txt args_fragment.txt