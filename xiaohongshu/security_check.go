package xiaohongshu

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/go-rod/rod"
	"github.com/xpzouying/xiaohongshu-mcp/errors"
)

// securityCheckPath 风控触发时页面被重定向到的路径，
// 形如 /website-login/captcha?redirectPath=<原页面>&verifyUuid=...
const securityCheckPath = "/website-login/captcha"

// checkSecurityVerification 当前页已被重定向到安全验证页时返回 ErrSecurityVerification。
//
// 重定向后 __INITIAL_STATE__ 永远等不到数据，若只靠轮询注水，要白等到超时才报一个
// 看不出原因的「没有捕获到数据」。读不到地址时不下结论，交给调用方原有的超时判定。
func checkSecurityVerification(page *rod.Page) error {
	res, err := page.Eval(`() => location.href`)
	if err != nil {
		return nil
	}
	href := res.Value.Str()
	if !isSecurityVerificationURL(href) {
		return nil
	}
	return fmt.Errorf("%w（页面被重定向到 %s），需要在有界面的浏览器里完成验证", errors.ErrSecurityVerification, href)
}

func isSecurityVerificationURL(href string) bool {
	u, err := url.Parse(href)
	return err == nil && strings.HasPrefix(u.Path, securityCheckPath)
}
