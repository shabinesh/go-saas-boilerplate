#!/bin/bash

rm transcription go.mod go.sum

# to be run from main branch
git rebase dev

find . -not -path './build_cookiecutter.sh' -type f -print0 | xargs -0 sed -i '' 's/transcription/{{cookiecutter.package_slug}}/g'
find . -not -path './build_cookiecutter.sh' -type f -print0 | xargs -0 sed -i '' 's/shabinesh/{{cookiecutter.gh_username}}/g'
