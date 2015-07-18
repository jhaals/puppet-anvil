#!/bin/bash

CWD=$(pwd)

rm -rf out modules
mkdir -p out modules/puppetlabs/apache modules/puppetlabs/concat modules/puppetlabs/stdlib

go build -o puppet-anvil

docker build -t puppet-anvil .


CID=$(docker run -d -v ${CWD}/modules:/modules -p 8080:8080 puppet-anvil)


wget https://forgeapi.puppetlabs.com/v3/files/puppetlabs-apache-1.5.0.tar.gz -O modules/puppetlabs/apache/puppetlabs-apache-1.5.0.tar.gz
wget https://forgeapi.puppetlabs.com/v3/files/puppetlabs-concat-1.2.3.tar.gz -O modules/puppetlabs/concat/puppetlabs-concat-1.2.3.tar.gz
wget https://forgeapi.puppetlabs.com/v3/files/puppetlabs-stdlib-4.6.0.tar.gz -O modules/puppetlabs/stdlib/puppetlabs-stdlib-4.6.0.tar.gz

puppet module install puppetlabs/apache --module_repository http://127.0.0.1:8080 --modulepath ./out

docker stop $CID
