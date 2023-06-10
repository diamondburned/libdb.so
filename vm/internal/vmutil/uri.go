package vmutil

import "net/url"

// MakeURI formats a URI from the given scheme, opaque, and values.
func MakeURI(scheme, opaque string, values url.Values) string {
	u := url.URL{
		Scheme:   scheme,
		Opaque:   opaque,
		RawQuery: values.Encode(),
	}
	return u.String()
}

// MakeTerminalWriteURI formats a URI that, when clicked in a terminal, will
// paste the given data.
func MakeTerminalWriteURI(data string) string {
	return MakeURI("terminal", "write", url.Values{"data": {data}})
}
