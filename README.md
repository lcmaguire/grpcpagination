# grpcpagination

package to paginate gRPC List endpoints that paginate via NextPageToken

## Requirements

following goolges aip's in particular 

- [AIP-132](https://google.aip.dev/132) for list endpoints. 
- [AIP-158](https://google.aip.dev/158) for pagination.

using the [protoc-gen-go-setters](https://github.com/lcmaguire/protoc-gen-go-setters) plugin.


## install

```sh
go get github.com/lcmaguire/grpcpagination@latest
```

## Purpose 

enables you to paginate any gRPC list endpoint across your team/organization via `next_page_token`

