# GoWrap
[![License](https://img.shields.io/badge/license-mit-green.svg)](https://github.com/hexdigest/gowrap/blob/master/LICENSE)
[![Build Status](https://travis-ci.org/hexdigest/gowrap.svg?branch=master)](https://travis-ci.org/hexdigest/gowrap)
[![Coverage Status](https://coveralls.io/repos/github/hexdigest/gowrap/badge.svg?branch=master)](https://coveralls.io/github/hexdigest/gowrap?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/hexdigest/gowrap?dropcache)](https://goreportcard.com/report/github.com/hexdigest/gowrap)
[![GoDoc](https://godoc.org/github.com/hexdigest/gowrap?status.svg)](http://godoc.org/github.com/hexdigest/gowrap)
[![Awesome](https://cdn.rawgit.com/sindresorhus/awesome/d7305f38d29fed78fa85652e3a63e154dd8e8829/media/badge.svg)](https://github.com/avelino/awesome-go#generation-and-generics)
[![Release](https://img.shields.io/github/release/hexdigest/gowrap.svg)](https://github.com/hexdigest/gowrap/releases/latest)

GoWrap is a command line tool that generates decorators for Go interface types using simple templates.
With GoWrap you can easily add metrics, tracing, fallbacks, pools, and many other features into your existing code in a few seconds.


## Demo

![demo](https://github.com/hexdigest/gowrap/blob/master/gowrap.gif)

## Installation

```
go get -u github.com/hexdigest/gowrap/cmd/gowrap
```

## Usage of gowrap

```
Usage: gowrap gen -p package -i interfaceName -t template -o output_file.go
  -g	don't put //go:generate instruction into the generated code
  -i string
    	the source interface name, i.e. "Reader"
  -o string
    	the output file name
  -p string
    	the source package import path, i.e. "io", "github.com/hexdigest/gowrap" or
    	a relative import path like "./generator"
  -t template
    	the template to use, it can be an HTTPS URL a local file or a
    	reference to one of the templates in the gowrap repository
  -v value
    	a key-value pair to parametrize the template,
    	arguments without an equal sign are treated as a bool values,
    	i.e. -v DecoratorName=MyDecorator -v disableChecks
```

This will generate an implementation of the io.Reader interface wrapped with prometheus metrics

```
  $ gowrap gen -p io -i Reader -t prometheus -o reader_with_metrics.go
```

This will generate a fallback decorator for the Connector interface that can be found in the ./connector subpackage:

```
  $ gowrap gen -p ./connector -i Connector -t fallback -o ./connector/with_metrics.go
```

Run `gowrap help` for more options

## Hosted templates

When you specify a template with the "-t" flag, gowrap will first search for and use the local file with this name.
If the file is not found, gowrap will look for the template [here](https://github.com/hexdigest/gowrap/tree/master/templates) and use it if found.

List of available templates:
  - [circuitbreaker](https://github.com/hexdigest/gowrap/tree/master/templates/circuitbreaker) stops executing methods of the wrapped interface after the specified number of consecutive errors and resumes execution after the specified delay
  - [fallback](https://github.com/hexdigest/gowrap/tree/master/templates/fallback) takes several implementations of the source interface and concurrently runs each implementation if the previous attempt didn't return the result in a specified period of time, it returns the first non-error result
  - [log](https://github.com/hexdigest/gowrap/tree/master/templates/log) instruments the source interface with logging using standard logger from the "log" package
  - [logrus](https://github.com/hexdigest/gowrap/tree/master/templates/logrus) instruments the source interface with logging using popular [sirupsen/logrus](https://github.com/sirupsen/logrus) logger
  - [opentracing](https://github.com/hexdigest/gowrap/tree/master/templates/opentracing) instruments the source interface with opentracing spans
  - [prometheus](https://github.com/hexdigest/gowrap/tree/master/templates/prometheus) instruments the source interface with prometheus metrics
  - [ratelimit](https://github.com/hexdigest/gowrap/tree/master/templates/ratelimit) instruments the source interface with RPS limit and concurrent calls limit
  - [retry](https://github.com/hexdigest/gowrap/tree/master/templates/retry) instruments the source interface with retries
  - [robinpool](https://github.com/hexdigest/gowrap/tree/master/templates/robinpool) puts several implementations of the source interface to the slice and for every method call it picks one implementation from the slice using the Round-robin algorithm
  - [syncpool](https://github.com/hexdigest/gowrap/tree/master/templates/syncpool) puts several implementations of the source interface to the sync.Pool and for every method call it gets one implementation from the pool and puts it back once finished

By default GoWrap places the `//go:generate` instruction into the generated code. 
This allows you to regenerate decorators' code just by typing `go generate ./...` when you change the source interface type declaration.
However if you used a remote template, the `//go:generate` instruction will contain the HTTPS URL of the template and therefore
you will need to have internet connection in order to regenerate decorators. In order to avoid this, you can copy templates from the GoWrap repository 
to local files and add them to your version control system:
```
$ gowrap template copy fallback templates/fallback
```

The above command will fetch the fallback template and copy it to the templates/fallback local file.
After template is copied, you can generate decorators using this local template:

```
$ gowrap gen -p io -i Reader -t templates/fallback reader_with_fallback.go
```

## Custom templates

You can always write your own template that will provide the desired functionality to your interfaces.
If you think that your template might be useful to others, please consider adding it to our [template repository](https://github.com/hexdigest/gowrap/tree/master/templates).

The structure of information passed to templates is documented with the [TemplateInputs](https://godoc.org/github.com/hexdigest/gowrap/generator#TemplateInputs) struct.
