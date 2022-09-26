package useragent

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os/exec"
	"strings"
)

func UserAgentMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		checkUserAgent(w, r)
		next.ServeHTTP(w, r)
		return
	})
}

func checkUserAgent(w http.ResponseWriter, r *http.Request) {
	ua := r.Header.Get("User-Agent")

	ualow := strings.ToLower(ua)

	logStr := r.URL.Path
	if strings.Contains(ualow, "google") {
		logStr += " Visited by Google bot"
	} else {
		logStr += " Visited by some thing else"
	}

	// get IP address
	ip, _, _ := net.SplitHostPort(r.RemoteAddr)

	// capture the output of host command to verify Google robots
	// based on https://support.google.com/webmasters/answer/80553

	cmd := exec.Command("host", ip)
	result, err := cmd.Output() // capture the exec output to variable result
	if err != nil {
		//fmt.Println(err)
		logStr += fmt.Sprintf(" But Host %s command execution failed.", ip)
		log.Println(logStr)
		return
	}

	// if result contain the word google, then it is genuine user agent
	// else fake

	if strings.Contains(strings.ToLower(string(result)), "google") {
		logStr += " and the user agent is real. "
	} else {
		logStr += " and the user agent is determine to be FAKED after verifying with host command. "
	}

	log.Println(logStr)

	return
}
