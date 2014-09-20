# Puppet Anvil

This is a minimal Go implementation of the Puppet Forge v3 API without external libraries. This project is inspired by [simple-puppet-forge](https://github.com/dalen/simple-puppet-forge).
No database is required, metadata is stored on disk.

### Installation

Modules must be stored in the following directory structure `user/module/user-module-version.tar.gz`
example:

    /var/lib/puppet-anvil/modules/puppetlabs/apache/puppetlabs-apache-1.1.0.tar.gz

You can create a .deb package for Ubuntu using `make deb`. fpm is required to create the package.

__Running Puppet Anvil__

    $ export MODULEPATH=/var/lib/puppet-anvil/modules
    $ export PORT=8080
    $ ./puppet-anvil
    Starting Puppet Anvil on port 8080 serving modules from /var/lib/puppet-anvil/modules

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