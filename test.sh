#!/usr/bin/env bash

make build rice-build

# Start bam_agent
./bam_agent -no-update & PID=$!

sleep 3

wget -q --delete-after http://localhost:1111/status
CODE=$?

if [ ${CODE} != 0 ] ; then
	echo "Failed to check status endpoint! Code: ${CODE}"
fi

# Stop by pid
kill -9 ${PID}

exit ${CODE}
