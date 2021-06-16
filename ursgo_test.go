package ursgo_test

import (
	"regexp"
	"strings"
	"testing"

	"github.com/mix3/ursgo"
	"github.com/stretchr/testify/assert"
)

func New(opts ...ursgo.Option) *regexp.Regexp {
	v, err := ursgo.New(opts...)
	if err != nil {
		panic(err)
	}
	return v
}

func TestUrlRegexSafe(t *testing.T) {
	t.Run("match exact URLs", func(t *testing.T) {
		fixtures := []string{
			"http://-.~_!$&'()*+';=:%40:80%2f::::::@example.com",
			"//223.255.255.254",
			"//a.b-c.de",
			"//foo.ws",
			"//localhost:8080",
			"//userid:password@example.com",
			"//➡.ws/䨹",
			"ftp://foo.bar/baz",
			"http://1337.net",
			"http://142.42.1.1/",
			"http://142.42.1.1:8080/",
			"http://223.255.255.254",
			"http://a.b-c.de",
			"http://a.b.c.d.e.f.g.h.i.j.k.l.m.n.o.p.q.r.s.t.u.v.w.x.y.z.com",
			"http://a_b.z.com",
			"http://code.google.com/events/#&product=browser",
			"http://example.com#foo",
			"http://example.com.",
			"http://example.com?foo=bar",
			"http://foo.bar/?q=Test%20URL-encoded%20stuff",
			"http://foo.com/(something)?after=parens",
			"http://foo.com/blah_(wikipedia)#cite-1",
			"http://foo.com/blah_(wikipedia)_blah#cite-1",
			"http://foo.com/blah_blah",
			"http://foo.com/blah_blah/",
			"http://foo.com/blah_blah_(wikipedia)",
			"http://foo.com/blah_blah_(wikipedia)_(again)",
			"http://foo.com/unicode_(✪)_in_parens",
			"http://j.mp",
			"http://localhost/",
			"http://mw1.google.com/mw-earth-vectordb/kml-samples/gp/seattle/gigapxl/$[level]/r$[y]_c$[x].jpg",
			"http://user:pass@example.com:123/one/two.three?q1=a1&q2=a2#body",
			"http://userid:password@example.com",
			"http://userid:password@example.com/",
			"http://userid:password@example.com:8080",
			"http://userid:password@example.com:8080/",
			"http://userid@example.com",
			"http://userid@example.com/",
			"http://userid@example.com:8080",
			"http://userid@example.com:8080/",
			"http://www.example.com/wpstyle/?p=364",
			"http://www.microsoft.xn--comindex-g03d.html.irongeek.com",
			"http://⌘.ws",
			"http://⌘.ws/",
			"http://☺.damowmow.com/",
			"http://✪df.ws/123",
			"http://➡.ws/䨹",
			"https://www.example.com/foo/?bar=baz&inga=42&quux",
			"ws://223.255.255.254",
			"ws://a.b-c.de",
			"ws://foo.ws",
			"ws://localhost:8080",
			"ws://userid:password@example.com",
			"ws://➡.ws/䨹",
			"www.google.com/unicorn",
		}

		urs := New(
			ursgo.Exact(true),
			ursgo.Auth(true),
			ursgo.Parens(true),
			ursgo.TrailingPeriod(true),
		)
		for _, f := range fixtures {
			assert.True(t, urs.MatchString(f))
		}
	})

	t.Run("match exact URLs with strict set to true", func(t *testing.T) {
		fixtures := []string{
			"http://مثال.إختبار",
			"http://उदाहरण.परीक्षा",
			"http://例子.测试",
		}
		urs := New(
			ursgo.Exact(true),
			ursgo.Strict(true),
			ursgo.Auth(true),
			ursgo.Parens(true),
		)
		for _, f := range fixtures {
			assert.True(t, urs.MatchString(f))
		}
	})

	t.Run("match URLs in text", func(t *testing.T) {
		fixture := `Foo //bar.net/?q=Query with spaces
Lorem ipsum //dolor.sit
<a href="http://example.com">example.com</a>
<a href="http://example.com/with-path">with path</a>
[and another](https://another.example.com) and`

		got := New(ursgo.Strict(true)).FindAllString(fixture, -1)
		want := []string{
			"//bar.net/?q=Query",
			"//dolor.sit",
			"http://example.com",
			"http://example.com/with-path",
			"https://another.example.com",
		}
		assert.Equal(t, want, got)
	})

	t.Run("do not match URLs", func(t *testing.T) {
		fixtures := []string{
			"http://",
			"http://.",
			"http://..",
			"http://../",
			"http://?",
			"http://??",
			"http://??/",
			"http://#",
			"http://##",
			"http://##/",
			"http://foo.bar?q=Spaces should be encoded",
			"//",
			"//a",
			"///a",
			"///",
			"http:///a",
			"rdar://1234",
			"h://test",
			"http:// shouldfail.com",
			":// should fail",
			"http://foo.bar/foo(bar)baz quux",
			"http://-error-.invalid/",
			"http://-a.b.co",
			"http://a.b-.co",
			"http://123.123.123",
			"http://3628126748",
			"http://.www.foo.bar/",
			"http://.www.foo.bar./",
			"http://go/ogle.com",
			"http://foo.bar/ /",
			"http://a.b_z.com",
			"http://ab_.z.com",
			"http://google\\.com",
			"http://www(google.com",
			"http://www.example.xn--overly-long-punycode-test-string-test-tests-123-test-test123/",
			"http://www=google.com",
			"https://www.g.com/error\n/bleh/bleh",
			"/foo.bar/",
			"///www.foo.bar./",
		}
		urs := New(
			ursgo.Exact(true),
		)
		for _, f := range fixtures {
			assert.False(t, urs.MatchString(f))
		}
	})

	t.Run("do not match URLs: foo.com", func(t *testing.T) {
		assert.False(t, New(
			ursgo.Exact(true),
			ursgo.Strict(true),
		).MatchString("foo.com"))
	})

	t.Run("match using list of TLDs", func(t *testing.T) {
		fixtures := []string{
			"-.~_!$&'()*+';=:%40:80%2f::::::@example.com",
			"//223.255.255.254",
			"//a.b-c.de",
			"//foo.ws",
			"//localhost:8080",
			"//userid:password@example.com",
			"//➡.ws/䨹",
			"1337.net",
			"142.42.1.1/",
			"142.42.1.1:8080/",
			"223.255.255.254",
			"a.b-c.de",
			"a.b.c.d.e.f.g.h.i.j.k.l.m.n.o.p.q.r.s.t.u.v.w.x.y.z.com",
			"code.google.com/events/#&product=browser",
			"example.com#foo",
			"example.com",
			"example.com.",
			"example.com?foo=bar",
			"foo.bar/?q=Test%20URL-encoded%20stuff",
			"foo.bar/baz",
			"foo.com/(something)?after=parens",
			"foo.com/blah_(wikipedia)#cite-1",
			"foo.com/blah_(wikipedia)_blah#cite-1",
			"foo.com/blah_blah",
			"foo.com/blah_blah/",
			"foo.com/blah_blah_(wikipedia)",
			"foo.com/blah_blah_(wikipedia)_(again)",
			"foo.com/unicode_(✪)_in_parens",
			"foo.ws",
			"google.com",
			"j.mp",
			"localhost/",
			"localhost:8080",
			"mw1.google.com/mw-earth-vectordb/kml-samples/gp/seattle/gigapxl/$[level]/r$[y]_c$[x].jpg",
			"user:pass@example.com:123/one/two.three?q1=a1&q2=a2#body",
			"userid:password@example.com",
			"userid:password@example.com/",
			"userid:password@example.com:8080",
			"userid:password@example.com:8080/",
			"userid@example.com",
			"userid@example.com/",
			"userid@example.com:8080",
			"userid@example.com:8080/",
			"www.example.com/foo/?bar=baz&inga=42&quux",
			"www.example.com/wpstyle/?p=364",
			"www.google.com/unicorn",
			"www.microsoft.xn--comindex-g03d.html.irongeek.com",
			"⌘.ws",
			"⌘.ws/",
			"☺.damowmow.com/",
			"✪df.ws/123",
			"➡.ws/䨹",
		}

		urs := New(
			ursgo.Exact(true),
			ursgo.Auth(true),
			ursgo.Parens(true),
			ursgo.TrailingPeriod(true),
		)
		for _, f := range fixtures {
			assert.True(t, urs.MatchString(f))
		}
	})

	t.Run("opt out of matching basic auth", func(t *testing.T) {
		fixtures := []string{
			"http://-.~_!$&'()*+';=:%40:80%2f::::::@example.com",
			"http://user:pass@example.com:123/one/two.three?q1=a1&q2=a2#body",
			"http://userid:password@example.com",
			"http://userid:password@example.com/with/path",
			"http://userid:password@example.com:8080",
			"http://userid:password@example.com:8080/path",
			"http://userid@example.com",
			"http://userid@example.com/with/path",
			"http://userid@localhost:8080",
			"http://userid@localhost:8080/path",
		}

		urs1 := New(
			ursgo.Exact(true),
			ursgo.Strict(true),
			ursgo.Auth(false),
		)
		urs2 := New(
			ursgo.Exact(true),
			ursgo.Strict(true),
			ursgo.Auth(false),
		)
		for _, f := range fixtures {
			assert.False(t, urs1.MatchString(f))
			assert.False(t, urs2.MatchString(strings.Replace(f, "http", "", 1)))
			assert.False(t, urs2.MatchString(strings.Replace(f, "http://", "", 1)))
		}

		fixture := `Lorem ipsum http://userid:password@example.com:8080 dolor sit
<a href="http://userid:password@example.com:8080/">example.com</a>
another //userid:password@example.com one
bites //userid:password@example.com/with/path the dust
and http://user:pass@example.com:123/one/two.three?q1=a1&q2=a2#body another one
and <a href="http://user:pass@example.com:123/one/two.three?q1=a1&q2=a2#body">another one</a>
and another <a href="userid:password@example.com">one gone</a>
and another userid@example.com one gone
another http://userid@example.com/ one
bites http://userid@localhost:8080 the
dust http://userid@localhost:8080/path`

		assert.Nil(t, New(
			ursgo.Exact(false),
			ursgo.Strict(true),
			ursgo.Auth(false),
		).FindAllString(fixture, -1))

		want := []string{
			"example.com:8080",
			"example.com:8080/",
			"example.com",
			"example.com",
			"example.com/with/path",
			"example.com:123/one/two.three?q1=a1&q2=a2#body",
			"example.com:123/one/two.three?q1=a1&q2=a2#body",
			"example.com",
			"example.com",
			"example.com/",
			"localhost:8080",
			"localhost:8080/path",
		}

		assert.Equal(t, want, New(
			ursgo.Exact(false),
			ursgo.Strict(false),
		).FindAllString(fixture, -1))

		assert.Equal(t, want, New(
			ursgo.Exact(false),
			ursgo.Strict(false),
		).FindAllString(strings.ReplaceAll(fixture, "http:", ""), -1))

		assert.Equal(t, want, New(
			ursgo.Exact(false),
			ursgo.Strict(false),
		).FindAllString(strings.ReplaceAll(fixture, "http://", ""), -1))
	})

	t.Run("match using explicit list of TLDs", func(t *testing.T) {
		fixtures := []string{
			"-.~_!$&'()*+';=:%40:80%2f::::::@example.com",
			"-.~_!$&'()*+';=:%40:80%2f::::::@example.onion",
			"//223.255.255.254",
			"//a.b-c.de",
			"//foo.ws",
			"//localhost:8080",
			"//userid:password@example.com",
			"//➡.onion/䨹",
			"//➡.ws/䨹",
			"1337.net",
			"142.42.1.1/",
			"142.42.1.1:8080/",
			"223.255.255.254",
			"a.b-c.de",
			"a.b.c.d.e.f.g.h.i.j.k.l.m.n.o.p.q.r.s.t.u.v.w.x.y.z.com",
			"code.google.com/events/#&product=browser",
			"example.com#foo",
			"example.com.",
			"example.com?foo=bar",
			"example.onion",
			"foo.bar/?q=Test%20URL-encoded%20stuff",
			"foo.bar/baz",
			"foo.com/(something)?after=parens",
			"foo.com/blah_(wikipedia)#cite-1",
			"foo.com/blah_(wikipedia)_blah#cite-1",
			"foo.com/blah_blah",
			"foo.com/blah_blah/",
			"foo.com/blah_blah_(wikipedia)",
			"foo.com/blah_blah_(wikipedia)_(again)",
			"foo.com/unicode_(✪)_in_parens",
			"foo.ws",
			"j.mp",
			"localhost/",
			"localhost:8080",
			"mw1.google.com/mw-earth-vectordb/kml-samples/gp/seattle/gigapxl/$[level]/r$[y]_c$[x].jpg",
			"mw1.unicorn.education/mw-earth-vectordb/kml-samples/gp/seattle/gigapxl/$[level]/r$[y]_c$[x].jpg",
			"unicorn.education",
			"user:pass@example.com:123/one/two.three?q1=a1&q2=a2#body",
			"userid:password@example.com",
			"userid:password@example.com/",
			"userid:password@example.com:8080",
			"userid:password@example.com:8080/",
			"userid:password@example.education",
			"userid@example.com",
			"userid@example.com/",
			"userid@example.com:8080",
			"userid@example.com:8080/",
			"www.example.com/foo/?bar=baz&inga=42&quux",
			"www.example.com/wpstyle/?p=364",
			"www.example.onion/wpstyle/?p=364",
			"www.google.com/unicorn",
			"www.microsoft.xn--comindex-g03d.html.irongeek.com",
			"⌘.ws",
			"⌘.ws/",
			"☺.damowmow.com/",
			"✪df.ws/123",
			"➡.ws/䨹",
		}
		urs := New(
			ursgo.Exact(true),
			ursgo.Auth(true),
			ursgo.Parens(true),
			ursgo.Tlds([]string{"com", "ws", "de", "net", "mp", "bar", "onion", "education"}),
			ursgo.TrailingPeriod(true),
		)
		for _, f := range fixtures {
			assert.True(t, urs.MatchString(f))
		}
	})

	t.Run("fail if not in explicit list of TLDs", func(t *testing.T) {
		fixtures := []string{
			"-.~_!$&'()*+';=:%40:80%2f::::::@example.biz",
			"//a.b-c.uk",
			"//foo.uk",
			"//userid:password@example.biz",
			"//➡.cn/䨹",
			"1337.biz",
			"a.b-c.cn",
			"a.b-c.ly",
			"a.b.c.d.e.f.g.h.i.j.k.l.m.n.o.p.q.r.s.t.u.v.w.x.y.z.biz",
			"code.google.biz/events/#&product=browser",
			"example.biz#foo",
			"example.biz.",
			"example.biz?foo=bar",
			"foo.baz/?q=Test%20URL-encoded%20stuff",
			"foo.baz/baz",
			"foo.baz/blah_blah",
			"foo.biz/(something)?after=parens",
			"foo.biz/blah_(wikipedia)#cite-1",
			"foo.biz/blah_(wikipedia)_blah#cite-1",
			"foo.biz/blah_blah_(wikipedia)",
			"foo.biz/unicode_(✪)_in_parens",
			"foo.co.uk/blah_blah/",
			"foo.jp",
			"foo.onion/blah_blah_(wikipedia)_(again)",
			"j.onion",
			"mw1.google.biz/mw-earth-vectordb/kml-samples/gp/seattle/gigapxl/$[level]/r$[y]_c$[x].jpg",
			"user:pass@example.biz:123/one/two.three?q1=a1&q2=a2#body",
			"userid:password@example.biz",
			"userid:password@example.biz/",
			"userid:password@example.biz:8080",
			"userid:password@example.biz:8080/",
			"userid@example.biz",
			"userid@example.biz/",
			"userid@example.biz:8080",
			"userid@example.biz:8080/",
			"www.example.biz/foo/?bar=baz&inga=42&quux",
			"www.example.education/wpstyle/?p=364",
			"www.google.biz/unicorn",
			"www.microsoft.xn--comindex-g03d.html.irongeek.biz",
			"⌘.onion",
			"⌘.onion/",
			"☺.damowmow.biz/",
			"✪df.onion/123",
			"➡.onion/䨹",
			"➡.uk/䨹",
		}
		urs := New(
			ursgo.Exact(true),
			ursgo.Auth(true),
			ursgo.Parens(true),
			ursgo.Tlds([]string{"com", "ws", "de", "net", "mp", "bar"}),
		)
		for _, f := range fixtures {
			assert.False(t, urs.MatchString(f))
		}
	})

	t.Run("do not match URLs with non-strict mode", func(t *testing.T) {
		assert.False(t, New(
			ursgo.Exact(true),
			ursgo.Auth(true),
			ursgo.Parens(true),
		).MatchString("018137.113.215.4074.138.129.172220.179.206.94180.213.144.175250.45.147.1364868726sgdm6nohQ"))
	})

	t.Run("IPv4", func(t *testing.T) {
		assert.True(t, New().MatchString("1.1.1.1"))
		assert.False(t, New(
			ursgo.IPv4(false),
		).MatchString("1.1.1.1"))
	})

	t.Run("IPv6", func(t *testing.T) {
		assert.True(t, New().MatchString("2606:4700:4700::1111"))
		assert.False(t, New(
			ursgo.IPv6(false),
		).MatchString("2606:4700:4700::1111"))
	})

	t.Run("parses similar to Gmail by default", func(t *testing.T) {
		want := []string{
			"bar.com", "bar.com", "foob.com", "example.com",
		}
		assert.Equal(t, want, New().FindAllString("foo@bar.com [foo]@bar.com foo bar @foob.com 'text@example.com, some text'", -1))
	})

	t.Run("apostrophes", func(t *testing.T) {
		assert.Equal(t, []string{"http://example.com/pic.jpg"}, New().FindAllString("background: url('http://example.com/pic.jpg');", -1))
		assert.Equal(t, []string{"http://example.com/pic.jpg'"}, New(
			ursgo.Apostrophes(true),
		).FindAllString("background: url('http://example.com/pic.jpg');", -1))
		assert.Equal(t, []string{"http://example.com/pic.jpg');"}, New(
			ursgo.Parens(true),
			ursgo.Apostrophes(true),
		).FindAllString("background: url('http://example.com/pic.jpg');", -1))
	})

	t.Run("localhost", func(t *testing.T) {
		assert.Equal(t, []string{"http://localhost/pic.jpg"}, New(
			ursgo.Localhost(true),
		).FindAllString("background: url('http://localhost/pic.jpg');", -1))
		assert.Equal(t, []string{"pic.jp"}, New(
			ursgo.Localhost(false),
		).FindAllString("background: url('http://localhost/pic.jpg');", -1))
	})

	t.Run("trailing period", func(t *testing.T) {
		assert.Equal(t, []string{"example.com.", "foobar.com"}, New(
			ursgo.TrailingPeriod(true),
		).FindAllString("background example.com. foobar.com", -1))
		assert.Equal(t, []string{"example.com", "foobar.com"}, New(
			ursgo.TrailingPeriod(false),
		).FindAllString("background example.com. foobar.com", -1))
	})
}
