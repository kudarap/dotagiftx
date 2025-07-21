package legacy

import "net/http"

const sitemapURL = "https://api.dotagiftx.com/sitemap.xml"

// pingGoogleSitemap tells google that sitemap has been updated.
func pingGoogleSitemap() error {
	_, err := http.Get("http://www.google.com/ping?sitemap=" + sitemapURL)
	return err
}
