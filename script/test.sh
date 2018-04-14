#!/bin/bash

set -eu

result_dir=test/result
mkdir -p $result_dir

cat << EOL | ./bin/diff2xlsx write -o $result_dir/standard.xlsx
 foobar
 hogepiyo
--- abababa
+++ abababa
 fugapiyo
- fugapiyo
+ fugapiy0
EOL

cat << EOL | ./bin/diff2xlsx write -n -o $result_dir/no-attribute.xlsx
 foobar
 hogepiyo
--- abababa
+++ abababa
 fugapiyo
- fugapiyo
+ fugapiy0
EOL

git --no-pager diff | ./bin/diff2xlsx write -o $result_dir/git-diff.xlsx
