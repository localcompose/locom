#!/bin/sh
prefix="$(basename "`git remote get-url origin`" .git)/"
tag=$(git describe --tags --abbrev=0 --match "${prefix}*" 2>/dev/null || echo notags)
tagged_commit=$(git rev-list -n 1 "$tag" 2>/dev/null || echo "")
current_commit=$(git rev-parse HEAD)
short_commit=$(git rev-parse --short HEAD)
version="${tag#$prefix}"

pre=""
base="$version"

if echo "$tag" | grep -q '-'; then
  pre="${tag#*-}"
  base="${tag%-*}"
fi

if [ "$tagged_commit" != "$current_commit" ]; then
  if [ -n "$pre" ]; then
    base="$base-$pre.$short_commit"
  else
    base="$base-$short_commit"
  fi
else
  if [ -n "$pre" ]; then
    base="$base-$pre"
  fi
fi

if [ -n "$(git status --porcelain)" ]; then
  ts=$(date +%s)
  commit_ts=$(git log -1 --format=%ct)
  extra=$((ts - commit_ts))
  echo -n "$base+dev$extra"
else
  echo -n "$base"
fi
