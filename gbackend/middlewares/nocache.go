package middlewares

import (
	"github.com/gin-gonic/gin"
)

// NoCache returns a middleware that disables client and CDN caching for responses.
// It sets HTTP response headers so that browsers and intermediaries (such as CDNs/proxies)
// will not cache the API responses.
func NoCache() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Prevent caching by browsers, CDNs and other intermediaries.
		// "no-store" tells caches not to store any part of the request/response.
		// "no-cache" forces caches to revalidate with the origin before reuse.
		// "must-revalidate" and "proxy-revalidate" enforce revalidation for shared caches.
		// "max-age=0" and "s-maxage=0" ensure immediate expiration for private and shared caches.
		ctx.Header("Cache-Control", "no-store, no-cache, must-revalidate, proxy-revalidate, max-age=0, s-maxage=0")
		ctx.Header("Pragma", "no-cache")
		ctx.Header("Expires", "0")
		// Hint to CDNs (e.g. Cloudflare, Fastly) not to cache.
		ctx.Header("CDN-Cache-Control", "no-store")
		ctx.Header("Surrogate-Control", "no-store")

		ctx.Next()
	}
}
