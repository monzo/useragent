/* Copyright (C) 2015 by Alexandru Cojocaru */

/* This program is free software: you can redistribute it and/or modify
   it under the terms of the GNU General Public License as published by
   the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.

   This program is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU General Public License for more details.

   You should have received a copy of the GNU General Public License
   along with this program.  If not, see <http://www.gnu.org/licenses/>. */

package useragent

import (
	"fmt"
	"github.com/blang/semver"
	"log"
	"testing"
)

func ExampleParse() {
	ua := Parse("Mozilla/5.0 (X11; Linux i686; rv:38.0) Gecko/20100101 Firefox/38.0")
	fmt.Print(ua)
	// Output: Type: Browser
	//Name: Firefox
	//Version: 38.0.0
	//OS: GNU/Linux
	//OSVersion: 0.0.0
	//Security: Unknown security
	//Mobile: false
	//Tablet: false
}

func ExampleParse_access() {
	ua := Parse("Mozilla/5.0 (X11; Linux i686; rv:38.0) Gecko/20100101 Firefox/38.0")
	if ua == nil {
		log.Fatal("cannot parse user agent string")
	}
	fmt.Println(ua.Type)
	fmt.Println(ua.Name)
	fmt.Println(ua.Version)
	fmt.Println(ua.OS)
	if ua.Security != SecurityUnknown {
		fmt.Println(ua.Security)
	}

	//Output:Browser
	//Firefox
	//38.0.0
	//GNU/Linux
}

func eqUA(a *UserAgent, b *UserAgent) bool {
	if a == nil || b == nil {
		return false
	}

	if a.Type != b.Type ||
		a.OS != b.OS ||
		!a.OSVersion.EQ(b.OSVersion) ||
		a.Name != b.Name ||
		!a.Version.EQ(b.Version) ||
		a.Security != b.Security ||
		a.Mobile != b.Mobile ||
		a.Tablet != b.Tablet {
		return false
	}
	return true
}

func mustParse(s string) semver.Version {
	v, err := semver.ParseTolerant(s)
	if err != nil {
		panic(`semver: Parse(` + s + `): ` + err.Error())
	}
	return v
}

func TestGecko(t *testing.T) {
	var got *UserAgent
	want := &UserAgent{}

	got = Parse(`Mozilla/5.0 (X11; U; Linux i686; rv:38.0) Gecko/20100101 Firefox/38.0`)
	want.Type = Browser
	want.OS = "GNU/Linux"
	want.OSVersion = semver.Version{}
	want.Name = "Firefox"
	want.Version = mustParse("38.0.0")
	want.Security = SecurityStrong
	if !eqUA(want, got) {
		t.Errorf("expected %+v, got %+v\n", want, got)
	}

	got = Parse(`Mozilla/5.0 (X11; U; Linux x86_64; sv-SE; rv:1.9.1.16) Gecko/20120714 IceCat/3.5.16 (like Firefox/3.5.16)`)
	want.Type = Browser
	want.OS = "GNU/Linux"
	want.OSVersion = semver.Version{}
	want.Name = "IceCat"
	want.Version = mustParse("3.5.16")
	want.Security = SecurityStrong
	if !eqUA(want, got) {
		t.Errorf("expected %+v, got %+v\n", want, got)
	}

	got = Parse(`Mozilla/5.0 (Windows x86; rv:19.0) Gecko/20100101 Firefox/19.0`)
	want.Type = Browser
	want.OS = "Windows"
	want.OSVersion = semver.Version{}
	want.Name = "Firefox"
	want.Version = mustParse("19.0.0")
	want.Security = SecurityUnknown
	if !eqUA(want, got) {
		t.Errorf("expected:\n%+v\ngot:\n%+v\n", want, got)
	}

	got = Parse(`Mozilla/5.0 (Mobile; rv:26.0) Gecko/26.0 Firefox/26.0`)
	want.Type = Browser
	want.OS = "Firefox OS"
	want.OSVersion = semver.Version{}
	want.Name = "Firefox"
	want.Version = mustParse("26.0.0")
	want.Security = SecurityUnknown
	want.Mobile = true
	if !eqUA(want, got) {
		t.Errorf("expected:\n%+v\ngot:\n%+v\n", want, got)
	}

	got = Parse(`Mozilla/5.0 (iPod touch; CPU iPhone OS 8_3 like Mac OS X) AppleWebKit/600.1.4 (KHTML, like Gecko) FxiOS/1.0 Mobile/12F69 Safari/600.1.4`)
	want.Type = Browser
	want.OS = "iOS"
	want.OSVersion = mustParse("8.3")
	want.Name = "Firefox"
	want.Version = mustParse("1.0.0")
	want.Security = SecurityUnknown
	want.Mobile = true
	if !eqUA(want, got) {
		t.Errorf("expected %+v, got %+v\n", want, got)
	}

	got = Parse(`Mozilla/5.0 (iPhone; CPU iPhone OS 8_3 like Mac OS X) AppleWebKit/600.1.4 (KHTML, like Gecko) FxiOS/1.0 Mobile/12F69 Safari/600.1.4`)
	want.Type = Browser
	want.OS = "iOS"
	want.OSVersion = mustParse("8.3")
	want.Name = "Firefox"
	want.Version = mustParse("1.0.0")
	want.Security = SecurityUnknown
	want.Mobile = true
	want.Tablet = false
	if !eqUA(want, got) {
		t.Errorf("expected %+v, got %+v\n", want, got)
	}

	got = Parse(`Mozilla/5.0 (iPad; CPU iPhone OS 8_3 like Mac OS X) AppleWebKit/600.1.4 (KHTML, like Gecko) FxiOS/1.0 Mobile/12F69 Safari/600.1.4`)
	want.Type = Browser
	want.OS = "iOS"
	want.OSVersion = mustParse("8.3")
	want.Name = "Firefox"
	want.Version = mustParse("1.0.0")
	want.Security = SecurityUnknown
	want.Mobile = false
	want.Tablet = true
	if !eqUA(want, got) {
		t.Errorf("expected %+v, got %+v\n", want, got)
	}

	// Silk on Kindle Fire: Tablet mode
	got = Parse(`Mozilla/5.0 (Linux; Android 4.4.3; KFTHWI Build/KTU84M) AppleWebKit/537.36 (KHTML, like Gecko) Silk/44.1.54 like Chrome/44.0.2403.63 Safari/537.36`)
	want.Type = Browser
	want.OS = "Android"
	want.OSVersion = mustParse("4.4.3")
	want.Name = "Silk"
	want.Version = mustParse("44.1.54")
	want.Security = SecurityUnknown
	want.Mobile = false
	want.Tablet = true
	if !eqUA(want, got) {
		t.Errorf("expected:\n%+v\ngot:\n%+v\n", want, got)
	}

	// Silk on Kindle Fire: Mobile mode
	got = Parse(`Mozilla/5.0 (Linux; U; Android 4.4.3; KFTHWI Build/KTU84M) AppleWebKit/537.36 (KHTML, like Gecko) Silk/44.1.54 like Chrome/44.0.2403.63 Mobile Safari/537.36`)
	want.Type = Browser
	want.OS = "Android"
	want.OSVersion = mustParse("4.4.3")
	want.Name = "Silk"
	want.Version = mustParse("44.1.54")
	want.Security = SecurityStrong
	want.Mobile = true
	want.Tablet = false
	if !eqUA(want, got) {
		t.Errorf("expected:\n%+v\ngot:\n%+v\n", want, got)
	}
}

func TestChrome(t *testing.T) {
	var got *UserAgent
	want := &UserAgent{}

	got = Parse(`Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/41.0.2227.0 Safari/537.36`)
	want.Type = Browser
	want.OS = "GNU/Linux"
	want.OSVersion = semver.Version{}
	want.Name = "Chrome"
	want.Version = mustParse("41.0.2227")
	want.Security = SecurityUnknown
	if !eqUA(want, got) {
		t.Errorf("expected %+v, got %+v\n", want, got)
	}

	got = Parse(`Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/41.0.2228.0 Safari/537.36`)
	want.Type = Browser
	want.OS = "Windows"
	want.OSVersion = mustParse("6.1")
	want.Name = "Chrome"
	want.Version = mustParse("41.0.2228")
	want.Security = SecurityUnknown
	if !eqUA(want, got) {
		t.Errorf("expected %+v, got %+v\n", want, got)
	}

	got = Parse(`Mozilla/5.0 (Linux; Android 4.0.4; Galaxy Nexus Build/IMM76B) AppleWebKit/535.19 (KHTML, like Gecko) Chrome/18.0.1025.133 Mobile Safari/535.19`)
	want.Type = Browser
	want.OS = "Android"
	want.OSVersion = mustParse("4.0.4")
	want.Name = "Chrome"
	want.Version = mustParse("18.0.1025")
	want.Security = SecurityUnknown
	want.Mobile = true
	if !eqUA(want, got) {
		t.Errorf("expected %+v, got %+v\n", want, got)
	}

	got = Parse(`Mozilla/5.0 (iPhone; U; CPU iPhone OS 5_1_1 like Mac OS X; en) AppleWebKit/534.46.0 (KHTML, like Gecko) CriOS/19.0.1084.60 Mobile/9B206 Safari/7534.48.3`)
	want.Type = Browser
	want.OS = "iOS"
	want.OSVersion = mustParse("5.1.1")
	want.Name = "Chrome"
	want.Version = mustParse("19.0.1084")
	want.Security = SecurityStrong
	want.Mobile = true
	if !eqUA(want, got) {
		t.Errorf("expected %+v, got %+v\n", want, got)
	}

	got = Parse(`Mozilla/5.0 (Linux; Android 6.0; Nexus 5 Build/MRA58N) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/46.0.2490.76 Safari/537.36`)
	want.Type = Browser
	want.OS = "Android"
	want.OSVersion = mustParse("6.0")
	want.Name = "Chrome"
	want.Version = mustParse("46.0.2490")
	want.Security = SecurityUnknown
	want.Mobile = false
	want.Tablet = true
	if !eqUA(want, got) {
		t.Errorf("expected %+v, got %+v\n", want, got)
	}

	got = Parse(`Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/53.0.2785.143 Safari/537.36`)
	want.Type = Browser
	want.OS = "Mac OS X"
	want.OSVersion = mustParse("10.12.0")
	want.Name = "Chrome"
	want.Version = mustParse("53.0.2785")
	want.Security = SecurityUnknown
	want.Mobile = false
	want.Tablet = false
	if !eqUA(want, got) {
		t.Errorf("expected %+v, got %+v\n", want, got)
	}
}

// Android's Chromium-based web rendering library
func TestWebView(t *testing.T) {
	var got *UserAgent
	want := &UserAgent{}

	got = Parse(`Mozilla/5.0 (Linux; Android 5.1.1; Nexus 5 Build/LMY48B; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/43.0.2357.65 Mobile Safari/537.36`)
	want.Type = Library
	want.OS = "Android"
	want.OSVersion = mustParse("5.1.1")
	want.Name = "WebView"
	want.Version = mustParse("43.0.2357")
	want.Security = SecurityUnknown
	want.Mobile = true
	want.Tablet = false
	if !eqUA(want, got) {
		t.Errorf("expected %+v, got %+v\n", want, got)
	}

	got = Parse(`Mozilla/5.0 (Linux; Android 5.0.2; SM-T350 Build/LRX22G; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/49.0.2623.105 Safari/537.36`)
	want.Type = Library
	want.OS = "Android"
	want.OSVersion = mustParse("5.0.2")
	want.Name = "WebView"
	want.Version = mustParse("49.0.2623")
	want.Security = SecurityUnknown
	want.Mobile = false
	want.Tablet = true
	if !eqUA(want, got) {
		t.Errorf("expected %+v, got %+v\n", want, got)
	}
}

func TestSafari(t *testing.T) {
	var got *UserAgent
	want := &UserAgent{}

	got = Parse(`Mozilla/5.0 (Macintosh; Intel Mac OS X 10_6_8) AppleWebKit/537.13+ (KHTML, like Gecko) Version/5.1.7 Safari/534.57.2`)
	want.Type = Browser
	want.OS = "Mac OS X"
	want.OSVersion = mustParse("10.6.8")
	want.Name = "Safari"
	want.Version = mustParse("5.1.7")
	want.Security = SecurityUnknown
	want.Mobile = false
	if !eqUA(want, got) {
		t.Errorf("expected %+v, got %+v\n", want, got)
	}

	got = Parse(`Mozilla/5.0 (iPhone; CPU iPhone OS 6_1_4 like Mac OS X) AppleWebKit/536.26 (KHTML, like Gecko) Version/6.0 Mobile/10B350 Safari/8536.25`)
	want.Type = Browser
	want.OS = "iOS"
	want.OSVersion = mustParse("6.1.4")
	want.Name = "Safari"
	want.Version = mustParse("6.0.0")
	want.Security = SecurityUnknown
	want.Mobile = true
	if !eqUA(want, got) {
		t.Errorf("expected %+v, got %+v\n", want, got)
	}

	got = Parse(`Mozilla/5.0 (iPad; U; CPU OS 3_2 like Mac OS X; en-us) AppleWebKit/531.21.10 (KHTML, like Gecko) Version/4.0.4 Mobile/7B334b Safari/531.21.10`)
	want.Type = Browser
	want.OS = "iOS"
	want.OSVersion = mustParse("3.2")
	want.Name = "Safari"
	want.Version = mustParse("4.0.4")
	want.Security = SecurityStrong
	want.Mobile = false
	want.Tablet = true
	if !eqUA(want, got) {
		t.Errorf("expected %+v, got %+v\n", want, got)
	}
}

func TestIE(t *testing.T) {
	var got *UserAgent
	want := &UserAgent{}

	got = Parse(`Mozilla/5.0 (compatible; MSIE 10.0; Windows NT 6.1; Trident/6.0)`)
	want.Type = Browser
	want.OS = "Windows"
	want.OSVersion = mustParse("6.1")
	want.Name = "MSIE"
	want.Version = mustParse("10.0.0")
	want.Security = SecurityUnknown
	if !eqUA(want, got) {
		t.Errorf("expected %+v, got %+v\n", want, got)
	}

	got = Parse(`Mozilla/5.0 (Windows NT 6.3; Trident/7.0; .NET4.0E; .NET4.0C; rv:11.0) like Gecko`)
	want.Type = Browser
	want.OS = "Windows"
	want.OSVersion = mustParse("6.3")
	want.Name = "MSIE"
	want.Version = mustParse("11.0.0")
	want.Security = SecurityUnknown
	if !eqUA(want, got) {
		t.Errorf("expected %+v, got %+v\n", want, got)
	}
}

func TestGeneric(t *testing.T) {
	var got *UserAgent
	want := &UserAgent{}

	got = Parse(`Dillo/0.8.6-i18n-misc`)
	want.Type = Browser
	want.OS = "unknown"
	want.OSVersion = semver.Version{}
	want.Name = "Dillo"
	want.Version = mustParse("0.8.6-i18n-misc")
	want.Security = SecurityUnknown
	//	want.URL = u("http://www.dillo.org/")
	if !eqUA(want, got) {
		t.Errorf("expected %+v, got %+v\n", want, got)
	}

	got = Parse(`Googlebot/2.1 (+http://www.google.com/bot.html)`)
	want.Type = Crawler
	want.OS = "unknown"
	want.OSVersion = semver.Version{}
	want.Name = "Googlebot"
	want.Version = mustParse("2.1.0")
	want.Security = SecurityUnknown
	if !eqUA(want, got) {
		t.Errorf("expected %+v, got %+v\n", want, got)
	}
}

func TestPhantomJS(t *testing.T) {
	var got *UserAgent
	want := &UserAgent{}
	want.Mobile = false
	want.Tablet = false

	got = Parse(`Mozilla/5.0 (Macintosh; Intel Mac OS X) AppleWebKit/538.1 (KHTML, like Gecko) PhantomJS/2.0.0 Safari/538.1`)
	want.Type = Library
	want.OS = "Mac OS X"
	want.OSVersion = semver.Version{}
	want.Name = "PhantomJS"
	want.Version = mustParse("2.0.0")
	want.Security = SecurityUnknown
	if !eqUA(want, got) {
		t.Errorf("expected %+v, got %+v\n", want, got)
	}

	got = Parse(`Mozilla/5.0 (Macintosh; Intel Mac OS X) AppleWebKit/534.34 (KHTML, like Gecko) PhantomJS/1.9.0 (development) Safari/534.34`)
	want.Type = Library
	want.OS = "Mac OS X"
	want.OSVersion = semver.Version{}
	want.Name = "PhantomJS"
	want.Version = mustParse("1.9.0")
	want.Security = SecurityUnknown
	if !eqUA(want, got) {
		t.Errorf("expected %+v, got %+v\n", want, got)
	}

	got = Parse(`Mozilla/5.0 (Unknown; Linux x86_64) AppleWebKit/538.1 (KHTML, like Gecko) PhantomJS/2.1.1 Safari/538.1`)
	want.Type = Library
	want.OS = "GNU/Linux"
	want.OSVersion = semver.Version{}
	want.Name = "PhantomJS"
	want.Version = mustParse("2.1.1")
	want.Security = SecurityUnknown
	if !eqUA(want, got) {
		t.Errorf("expected %+v, got %+v\n", want, got)
	}
}

func TestAndroidAPI(t *testing.T) {
	var got *UserAgent
	want := &UserAgent{}
	want.Mobile = true
	want.Tablet = false

	got = Parse(`Dalvik/2.1.0 (Linux; U; Android 5.1.1; E6603 Build/32.0.A.6.152)`)
	want.Type = Library
	want.OS = OSAndroid
	want.OSVersion = mustParse("5.1.1")
	want.Name = "E6603"
	want.Version = mustParse("2.1.0")
	want.Security = SecurityStrong
	if !eqUA(want, got) {
		t.Errorf("expected %+v, got %+v\n", want, got)
	}

	got = Parse(`Dalvik/2.1.0 (Linux; U; Android 6.0.1; Nexus 6P Build/MHC19I)`)
	want.Type = Library
	want.OS = OSAndroid
	want.OSVersion = mustParse("6.0.1")
	want.Name = "Nexus 6P"
	want.Version = mustParse("2.1.0")
	want.Security = SecurityStrong
	if !eqUA(want, got) {
		t.Errorf("expected %+v, got %+v\n", want, got)
	}
}
