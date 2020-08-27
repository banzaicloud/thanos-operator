#!/bin/bash

set -euf

RELEASE_TAG="$1"

function main()
{
    curl \
        -X POST \
        -H "Accept: application/vnd.github.everest-preview+json" \
        -H "Content-Type: application/json" \
        -H "Authorization: token ${GITHUB_TOKEN}" \
        -d "{
            \"ref\": \"docs-generator-trigger\",
            \"inputs\": {
                \"release-tag\": \"${RELEASE_TAG}\"
            }
        }" "https://api.github.com/repos/banzaicloud/thanos-operator-docs/actions/workflows/update-generated-docs.yml/dispatches"
}

main "$@"
