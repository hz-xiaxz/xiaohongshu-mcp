package errors

import "errors"

var ErrNoFeeds = errors.New("没有捕获到 feeds 数据")
var ErrNoFeedDetail = errors.New("没有捕获到 feed 详情数据")

// ErrSecurityVerification 页面被小红书风控重定向到「安全验证」页（滑块/点选）。
// 无头浏览器过不了，调用方应当停下来让人在有界面的浏览器里完成验证，而不是重试。
var ErrSecurityVerification = errors.New("触发小红书安全验证")
