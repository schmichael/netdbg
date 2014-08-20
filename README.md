netdbg
------

Network protocol proxy/debugger. A transparent proxy between clients and
servers to ease protocol debugging and degenerate case testing.

Very early/experimental. Not even a toy at this point.

```sh
go get github.com/schmichael/netdbg/cmd/netdbg

# Prints every time a packet is sent or received
netdbg prog,prog localhost:8080 google.com:80

# Prints escaped version of sent and received data
netdbg log,log localhost:8080 google.com:80

# In another terminal
curl localhost:8080 > /dev/null
```

Filters
-------

Filters are implemented as an interface which implements Accept and
Close methods. However, the most important part of a filters
functionality is that it's given in and out chans for manipulating the
data stream.

Each filter can be used on data the *client sends* or on data the *server sends*.

To cleanly shutdown the filter pipeline a filter must close its `out`
chan when their `in` chan is closed. 
