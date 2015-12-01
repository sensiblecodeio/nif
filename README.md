# nif - Simple network interface info tool

Stop all that `ip`, `ifconfig`, `/dev/net`, `grep`, `cut`, `sed`, `awk` crap!
Use a single binary to do all that for you, reliably, cross-platform.

Sometimes you just need a list of network interface names, or just a single
best-guessed network interface name, or just its best-guessed internal IP
address. This tool helps you with that.


## Examples

### List all sensible network interfaces

    $ nif
    en0
    en1
    en2
    p2p0
    awdl0
    vboxnet0
    vboxnet1
    vboxnet2
    vmnet1
    vmnet8
    
### List the first sensible network interface

    $ nif -1
    en0

### List only the IPv4 of the first sensible network interface

    $ nif -1 -i -4
    192.168.0.6

## Usage

    $ nif --help
    NAME:
       nif - Simple network interface info tool
    
    USAGE:
       nif [global options] command [command options] [arguments...]
    
    VERSION:
       2.0
    
    COMMANDS:
       help, h      Shows a list of commands or help for one command
    
    GLOBAL OPTIONS:
    --all, -a               List all available network interfaces
    --one, -o, -1           Show only single best guessed network interfaces and/or IP address
    --ipv4, -4              Show IPv4 addresses next to network interface
    --ipv6, -6              Show IPv6 addresses next to network interface
    --only-ip, -i           Only show IP addresses of network interface
    --retry, -r "0"         Retry n times in intervals of 1sec if no interface addresses could be found
    --debug, -d             Show additional debug information
    --help, -h              show help
    --version, -v           print the version
    

## Installation

    $ go get github.com/scraperwiki/nif


## Build

    $ go build
