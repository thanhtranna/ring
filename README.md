# Ring
Consistent hashing paper implementation using Red Black Tree 
![ring](./assets/consistent-hashing.png)

## Example Usage

```go
ring:=NewRing([]string{"server-1","server-2","server-3"},1)
node:=ring.Get("foo")
```


## TODO

- More test cases
- Performance test for xxhash