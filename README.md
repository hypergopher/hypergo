# RenderFish

> Status: **Experimental**

RenderFish is a Go package designed to simplify the process of rendering data into various output formats for web applications. Its name is inspired by the Babel Fish from "The Hitchhiker's Guide to the Galaxy", reflecting its role in transforming data into different view outputs.

## Overview

RenderFish provides a fluent interface for building HTTP responses, primarily focusing on rendering data for HTML templates, JSON responses, and HTMX snippets. It uses an adapter-based architecture, allowing for flexible output formatting and easy extension to new formats.

## Goals

- Fluent API: Utilizes a builder pattern for constructing responses.
- Multi-format Support: Renders output in HTML, JSON, HTMX, and other formats through adapters.
- Template Integration: Works seamlessly with Go's `html/template` package.
- Extensibility: Allows creation of custom adapters for additional output formats.
- Performance: Designed with efficiency in mind for use in high-traffic applications.

## Use Cases

RenderFish could be useful for:

- Web applications requiring multiple output formats (or just one, the fluent approach makes for a clean way to build responses)
- APIs that need to serve both HTML and JSON responses
- Projects using server-side rendering with occasional HTMX integration
