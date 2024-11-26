#!/bin/bash

echo "Running Rector"

composer install --no-ansi --no-interaction --no-progress --prefer-dist --ignore-platform-reqs

"$GITHUB_WORKSPACE"/vendor/bin/rector process