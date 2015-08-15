# AuthGuard

AuthGuard is a tool that allows redirecting an HTTP-request and add HTTP Basic Auth as well as TLS into the chain.

This could, for example, be used to add Basic Auth and TLS to a webservice that doesn't provide that.
Just firewall that service for everything but localhost.
Then fire up AuthGuard next to that service and use it as a proxy to your service with enabled Authentication and TLS.

One concrete example is the [Prometheus](www.prometheus.io)-monitoring-system, which doesn't provide authentication out of the box.


## Building and running

### manually

    # get dependencies
    go get -u auth "github.com/abbot/go-http-auth"
    
    # actually build and run
    git clone https://github.com/cherti/authguard.git
    go build authguard.go
    ./authguard -help


### automatically using go-toolchain

    go get -u "github.com/cherti/authguard"
    ./authguard -help


## Configuration

Configuration is done soley via commandline options:

    Usage of ./authguard:
      -web.listen-address=":8081": address exposed to outside
      -web.proxy-to="127.0.0.1:8080": address to proxy to
      -scheme="http": scheme to use for connection to target (either http or https)
    
      -auth=true: use HTTP-Basic-Auth for outer connection
      -user="authguard": user for HTTP basic auth outwards
      -pass="authguard": password for HTTP basic auth outwards
    
      -tls=true: use TLS for outer connection
      -crt="": path to TLS public key file for outer connection
      -key="": path to TLS private key file for outer connection


## License

This works is released under the [GNU General Public License v3](https://www.gnu.org/licenses/gpl-3.0.txt). You can find a copy of this license at https://www.gnu.org/licenses/gpl-3.0.txt.
