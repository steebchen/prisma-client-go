// Package gocase is a package to convert normal CamelCase to Golang's CamelCase and vice versa.
// Golang's CamelCase means a string that takes into account to Go's common initialisms.
// For more details, see `golang.org/x/lint/golint`.
package gocase

import "strings"

// To returns a string converted to Go case.
func To(s string) string {
	for _, ci := range commonInitialisms {
		s = strings.ReplaceAll(s, ci[1], ci[0])
	}
	return s
}

// Revert returns a string converted from Go case to normal case.
func Revert(s string) string {
	for _, ci := range commonInitialisms {
		s = strings.ReplaceAll(s, ci[0], ci[1])
	}
	return s
}

// commonInitialisms is a list of common initialisms.
// This list must match `golang.org/x/lint/golint`.
var commonInitialisms = [][]string{
	{"ACL", "Acl"},
	{"API", "Api"},
	{"ASCII", "Ascii"},
	{"CPU", "Cpu"},
	{"CSS", "Css"},
	{"DNS", "Dns"},
	{"EOF", "Eof"},
	{"GUID", "Guid"},
	{"HTML", "Html"},
	{"HTTP", "Http"},
	{"HTTPS", "Https"},
	{"ID", "Id"},
	{"IP", "Ip"},
	{"JSON", "Json"},
	{"LHS", "Lhs"},
	{"QPS", "Qps"},
	{"RAM", "Ram"},
	{"RHS", "Rhs"},
	{"RPC", "Rpc"},
	{"SLA", "Sla"},
	{"SMTP", "Smtp"},
	{"SQL", "Sql"},
	{"SSH", "Ssh"},
	{"TCP", "Tcp"},
	{"TLS", "Tls"},
	{"TTL", "Ttl"},
	{"UDP", "Udp"},
	{"UI", "Ui"},
	{"UID", "Uid"},
	{"UUID", "Uuid"},
	{"URI", "Uri"},
	{"URL", "Url"},
	{"UTF8", "Utf8"},
	{"VM", "Vm"},
	{"XML", "Xml"},
	{"XMPP", "Xmpp"},
	{"XSRF", "Xsrf"},
	{"XSS", "Xss"},
}
