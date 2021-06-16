package ursgo

import (
	"regexp"
	"sort"
	"strings"
)

type option struct {
	exact          bool
	strict         bool
	auth           bool
	localhost      bool
	parens         bool
	apostrophes    bool
	trailingPeriod bool
	ipv4           bool
	ipv6           bool
	tlds           []string
}

type Option func(*option)

func Exact(v bool) Option {
	return func(u *option) {
		u.exact = v
	}
}

func Strict(v bool) Option {
	return func(u *option) {
		u.strict = v
	}
}

func Auth(v bool) Option {
	return func(u *option) {
		u.auth = v
	}
}

func Localhost(v bool) Option {
	return func(u *option) {
		u.localhost = v
	}
}

func Parens(v bool) Option {
	return func(u *option) {
		u.parens = v
	}
}

func Apostrophes(v bool) Option {
	return func(u *option) {
		u.apostrophes = v
	}
}

func TrailingPeriod(v bool) Option {
	return func(u *option) {
		u.trailingPeriod = v
	}
}

func IPv4(v bool) Option {
	return func(u *option) {
		u.ipv4 = v
	}
}

func IPv6(v bool) Option {
	return func(u *option) {
		u.ipv6 = v
	}
}

func Tlds(v []string) Option {
	return func(u *option) {
		u.tlds = v
	}
}

var (
	v4    = "(?:25[0-5]|2[0-4]\\d|1\\d\\d|[1-9]\\d|\\d)(?:\\.(?:25[0-5]|2[0-4]\\d|1\\d\\d|[1-9]\\d|\\d)){3}"
	v6seg = "[a-fA-F\\d]{1,4}"
)

func New(opts ...Option) (*regexp.Regexp, error) {
	opt := &option{
		localhost: true,
		ipv4:      true,
		ipv6:      true,
		tlds:      tlds,
	}
	for _, f := range opts {
		f(opt)
	}

	protocol := "(?:(?:[a-z]+:)?//)"
	if !opt.strict {
		protocol += "?"
	}

	var auth string
	if opt.auth {
		auth = "(?:\\S+(?::\\S*)?@)?"
	}

	host := "(?:(?:[a-z\\u00a1-\\uffff0-9][-_]*)*[a-z\\u00a1-\\uffff0-9]+)"
	domain := "(?:\\.(?:[a-z\\u00a1-\\uffff0-9]-*)*[a-z\\u00a1-\\uffff0-9]+)*"

	tld := "(?:\\."
	if opt.strict {
		tld += "(?:[a-z\\u00a1-\\uffff]{2,})"
	} else {
		sort.SliceStable(opt.tlds, func(i, j int) bool {
			a, b := len([]rune(opt.tlds[i])), len([]rune(opt.tlds[j]))
			return b < a // desc
		})
		tld += "(?:" + strings.Join(opt.tlds, "|") + ")"
	}
	tld += ")"
	if opt.trailingPeriod {
		tld += "\\.?"
	}

	port := "(?::\\d{2,5})?"

	var path string
	if opt.parens {
		if opt.apostrophes {
			path = "(?:[/?#][^\\s\"]*)?"
		} else {
			path = "(?:[/?#][^\\s\"']*)?"
		}
	} else {
		if opt.apostrophes {
			path = "(?:[/?#][^\\s\"\\)]*)?"
		} else {
			path = "(?:[/?#][^\\s\"\\)']*)?"
		}
	}

	regex := "(?:" + protocol + "|www\\.)" + auth + "(?:"

	if opt.localhost {
		regex += "localhost|"
	}
	if opt.ipv4 {
		regex += v4 + "|"
	}
	if opt.ipv6 {
		reg := strings.Join([]string{
			"(?:",
			"(?:${v6seg}:){7}(?:${v6seg}|:)|",        // 1:2:3:4:5:6:7::  1:2:3:4:5:6:7:8
			"(?:${v6seg}:){6}(?:${v4}|:${v6seg}|:)|", // 1:2:3:4:5:6::    1:2:3:4:5:6::8   1:2:3:4:5:6::8  1:2:3:4:5:6::1.2.3.4
			"(?:${v6seg}:){5}(?::${v4}|(?::${v6seg}){1,2}|:)|",                   // 1:2:3:4:5::      1:2:3:4:5::7:8   1:2:3:4:5::8    1:2:3:4:5::7:1.2.3.4
			"(?:${v6seg}:){4}(?:(?::${v6seg}){0,1}:${v4}|(?::${v6seg}){1,3}|:)|", // 1:2:3:4::        1:2:3:4::6:7:8   1:2:3:4::8      1:2:3:4::6:7:1.2.3.4
			"(?:${v6seg}:){3}(?:(?::${v6seg}){0,2}:${v4}|(?::${v6seg}){1,4}|:)|", // 1:2:3::          1:2:3::5:6:7:8   1:2:3::8        1:2:3::5:6:7:1.2.3.4
			"(?:${v6seg}:){2}(?:(?::${v6seg}){0,3}:${v4}|(?::${v6seg}){1,5}|:)|", // 1:2::            1:2::4:5:6:7:8   1:2::8          1:2::4:5:6:7:1.2.3.4
			"(?:${v6seg}:){1}(?:(?::${v6seg}){0,4}:${v4}|(?::${v6seg}){1,6}|:)|", // 1::              1::3:4:5:6:7:8   1::8            1::3:4:5:6:7:1.2.3.4
			"(?::(?:(?::${v6seg}){0,5}:${v4}|(?::${v6seg}){1,7}|:))",             // ::2:3:4:5:6:7:8  ::2:3:4:5:6:7:8  ::8             ::1.2.3.4
			")(?:%[0-9a-zA-Z]{1,})?", // %eth0            %1
		}, "")
		regex += reg + "|"
	}
	regex += host + domain + tld + ")" + port + path

	if opt.exact {
		regex = "(?i)(?:^" + regex + "$)"
	} else {
		regex = "(?i)" + regex
	}

	regex = strings.ReplaceAll(regex, "\\u00a1", "\\x{00a1}")
	regex = strings.ReplaceAll(regex, "\\uffff", "\\x{ffff}")
	regex = strings.ReplaceAll(regex, "${v4}", v4)
	regex = strings.ReplaceAll(regex, "${v6seg}", v6seg)

	return regexp.Compile(regex)
}
