# HyperView

> Status: **Experimental**

HyperView is a Go package designed to simplify the process of rendering data into various output formats for web applications. 

## Overview

HyperView provides a fluent interface for building HTTP responses, primarily focusing on rendering data for HTML templates, JSON responses, and HTMX snippets. It uses an adapter-based architecture, allowing for flexible output formatting and easy extension to new formats.

## Goals

- Fluent API: Utilizes a builder pattern for constructing responses.
- Multi-format Support: Renders output in HTML, JSON, HTMX, and other formats through adapters.
- Template Integration: Works seamlessly with Go's `html/template` package.
- Extensibility: Allows creation of custom adapters for additional output formats.
- Performance: Designed with efficiency in mind for use in high-traffic applications.

## Use Cases

HyperView could be useful for:

- Web applications requiring multiple output formats (or just one, the fluent approach makes for a clean way to build responses)
- APIs that need to serve both HTML and JSON responses
- Projects using server-side rendering with occasional HTMX integration

## Layouts

Layouts are used to define the structure of a page, including common elements like headers, footers, and navigation. They are typically used to wrap content from a view or partial.

They are expected to be in the `layouts` directory of the configured template path.

> Each layout should be defined as a Go template file and named with the `layout:layoutName` format.

For example, the following layout file defines a `layout:base` name:

```html
{{define "layout:dashboard"}}
<!DOCTYPE html>
<html lang="en">
...
</html>
{{end}}
```

When referring to layouts, however, the `layout:` prefix is omitted. For example, to use the `base` layout in a response.

```go
resp := response.NewResponse().
    Layout("dashboard").
    Path("dashboard/account").
    Title("Current Account").
    Data(data)
```

## Partials

Partials are used to define reusable components that can be included in multiple views. They are typically used for elements like navigation menus, sidebars, and widgets.

They are expected to be in the `partials` directory of the configured template path.

Partials can be named however you like, but it is recommended to use a descriptive name that reflects the content of the partial.

For example, I like to use an `@` prefix to indicate that a file is a partial, along with the relative directory path from the `partials` directory. This makes it easy to identify partials in the template directory.

For example, if i have a partial that renders a navigation menu, it might exist in the `partials/navbar.html` file:

```html
{{define "@navbar"}}
<nav>
...
</nav>
{{end}}
```

Or, if I have a partial in a nested directory, like `partials/widgets/card.html`:

```html
{{define "@widgets/card"}}
<div class="card">
...
</div>
{{end}}
```

> For the purposes of HyperView, however, this is arbitrary and you can name your partials however you like.

## Views

Views are used to define the content of a page. They are typically used to render the main content of a page.

They are expected to be in the `views` directory of the configured template path. 

You can use any naming convention you like for defined templates in the views, but it is recommended to use a descriptive name that reflects the content of the view.

For example, I like to use a `page:` prefix to name page-related templates that are used in the layouts. I typically use the following: 

- `page:main` for the main content of a page
- `page:title` for the title of a page

```html
{{define "page:title"}}Some title{{end}}

{{define "page:main"}}
    <p>Some content</p>
{{end}}
```

These are then used in the layout file like so:

```html
{{template "page:title" .}}
{{template "page:main" .}}
```

For the purposes of HyperView, however, this is arbitrary and you can name your defined templates however you like.

When indicating the view path in a response, the `page:` prefix is omitted and only the relative directory path from the `views` directory is used. 

For example, to show the view `views/dashboard/account.html`:

```go
resp := response.NewResponse().
    Layout("dashboard").
    Path("dashboard/account").
    Title("Current Account").
    Data(data)
```
