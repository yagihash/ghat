#!/bin/bash
set -euo pipefail

INPUT_JSON="$1"
ACTION_YML="action.yml"
NEW_YML="action.yml.new"

jq -r '.components.schemas["app-permissions"].properties | to_entries | .[] |
  "  permission_" + .key + ":\n    description: \"" + .value.description + " (" + (.value.enum | join("/")) + ")\"\n    required: false"' \
  "$INPUT_JSON" > inputs_fragment.txt

awk '
  /^  permission[-_].*:/ { skip = 1; next }
  skip && /^    / { next }
  skip && /^  [^ ]/ { skip = 0 }
  skip && /^$/ { skip = 0 }

  { print }

  /repositories:/ { in_repos = 1 }
  in_repos && /required: false/ {
    printf "\n"
    system("cat inputs_fragment.txt")
    in_repos = 0
  }
' "$ACTION_YML" > "$NEW_YML"

cat -s "$NEW_YML" > "$ACTION_YML"
rm "$NEW_YML" inputs_fragment.txt

echo "Successfully synchronized $ACTION_YML"