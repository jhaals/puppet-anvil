# go-puppet-forge

This is a minimal Go implementation of the Puppet Forge v3 API without external libraries. This project is inspired by [simple-puppet-forge](https://github.com/dalen/simple-puppet-forge).
No database is required, metadata is stored on disk.

### Installation
Pre-built binaries [here](http://mumble.ifup.se/go-puppet-forge/)

Requires GNU tar

Modules should be stored in `/var/lib/go-puppet-forge/modules`
with the following directory structure `user/module/user-module-version.tar.gz`

go-puppet-forge is listening on port 8080

example:

    /var/lib/go-puppet-forge/modules/puppetlabs/apache/puppetlabs-apache-1.1.0.tar.gz

#### Usage
A custom module_repository can be specified in the puppet config file.

    module_repository=http://my-forge.com/

Or specified directly install command

    ~ puppet module install puppetlabs/apache --module_repository http://127.0.0.1:8080 --modulepath modules
    Notice: Preparing to install into /Users/jhaals/modules ...
    Notice: Created target directory /Users/jhaals/modules
    Notice: Downloading from http://127.0.0.1:8080 ...
    Notice: Installing -- do not interrupt ...
    /Users/jhaals/modules
    └─┬ puppetlabs-apache (v1.1.0)
      ├── puppetlabs-concat (v1.1.0)
      └── puppetlabs-stdlib (v4.2.2)