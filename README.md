# scim-patch
Go implementation of SCIM 2.0 Patch operations.

> [!CAUTION]
> It has not been fully implemented and is not ready for use.

# Overview

The specification of SCIM 2.0 Patch operations is broad, and absorbing the differences for each IdP is also challenging.
Therefore, this library aims to comprehensively handle "schema manipulation via Patch".

Since this does not directly manipulate the application's data, overall processing and data storage may become redundant.
However, instead of that, you only need to consider mapping between the SCIM schema and the schema used in your application, helping to reduce tight coupling.

Additionally, this library depends on the following SCIM-related implementations for handling Schema and filters:

- https://github.com/elimity-com/scim
- https://github.com/scim2/filter-parser

### Expected Usage Example

The following issue is relevant, and we aim to implement it to be usable in a form similar to this example:
https://github.com/elimity-com/scim/issues/171

For an example of usage in the current implementation, please refer to [example](./_example/README.md).
