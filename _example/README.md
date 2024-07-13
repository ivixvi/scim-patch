# Example server

This is a sample server using an in-memory data store.
The following 4 endpoints necessary for testing PATCH operations are implemented.

- POST /Users 
- GET /Users/{id}
- PATCH /Users/{id}
- DELETE /Users/{id}

## Usage

```shell
$ cd path/to/scim-patch/_example
$ go run .\...
```
