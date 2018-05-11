#!/bin/bash
set -e
git clone https://github.com/josephtchung/webui litwebui
cd litwebui
echo "{" >> package.new.json
echo "\"homepage\":\"/litwebui\"," >> package.new.json 
tail -n +2 package.json >> package.new.json 
rm package.json
mv package.new.json package.json
npm install
npm run build
