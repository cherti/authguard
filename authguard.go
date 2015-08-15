package main

import (
	"flag"
	"fmt"
	auth "github.com/abbot/go-http-auth"
	"net/http"
	"net/http/httputil"
)

// addresses and protocols
var outerAddress = flag.String("web.listen-address", ":8081", "address exposed to outside")
var innerAddress = flag.String("web.proxy-to", "127.0.0.1:8080", "address to proxy to")
var innerScheme = flag.String("scheme", "http", "scheme to use for connection to target (either http or https)")

// HTTP basic auth
var useAuth = flag.Bool("auth", true, "use HTTP-Basic-Auth for outer connection")
var user = flag.String("user", "authguard", "user for HTTP basic auth outwards")
var pass = flag.String("pass", "authguard", "password for HTTP basic auth outwards")

// TLS
var useTLS = flag.Bool("tls", true, "use TLS for outer connection")
var crt = flag.String("crt", "", "path to TLS public key file for outer connection")
var key = flag.String("key", "", "path to TLS private key file for outer connection")

// return secret for basic http-auth
func Secret(puser, realm string) string {
	if *user == puser {

		magic := []byte("$magic$")
		salt := []byte("salt")
		e := auth.MD5Crypt([]byte(*pass), salt, magic)
		return string(e)
	}
	return ""
}

// modifies incoming http.request to go to target
func director(r *http.Request) {
	r.URL.Scheme = *innerScheme
	r.URL.Host = *innerAddress
}

func redirectIt(w http.ResponseWriter, r *http.Request) {
	proxy := &httputil.ReverseProxy{Director: director}
	proxy.ServeHTTP(w, r)
}

func main() {
	flag.Parse()

	fmt.Println("starting redirector from", *outerAddress, "to", *innerAddress)

	if *useAuth {
		fmt.Println("HTTP Basic Auth enabled")
		authenticator := auth.NewBasicAuthenticator("", Secret)
		http.HandleFunc("/", auth.JustCheck(authenticator, redirectIt))
	} else {
		fmt.Println("HTTP Basic Auth disabled")
		http.HandleFunc("/", redirectIt)
	}

	var err error

	if *useTLS {
		fmt.Println("TLS enabled")
		err = http.ListenAndServeTLS(*outerAddress, *crt, *key, nil)
	} else {
		fmt.Println("TLS disabled")
		err = http.ListenAndServe(*outerAddress, nil)
	}

	if err != nil {
		fmt.Println(err)
	}

}
