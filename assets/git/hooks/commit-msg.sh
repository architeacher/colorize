#!/bin/bash
#
# An example hook script to check the commit log message.
# Called by "git commit" with one argument, the name of the file
# that has the commit message.  The hook should exit with non-zero
# status after issuing an appropriate message if it wants to stop the
# commit.  The hook is allowed to edit the commit message file.
#
# To enable this hook, rename this file to "commit-msg".

# Uncomment the below to add a Signed-off-by line to the message.
# Doing this in a hook is a bad idea in general, but the prepare-commit-msg
# hook is more suited to it.
#
# SOB=$(git var GIT_AUTHOR_IDENT | sed -n 's/^\(.*>\).*$/Signed-off-by: \1/p')
# grep -qs "^$SOB" "$1" || echo "$SOB" >> "$1"

# This example catches duplicate Signed-off-by lines.

test "" = "$(grep '^Signed-off-by: ' "$1" |
	 sort | uniq -c | sed -e '/^[ 	]*1[ 	]/d')" || {
	echo >&2 Duplicate Signed-off-by lines.
	exit 1
}

if grep -q -i -e "WIP" -e "work in progress" "$1"; then
    read -n 1 -r -p "You're about to add a WIP commit, do you want to run the CI? [y|n] " < /dev/tty
    echo
    if echo "$REPLY" | grep -E '^[Nn]$' > /dev/null; then
        echo "[Skipping CI]" >> "$1"
    fi
fi

if ! grep -qE "^:[a-z0-9_-]+:" "$1"; then
    read -n 1 -r -p "You're about to commit without an icon, do you want to continue? [y|n] " < /dev/tty
    echo
    if echo "$REPLY" | grep -E '^[Nn]$' > /dev/null; then
        echo "Skipping commit message's icon check."
        exit 1
    fi
fi

if ! grep -qE "^:[a-z0-9_-]+:\s+.*" "$1"; then
    read -n 1 -r -p "You're about to commit without a commit message, do you want to continue? [y|n] " < /dev/tty
    echo
    if echo "$REPLY" | grep -E '^[Nn]$' > /dev/null; then
        echo "Skipping commit message's check."
        exit 1
    fi
fi

if ! grep -iqE "\.$" "$1"; then
    read -n 1 -r -p "You're about to commit without ending with full stop, do you want to continue? [y|n] " < /dev/tty
    echo
    if echo "$REPLY" | grep -E '^[Nn]$' > /dev/null; then
        echo "Skipping commit message's format check."
        exit 1
    fi
fi

if ! grep -qE ".*\[(major|minor|patch)\].*" "$1"; then
    read -n 1 -r -p "You're about to commit without release type [major|minor|patch], do you want to continue? [y|n] " < /dev/tty
    echo
    if echo "$REPLY" | grep -E '^[Nn]$' > /dev/null; then
        echo "Skipping commit message's release check."
        exit 1
    fi
fi

if ! grep -qE ".*Signed-off-by:\s.*" "$1"; then
    read -n 1 -r -p "You're about to commit a non-signed commit, do you want to continue? [y|n] " < /dev/tty
    echo
    if echo "$REPLY" | grep -E '^[Nn]$' > /dev/null; then
        echo "Skipping commit message's signage check."
        exit 1
    fi
fi
