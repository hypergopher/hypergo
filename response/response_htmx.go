package response

import (
	"github.com/hypergopher/hypergo/htmx"
	"github.com/hypergopher/hypergo/htmx/location"
	"github.com/hypergopher/hypergo/htmx/swap"
)

// HxLocation sets the HX-Location header, which instructs the browser to navigate to the given path without reloading the page.
//
// For simple navigations, use a path and no HxLocation options.
// For more complex navigations, use the HxLocation options to fine-tune the navigation.
//
// Sets the HX-Location header with the given path.
//
// For more information, see: https://htmx.org/headers/hx-location
func (resp *Response) HxLocation(path string, opt ...location.Option) *Response {
	if opt == nil {
		resp.headers[htmx.HXLocation] = path
	} else {
		loc := location.NewLocation(path, opt...)
		resp.headers[htmx.HXLocation] = loc.String()
	}
	return resp
}

// HxPushURL sets the HX-Push-Url header, which instructs the browser to navigate to the given path without reloading the page.
//
// To prevent the browser from updating the page, set the HX-Push-Url header to an empty string or "false".
//
// For more information, see: https://htmx.org/headers/hx-push-url
func (resp *Response) HxPushURL(path string) *Response {
	resp.headers[htmx.HXPushURL] = path
	return resp
}

// HxNoPushURL prevents the browser from updating the history stack by setting the HX-Push-Url header to "false".
//
// For more information, see: https://htmx.org/headers/hx-no-push-url
func (resp *Response) HxNoPushURL() *Response {
	resp.headers[htmx.HXPushURL] = "false"
	return resp
}

// HxRedirect sets the HX-Redirect header, which instructs the browser to navigate to the given path (this will reload the page).
//
// For more information, see: https://htmx.org/reference/#response_headers
func (resp *Response) HxRedirect(path string) *Response {
	resp.headers[htmx.HXRedirect] = path
	return resp
}

// HxRefresh sets the HX-Refresh header, which instructs the browser to reload the page.
//
// For more information, see: https://htmx.org/reference/#response_headers
func (resp *Response) HxRefresh() *Response {
	resp.headers[htmx.HXRefresh] = "true"
	return resp
}

// HxNoRefresh prevents the browser from reloading the page by setting the HX-Refresh header to "false".
//
// For more information, see: https://htmx.org/reference/#response_headers
func (resp *Response) HxNoRefresh() *Response {
	resp.headers[htmx.HXRefresh] = "false"
	return resp
}

// HxReplaceURL sets the HX-Replace-Url header, which instructs the browser to replace the history stack with the given path.
//
// For more information, see: https://htmx.org/headers/hx-replace-url
func (resp *Response) HxReplaceURL(path string) *Response {
	resp.headers[htmx.HXReplaceURL] = path
	return resp
}

// HxNoReplaceURL prevents the browser from updating the history stack by setting the HX-Replace-Url header to "false".
//
// For more information, see: https://htmx.org/headers/hx-replace-url
func (resp *Response) HxNoReplaceURL() *Response {
	resp.headers[htmx.HXReplaceURL] = "false"
	return resp
}

// HxReswap sets the HX-Reswap header, which instructs HTMX to change the swap behavior of the target element.
//
// For more information, see: https://htmx.org/attributes/hx-swap
func (resp *Response) HxReswap(swap *swap.Style) *Response {
	resp.headers[htmx.HXReswap] = swap.String()
	return resp
}

// HxRetarget sets the HX-Retarget header, which instructs HTMX to update the target element.
//
// For more information, see: https://htmx.org/reference/#response_headers
func (resp *Response) HxRetarget(target string) *Response {
	resp.headers[htmx.HXRetarget] = target
	return resp
}

// HxReselect sets the HX-ReSelect header, which instructs HTMX to update which part of the response is selected.
//
// For more information, see: https://htmx.org/reference/#response_headers and https://htmx.org/attributes/hx-select
func (resp *Response) HxReselect(reselect string) *Response {
	resp.headers[htmx.HXReselect] = reselect
	return resp
}

// HxTrigger sets a HX-Trigger header
//
// For more information, see: https://htmx.org/headers/hx-trigger/
func (resp *Response) HxTrigger(event string, value any) *Response {
	resp.triggers.Set(event, value)
	return resp
}

// HxTriggerAfterSettle sets a HX-Trigger-After-Settle header
//
// For more information, see: https://htmx.org/headers/hx-trigger/
func (resp *Response) HxTriggerAfterSettle(event string, value any) *Response {
	resp.triggers.SetAfterSettle(event, value)
	return resp
}

// HxTriggerAfterSwap sets a HX-Trigger-After-Swap header
//
// For more information, see: https://htmx.org/headers/hx-trigger/
func (resp *Response) HxTriggerAfterSwap(event string, value any) *Response {
	resp.triggers.SetAfterSwap(event, value)
	return resp
}
