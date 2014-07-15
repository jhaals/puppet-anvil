# go-puppet-forge

This is a minimal Go implementation of the Puppet Forge v3 API without external libraries. This project is inspired by [simple-puppet-forge](https://github.com/dalen/simple-puppet-forge).
No database is required, metadata is stored on disk.

### Installation
Pre-built binaries [here](http://dl.bintray.com/jhaals/generic/go-puppet-forge/)

Requires GNU tar

Port and module directory is configured via environment variables.

Modules must be stored in the following directory structure `user/module/user-module-version.tar.gz`
example:

    /var/lib/go-puppet-forge/modules/puppetlabs/apache/puppetlabs-apache-1.1.0.tar.gz

__Running go-puppet-forge__

    $ export MODULEPATH=/var/lib/go-puppet-forge/modules
    $ export PORT=8080
    $ ./go-puppet-forge
    Starting go-puppet-forge on port 8080 serving modules from /var/lib/go-puppet-forge/modules

#### Usage with Puppet
A custom module_repository can be specified in the puppet config file.

    module_repository=http://my-forge.com/

Or directly on command line

    ~ puppet module install puppetlabs/apache --module_repository http://127.0.0.1:8080 --modulepath modules
    Notice: Preparing to install into /Users/jhaals/modules ...
    Notice: Created target directory /Users/jhaals/modules
    Notice: Downloading from http://127.0.0.1:8080 ...
    Notice: Installing -- do not interrupt ...
    /Users/jhaals/modules
    └─┬ puppetlabs-apache (v1.1.0)
      ├── puppetlabs-concat (v1.1.0)
      └── puppetlabs-stdlib (v4.2.2)