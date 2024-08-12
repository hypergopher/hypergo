package funcs

import (
	"strings"
	"text/template"
	"time"
)

var FuncMap = template.FuncMap{
	// Boolean
	"yesno": YesNo,

	// Forms
	"inputAttrs": InputAttrs,

	// HTML
	"safeHTML": safeHTML,
	"safeAttr": safeAttr,
	"safeCSS":  safeCSS,
	"safeJS":   safeJS,
	"safeURL":  safeURL,

	// Maps
	"classMap": ClassMap,

	// Math
	"isEven": isEven,
	"isOdd":  isOdd,

	// Numbers
	"int": toInt64,

	// Slices
	"slice": slice,

	// Strings
	"contains":   strings.Contains,
	"hasPrefix":  strings.HasPrefix,
	"hasSuffix":  strings.HasSuffix,
	"humanize":   Humanize,
	"isBlank":    IsBlank,
	"join":       strings.Join,
	"lower":      strings.ToLower,
	"notBlank":   NotBlank,
	"pluralize":  Pluralize,
	"replaceAll": strings.ReplaceAll,
	"replace":    strings.Replace,
	"slugify":    Slugify,
	"split":      strings.Split,
	"trim":       strings.TrimSpace,
	"trimPrefix": strings.TrimPrefix,
	"trimSuffix": strings.TrimSuffix,
	"truncate":   Truncate,
	"upper":      strings.ToUpper,

	// Time
	"now":   time.Now,
	"since": time.Since,
	"until": time.Until,
}
