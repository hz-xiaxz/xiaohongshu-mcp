package xiaohongshu

import "testing"

func TestIsSecurityVerificationURL(t *testing.T) {
	cases := map[string]bool{
		"https://www.xiaohongshu.com/website-login/captcha?redirectPath=https%3A%2F%2Fwww.xiaohongshu.com%2Fsearch_result%3Fkeyword%3D%25E7%2589%259B%25E6%258E%2592&verifyUuid=abc&verifyType=124": true,
		"https://www.xiaohongshu.com/search_result?keyword=%E7%89%9B%E6%8E%92&source=web_explore_feed":                                                                                              false,
		"https://www.xiaohongshu.com/explore/abc?xsec_token=x&redirect=/website-login/captcha":                                                                                                      false,
		"about:blank": false,
		"":            false,
	}
	for href, want := range cases {
		if got := isSecurityVerificationURL(href); got != want {
			t.Errorf("isSecurityVerificationURL(%q) = %v, want %v", href, got, want)
		}
	}
}
