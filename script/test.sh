#!/bin/bash

set -eu

cat << EOL | ./bin/diff2xlsx write -o test/result/standard.xlsx
 foobar
 hogepiyo
--- abababa
+++ abababa
 fugapiyo
- fugapiyo
+ fugapiy0
EOL

cat << EOL | ./bin/diff2xlsx write -n -o test/result/no-attribute.xlsx
 foobar
 hogepiyo
--- abababa
+++ abababa
 fugapiyo
- fugapiyo
+ fugapiy0
EOL
