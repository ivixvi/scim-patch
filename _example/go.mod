module github.com/ivixvi/scim-patch/_example

go 1.22.5

replace github.com/ivixvi/scim-patch => ../

require (
	github.com/elimity-com/scim v0.0.0-20240320110924-172bf2aee9c8
	github.com/ivixvi/scim-patch v0.0.0-00010101000000-000000000000
)

require (
	github.com/di-wu/parser v0.2.2 // indirect
	github.com/di-wu/xsd-datetime v1.0.0 // indirect
	github.com/scim2/filter-parser/v2 v2.2.0 // indirect
)
