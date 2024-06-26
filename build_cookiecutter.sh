#!/bin/bash

rm transcription go.mod go.sum

# to be run from main branch
git branch -D main
git switch -c main
git rebase dev


find \{\{cookiecutter.project_slug\}\}/ -not -path './test.db' -not -path './vendor/*' -type f -print0 | xargs -0 sed -i '' -e 's/transcription/{{cookiecutter.package_slug}}/g' -e 's/shabinesh/{{cookiecutter.gh_username}}/g'
#find . -not -path './build_cookiecutter.sh' -not -path './cookiecutter.json' -not -path './test.db' -not -path './.git/*' -not -path './vendor/*' -type f -print0 | xargs -0 sed -i '' 's/shabinesh/{{cookiecutter.gh_username}}/g'

sed -e 's/transcription/{{cookiecutter.package_slug}}/g' -e 's/shabinesh/{{cookiecutter.gh_username}}/g' main.go > main.go.tmp && mv main.go.tmp main.go

last_commit_message=$(git log -1 --pretty=%B dev)

git add .
git commit -m "$last_commit_message"