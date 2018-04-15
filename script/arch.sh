#!/bin/bash

set -eu

DIST_DIR=$1

ls -d $DIST_DIR/* | 
	while read -r d
	do
		dn=`dirname $d`
		bn=`basename $d`
		tar czf $d.tar.gz -C $dn $bn
	done
