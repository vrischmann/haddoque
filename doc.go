/*
Package haddoque implements a query engine for JSON-decoded map.

The idea is to be able to quickly query an arbitrary JSON-decoded map, without having to resort to write a script, or import the data into a database.

Query syntax

The syntax is loosely based on SQL and Go's text/template.

A query is composed of:

  - a mandatory field selector
  - optional filter conditions

For example:

    .name where (.id == 1)

Field selector

A field selector is a representation of the path to get to the field you want in the map.

In other words, ".data.id" means getting the field map["data"]["id"].

There's a special case with the "." field selector: it returns the entire source map

Condition

A condition is a binary expression, which has to evaluate to true for the query to return something.

Example:

    .id == 1
    .names contains ["foo", "bar"]
    .data.version < 10

Those are all valid conditions.

To be valid in a query though you need to enclose all expressions into parentheses like we did above.
This is a limitation of the engine that may or may not be removed in the future.

Combining conditions

You can combine any number of conditions, like so:

    . where ( (.id == 1) and (.name == "foobar) ) or (.mobile.type == "iPhone")
*/
package haddoque
