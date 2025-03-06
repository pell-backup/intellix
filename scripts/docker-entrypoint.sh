#!/bin/bash
set -e

if [ ! -d "$PELLDVS_HOME/config" ]; then
	echo "Running pelldvs init to create (default) configuration for docker run."
	pelldvs init
fi

exec intellixd "$@"
