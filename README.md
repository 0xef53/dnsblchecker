DNSBL Checker
-------------
[![Build Status](https://drone.io/github.com/0xef53/dnsblchecker/status.png)](https://drone.io/github.com/0xef53/dnsblchecker/latest)

DNSBL Checker is a simple tool to check IPs in DNS-based blackhole lists.

### Getting binary

Latest version is available [here](https://drone.io/github.com/0xef53/dnsblchecker/files/dnsblchecker)

### Installing from source

    go build dnsblchecker.go

### How to use

Create a dnsbl servers list first. For example:

    # cat > dnsbl.txt <<EOF
    bb.barracudacentral.org
    block.stopspam.org
    bl.spamcop.net
    bl.spameatingmonkey.net
    cidr.bl.mcafee.com
    dnsbl.sorbs.net
    multi.surbl.org
    multi.uribl.com
    rhsbl.ahbl.org
    spam.dnsbl.sorbs.net
    zen.spamhaus.org
    EOF

And now you can run:

    # ./dnsblchecker 8.8.8.8

The help page is available on command:

    # ./dnsblchecker --help
