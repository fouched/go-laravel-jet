package celeritas

import (
	"github.com/justinas/nosurf"
	"net/http"
	"strconv"
)

func (c *Celeritas) SessionLoad(next http.Handler) http.Handler {
	c.InfoLog.Println("SessionLoad")
	return c.Session.LoadAndSave(next)
}

func (c *Celeritas) NoSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)
	secure, _ := strconv.ParseBool(c.config.cookie.secure)

	// to allow some URLS
	csrfHandler.ExemptGlob("/api/*")

	// SameSite=Strict—the cookie is only sent for requests that originate on the same domain.
	// Even arriving at the site from an off-site link will not see the cookie,
	// unless you subsequently refresh the page or navigate within the site

	// SameSite=Lax—cookie is sent if you navigate to the site through following a link from
	// another domain but not if you submit a form. This is generally what you want to protect
	// against CSRF attacks!

	// Why not habitually use SameSite=Strict? Because then if someone follows a link to your site
	// their first request will be treated as if they are not signed in at all. That’s bad!
	// So explicitly setting a cookie with SameSite=Lax should be enough to protect your application
	// from CSRF vulnerabilities... provided your users have a browser that supports it.

	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   secure,
		SameSite: http.SameSiteStrictMode,
		Domain:   c.config.cookie.domain,
	})

	return csrfHandler
}
