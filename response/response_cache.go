package response

// NoCacheStrict sets the Cache-Control header to "no-cache, no-store, must-revalidate".
func (resp *Response) NoCacheStrict() {
	resp.headers["Cache-Control"] = "no-cache, no-store, must-revalidate"
}

// CacheControl sets the Cache-Control header to the given value.
func (resp *Response) CacheControl(cacheControl string) {
	resp.headers["Cache-Control"] = cacheControl
}

// ETag sets the ETag header to the given value.
func (resp *Response) ETag(etag string) {
	resp.headers["ETag"] = etag
}

// LastModified sets the Last-Modified header to the given value.
func (resp *Response) LastModified(lastModified string) {
	resp.headers["Last-Modified"] = lastModified
}
