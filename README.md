# takemeon
nxdomain subdomain enumeration. Helps in scaling the automation. Currently, it only helps to resolve the `nxdomain` if possible. 

## Usage:

Typically, it's better to use it with `-mdns` flag since that would make this tool run a little faster.
```
$ cat test.txt | takemeon -mdns 8.8.8.8
test.milindpurswani.com | totallynonexistingdomain.com
test3.milindpurswani.com | totallynonexistingdomain.com
```

```
$ takemeon -h 
Usage of ./main:
  -c int
        set the concurrency level (default 1)
  -mdns string
        Specify dns server IP address. (Makes this tool run a little faster) (default "/etc/resolv.conf")
```

## Installation

```
go get -u github.com/milindpurswani/takemeon
```

