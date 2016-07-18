#!/bin/sh

for i in `seq 1 16`; do
	DATE=`date -u +"%Y-%m-%dT%H:%M:%SZ"`
	sed 's/{{ DATE }}/'$DATE'/g' < change.json > /tmp/change
	curl -XPOST -T /tmp/change http://localhost:7778/update
done
