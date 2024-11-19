# yum-updateinfo

Build an update-info package for rpm packages to give the ability to do security updates in an automated manner.

## Usage

```
$ ./yum-updateinfo -h
YUM UpdateInfo Generator, Version # (https://github.com/pschou/yum-updateinfo)
Usage:
  yum-updateinfo [options] path_to_repodata/repomd.xml
Options:
  -collection string
        Name for collection (default "update-rpms")
  -conf string
        Config file (default "packages.yml")
  -rel string
        Set RELEASE value (default "RedHat Compatible")
```
