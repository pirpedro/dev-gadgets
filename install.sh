#!/usr/bin/env bash

make_install() {
  if [ -n "$PREFIX" ]; then
    PREFIX="$PREFIX" make install
  else
    sudo make install
  fi
}

dir=$(mktemp -t -d git-gadgets-install.XXXXXXXXXX)
trap 'rm -rf "$dir"' EXIT
cd "$dir" &&
  echo "Setting up 'git-gadgets'...." &&
  git clone https://github.com/pirpedro/git-gadgets.git &>/dev/null &&
  cd git-gadgets &&
  git checkout \
    $(git describe --tags $(git rev-list --tags --max-count=1)) \
    &>/dev/null &&
  make_install
