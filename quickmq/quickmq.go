package quickmq

import (
	"fmt"
	"net"
	"net/url"
)

// quickUrl should be of the format quick://username:password@host:port
func parseUrl(quickUrl string) (newUrl, username, password string, err error) {
	u, err := url.Parse(quickUrl)
	if err != nil {
		return "", "", "", err
	}

	host, port, err := net.SplitHostPort(u.Host)
	if err != nil {
		return "", "", "", err
	}

	username = u.User.Username()
	if username == "" {
		username = "admin"
	}

	password, ok := u.User.Password()
	if !ok {
		password = "admin"
	}

	newUrl = fmt.Sprintf("http://%s:%s", host, port)
	return newUrl, username, password, nil
}
