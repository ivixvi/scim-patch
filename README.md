# scim-patch

[![Go Reference](https://pkg.go.dev/badge/github.com/ivixvi/scim-patch.svg)](https://pkg.go.dev/github.com/ivixvi/scim-patch)

Go implementation of SCIM2.0 Patch operation

> [!CAUTION]
> This library is not stable and not ready for production use.

# Overview

The SCIM2.0 Patch operation specification is broad, and absorbing differences between IdPs is challenging.
This library aims to handle "schema operations via Patch" comprehensively.

Since it does not directly manipulate application data, the overall processing and data storage may become redundant.
However, it helps reduce tight coupling by only requiring you to consider the mapping between the SCIM schema and the schema used in your application.

Additionally, this library depends on the following SCIM-related implementations for handling Schema and filter:

- https://github.com/elimity-com/scim
- https://github.com/scim2/filter-parser

## Expected Use Cases

The concept is close to the issue described here:
https://github.com/elimity-com/scim/issues/171

For usage examples in the current implementation, please refer to [example](./_example/README.md).

### Logger

Logging of internal processing in the Patcher can be achieved by passing a logger via context.
You can use a logger that implements the PatcherLogger interface.

For specific usage examples, please refer to [example](./_example/README.md).
