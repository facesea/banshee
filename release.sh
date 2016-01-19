#!/bin/bash
echo "====== pull latest upstream/master ========"
git pull upstream master
git rev-parse upstream/master > commit
echo "====== install static deps ========="
cd static && npm install -q && cd ..
cd static/public && npm install -q && cd ../..
echo "====== build static files ========"
cd static && gulp build && cd ..
echo "====== build binary ========"
godep go build
echo "====== create package ========"
version=$(./banshee -v)
os=$(uname | awk '{print tolower($0)}')
arch=$(go env GOARCH)
printf "====== banshee version: %s\n" ${version}
printf "====== build env: %s %s" ${os} ${arch}
dir=$(printf "banshee%s.%s-%s" ${version} ${os} ${arch})
mkdir -p ${dir}
mkdir -p ${dir}/static
mv commit ${dir}/
cp ./banshee ${dir}/
cp -r ./static/dist ${dir}/static || true
cp ./LICENSE ${dir}/
cp ./README.md ${dir}/
pkg=$(printf "%s.tar.gz" ${dir})
tar cvzf ${pkg} ${dir}
rm -rf ${dir}
du -h ${pkg}
