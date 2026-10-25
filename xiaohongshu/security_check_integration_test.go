//go:build integration

// 集成测试：起无头浏览器，用本地 httptest 页面模拟「被重定向到安全验证页」，不触网、不需登录。
// 手动跑：go test -tags integration ./xiaohongshu/ -run TestSecurityVerification
package xiaohongshu

import (
	stderrors "errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/xpzouying/xiaohongshu-mcp/browser"
	"github.com/xpzouying/xiaohongshu-mcp/errors"
)

// 原页面先挂载、过 300ms 再被脚本跳到验证页，和线上观察到的顺序一致。
func newCaptchaRedirectServer() *httptest.Server {
	mux := http.NewServeMux()
	redirect := `<html><body><script>
		setTimeout(() => location.replace("/website-login/captcha?redirectPath=x&verifyUuid=1"), 300)
	</script></body></html>`
	mux.HandleFunc("/search_result", func(w http.ResponseWriter, _ *http.Request) { _, _ = w.Write([]byte(redirect)) })
	mux.HandleFunc("/explore/abc", func(w http.ResponseWriter, _ *http.Request) { _, _ = w.Write([]byte(redirect)) })
	mux.HandleFunc("/website-login/captcha", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`<html><body>安全验证</body></html>`))
	})
	return httptest.NewServer(mux)
}

func TestSecurityVerificationFailsFast(t *testing.T) {
	srv := newCaptchaRedirectServer()
	defer srv.Close()

	b := browser.NewBrowser(true)
	defer b.Close()
	page := b.NewPage()
	defer func() { _ = page.Close() }()

	t.Run("search", func(t *testing.T) {
		require.NoError(t, page.Navigate(srv.URL+"/search_result?keyword=x"))
		start := time.Now()
		err := waitFeedsLoaded(page, 20*time.Second)
		require.True(t, stderrors.Is(err, errors.ErrSecurityVerification), "got %v", err)
		require.Less(t, time.Since(start), 5*time.Second)
	})

	t.Run("detail", func(t *testing.T) {
		require.NoError(t, page.Navigate(srv.URL+"/explore/abc"))
		start := time.Now()
		loaded, err := waitNoteDetailLoaded(page, "abc", 15*time.Second)
		require.False(t, loaded)
		require.True(t, stderrors.Is(err, errors.ErrSecurityVerification), "got %v", err)
		require.Less(t, time.Since(start), 5*time.Second)
	})
}
