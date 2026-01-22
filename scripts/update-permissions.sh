#!/bin/bash
set -euo pipefail

INPUT_JSON="$1"
ACTION_YML="action.yml"
NEW_YML="action.yml.new"

jq -r '.components.schemas["app-permissions"].properties | to_entries | .[] |
  "  permission_" + (.key) + ":\n    description: \"" + .value.description + " (" + (.value.enum | join("/")) + ")\"\n    required: false"' \
  "$INPUT_JSON" > inputs_fragment.txt

awk '
  { print }
  /repositories:/ { in_repos = 1 }
  in_repos && /required: false/ { system("cat inputs_fragment.txt"); in_repos = 0 }
' "$ACTION_YML" > "$NEW_YML"

mv "$NEW_YML" "$ACTION_YML"
rm inputs_fragment.txt