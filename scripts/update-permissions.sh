#!/bin/bash
set -euo pipefail

INPUT_JSON="$1"
ACTION_YML="action.yml"
NEW_YML="action.yml.new"

# The schema is third-party content (octokit/openapi): keys must stay in a safe
# charset and descriptions are JSON-escaped (@json) so they cannot break out of
# the YAML double-quoted scalar.
jq -r '.components.schemas["app-permissions"].properties | to_entries | .[] |
  "  permission_" + (.key | if test("^[a-zA-Z0-9_-]+$") then . else error("unexpected permission key: " + .) end)
  + ":\n    description: " + ((.value.description + " (" + (.value.enum | join("/")) + ")") | @json)
  + "\n    required: false"' \
  "$INPUT_JSON" > inputs_fragment.txt

awk '
  /^  permission[-_].*:/ { skip = 1; next }
  skip && /^    / { next }
  skip && /^  [^ ]/ { skip = 0 }
  skip && /^$/ { skip = 0 }

  { print }

  /repositories:/ { in_repos = 1 }
  in_repos && /required: false/ {
    system("cat inputs_fragment.txt")
    in_repos = 0
  }
' "$ACTION_YML" > "$NEW_YML"

cat -s "$NEW_YML" > "$ACTION_YML"
rm "$NEW_YML" inputs_fragment.txt