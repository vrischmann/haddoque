haddoque
========

[![Build Status](https://travis-ci.org/vrischmann/haddoque.svg?branch=master)](https://travis-ci.org/vrischmann/haddoque)
[![GoDoc](https://godoc.org/github.com/vrischmann/haddoque?status.svg)](https://godoc.org/github.com/vrischmann/haddoque)

haddoque is a query engine library for JSON data.

**Note**: documentation is currently being written, please bear with me. In the meantime if you're interested, check out the tests in `engine_test.go`.

Use cases
---------

Anytime you want to pre-process arbitrary JSON data with user-defined queries, you can use haddoque.

Granted it's a limited engine, but here are a few use cases:

  * search a stream of arbitrary JSON data, such as for example a event tracker with all sorts of events.
  * quickly search through a bunch of JSON files

As an example, I wrote haddoque as a way to query data in a Kafka stream.

Supported data
--------------

haddoque is intended to be used with `map[string]interface{}` data as input. In the future I may look into supported struct-decoded data.

License
-------

haddoque is MIT licensed. See the LICENSE file for more details.
