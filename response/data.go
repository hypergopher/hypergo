package response

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/hypergopher/renderfish/constants"
	"github.com/hypergopher/renderfish/htmx"
	"github.com/hypergopher/renderfish/request"
)

// Data is the struct that all view models must implement. It provides common data for all templates
// and represents the data that is passed to the template.
//
// This is a short-lived object that is used to work with data passed to the template. It is not thread-safe.
//
// Data should not be used directly. Instead, use the NewData function to create an instance
// of Data that contains the data you want to pass to the template.
//
// Example: NewData(request, map[string]any{"title": "Hello World"})
//
// The environment variables are generally added via config/config, but if you're not using that package,
// you can set them manually.
//
//goland:noinspection GoNameStartsWithPackageName
type Data struct {
	title       string
	request     *http.Request
	pageData    map[string]any
	csrfToken   string
	environment string
}

// NewData creates a new Data instance.
// If you are using this outside the normal RenderFish rendering process, be sure to set the request manually
// via Data.SetRequest as the request is deliberately set later in the normal rendering flow.
func NewData(pageData map[string]any) *Data {
	pageData = initData(pageData)
	return &Data{
		pageData: pageData,
	}
}

// SetTitle sets the title of the page.
func (v *Data) SetTitle(title string) {
	v.title = title
}

// SetRequest sets the request for the Data instance.
func (v *Data) SetRequest(r *http.Request) {
	v.request = r
}

func initData(data map[string]any) map[string]any {
	if data == nil {
		data = map[string]any{}
	}

	// If no "Error" key is set, set it to an empty string
	if _, ok := data["Error"]; !ok {
		data["Error"] = ""
	}

	// If no "Errors" key is set, set it to an empty map
	if _, ok := data["Errors"]; !ok {
		data["Errors"] = map[string]string{}
	}

	return data
}

// Data returns the data map that will be passed to the template.
func (v *Data) Data() map[string]any {
	v.pageData = initData(v.pageData)
	v.pageData["View"] = v
	return v.pageData
}

// AddData adds a map of data to the existing view data model.
func (v *Data) AddData(data map[string]any) {
	for key, value := range data {
		v.pageData[key] = value
	}
}

// AddDataItem adds a single key-value pair to the existing view data model.
func (v *Data) AddDataItem(key string, value any) {
	v.pageData[key] = value
}

// AddErrors adds an error message and a map of field errors to the view data model.
func (v *Data) AddErrors(msg string, fieldErrors map[string]string) {
	v.pageData["Error"] = msg
	v.pageData["Errors"] = fieldErrors
}

// Get returns the value of the specified key from the view data model.
func (v *Data) Get(key string) any {
	val, ok := v.pageData[key]
	if ok {
		return val
	}

	return ""
}

// GetString returns the value of the specified key from the view data model as a string.
func (v *Data) GetString(key string) string {
	val, ok := v.Get(key).(string)
	if ok {
		return val
	}

	return ""
}

// Title returns the title of the page.
func (v *Data) Title() string {
	return v.title
}

// ------ Error Helpers --------

// HasError returns true if the view data model contains an error message.
func (v *Data) HasError() bool {
	return v.GetString("Error") != ""
}

// Error returns the error message from the view data model.
func (v *Data) Error() string {
	return v.GetString("Error")
}

// HasErrors returns true if the view data model contains field errors.
func (v *Data) HasErrors() bool {
	return len(v.Errors()) > 0
}

// Errors returns a map of field errors from the view data model.
func (v *Data) Errors() map[string]string {
	val, ok := v.Get("Errors").(map[string]string)
	if ok {
		return val
	}

	return map[string]string{}
}

// ------ Common Helpers --------

// BaseURL returns the base URL of the request.
func (v *Data) BaseURL() string {
	return request.BaseURL(v.request)
}

// Context returns the context of the request.
func (v *Data) Context() context.Context {
	return v.request.Context()
}

// CurrentYear returns the current year.
func (v *Data) CurrentYear() int {
	return time.Now().Year()
}

// Nonce returns the nonce value from the request context, if available.
func (v *Data) Nonce() string {
	nonce, ok := v.request.Context().Value(constants.NonceContextKey).(string)
	if ok {
		return nonce
	}

	return ""
}

// HTMXNonce returns the HTMX nonce value from the request context, if available.
// This adds the inlineScriptNonce key to a JSON object with the nonce value and can be used in an HTMX meta tag.
func (v *Data) HTMXNonce() string {
	return fmt.Sprintf("{\"includeIndicatorStyles\":false,\"inlineScriptNonce\": \"%s\"}", v.Nonce())
}

// RequestPath returns the path of the request.
func (v *Data) RequestPath() string {
	return request.URLPath(v.request)
}

// RequestMethod returns the method of the request.
func (v *Data) RequestMethod() string {
	return request.Method(v.request)
}

// IsHtmxRequest returns true if the request is an HTMX request, but not a boosted request.
func (v *Data) IsHtmxRequest() bool {
	return htmx.IsHtmxRequest(v.request)
}

// IsBoostedRequest returns true if the request is a boosted request.
func (v *Data) IsBoostedRequest() bool {
	return htmx.IsBoostedRequest(v.request)
}
