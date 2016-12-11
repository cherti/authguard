package main

import (
	"crypto/subtle"
	"flag"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
)

var (
	// logger
	logInfo  = log.New(os.Stdout, "", 0)
	logError = log.New(os.Stdout, "ERROR: ", 0)

	// addresses and protocols
	outerAddress = flag.String("web.listen-address", ":8080", "address exposed to outside")
	innerAddress = flag.String("web.proxy-to", "127.0.0.1:8080", "address to proxy to")
	innerScheme  = flag.String("scheme", "http", "scheme to use for connection to target (either http or https)")

	// HTTP basic auth
	useAuth = flag.Bool("auth", true, "use HTTP-Basic-Auth for outer connection")
	user    = flag.String("user", "authguard", "user for HTTP basic auth outwards")
	pass    = flag.String("pass", "authguard", "password for HTTP basic auth outwards")

	// TLS
	crt = flag.String("crt", "", "path to TLS public key file for outer connection")
	key = flag.String("key", "", "path to TLS private key file for outer connection")

	// misc
	logTimestamps = flag.Bool("config.log-timestamps", false, "Log with timestamps")
)

// director modifies the incoming http.request to go to the specified innerAddress
func director(r *http.Request) {
	r.URL.Scheme = *innerScheme
	r.URL.Host = *innerAddress
}

// performRedirect redirects the incoming request to what is specified in the innerAddress-field
func performRedirect(w http.ResponseWriter, r *http.Request) {
	proxy := &httputil.ReverseProxy{Director: director}
	proxy.ServeHTTP(w, r)
}

// redirectAfterAuthCheck checks for correct authentication-credentials and either applies
// the intended redirect or asks for authentication-credentials once again.
func redirectAfterAuthCheck(w http.ResponseWriter, r *http.Request) {
	u, p, ok := r.BasicAuth()
	if ok && subtle.ConstantTimeCompare([]byte(u), []byte(*user)) == 1 && subtle.ConstantTimeCompare([]byte(p), []byte(*pass)) == 1 {
		performRedirect(w, r)
	} else {
		// send out unauthenticated response asking for basic auth
		// (to make sure people that mistyped can retry)
		w.Header().Set("WWW-Authenticate", `Basic realm="all"`)
		http.Error(w, "Unauthenticated", 401)
	}
}

func main() {
	flag.Parse()
	if *logTimestamps {
		logInfo.SetFlags(3)
		logError.SetFlags(3)
	}

	logInfo.Println("starting redirector from", *outerAddress, "to", *innerAddress)

	if *useAuth {
		logInfo.Println("HTTP Basic Auth enabled")
		http.HandleFunc("/", redirectAfterAuthCheck)
	} else {
		logInfo.Println("HTTP Basic Auth disabled")
		http.HandleFunc("/", performRedirect)
	}

	useTLS := *crt != "" && *key != ""
	if useTLS {
		logInfo.Println("TLS enabled")
		logError.Fatal(http.ListenAndServeTLS(*outerAddress, *crt, *key, nil))
	} else {
		logInfo.Println("TLS disabled")
		logError.Fatal(http.ListenAndServe(*outerAddress, nil))
	}
}
