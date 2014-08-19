netdbg
------

Network protocol proxy/debugger. A transparent proxy between clients and
servers to ease protocol debugging and degenerate case testing.

Very early/experimental. Not even a toy at this point.

```sh
go get github.com/schmichael/netdbg/cmd/netdbg

netdbg prog localhost:8080 google.com:80

# In another terminal
curl localhost:8080 > /dev/null
```
