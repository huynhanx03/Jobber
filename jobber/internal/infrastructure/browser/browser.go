package browser

import (
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/proto"
	"github.com/go-rod/stealth"
)

type Browser struct {
	Rod *rod.Browser
}

func New() *Browser {
	url := launcher.New().
		Headless(true).
		Set("disable-gpu").
		Set("no-sandbox").
		Set("disable-setuid-sandbox").
		Set("disable-blink-features", "AutomationControlled").
		Set("window-size", "1280,800").
		MustLaunch()

	b := rod.New().ControlURL(url).MustConnect()

	return &Browser{Rod: b}
}

func (b *Browser) Close() {
	b.Rod.MustClose()
}

func (b *Browser) StealthPage(url string) *rod.Page {
	page := stealth.MustPage(b.Rod)

	page.MustSetUserAgent(&proto.NetworkSetUserAgentOverride{
		UserAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/121.0.0.0 Safari/537.36",
	})

	page.MustNavigate(url).MustWaitLoad()
	return page
}

func RandomDelay(min, max time.Duration) {
	delay := min + time.Duration(time.Now().UnixNano()%int64(max-min))
	time.Sleep(delay)
}
