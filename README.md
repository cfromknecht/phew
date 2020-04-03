# phew
A hex-encoded, gzipped encoder/decoder for [spew](github.com/davecgh/spew).

When debugging go programs `spew` offers a wealth of information, but it can
also be useful in production for capturing critical information about the
execution of programs under real-world conditions.

The latter becomes somewhat challenging, however, since the `spew`
pretty-printer deviates from the standard log format where everything fits on a
single line. When printing large structs, the output can span several hundred or
thousands of lines, obliterating the readability of the production logs.

Consider the following struct:
```
type VeryLargeThing struct {
    Number      uint32
    LotsOfBytes []byte
}
```

Calling `spew.Sdump` on an instance of `VeryLargeThing` might resemble:
```
(main.VeryLargeThing) {
 Number: (uint32) 42,
 LotsOfBytes: ([]uint8) (len=8 cap=8) {
  00000000  62 69 67 20 64 61 74 61                           |big data|
 }
}
```



Often times we don't necessarily need to observe the spewed output, but rather
we would like the information to be availabe if need to investigate an issue.
However, embedding the above `Sdump` can cause production logs to become
bloated, even if there is no active investigation.

Enter `phew`: instead of logging the pretty-printed struct directly, we log
hex-encoded, gzipped output of `spew` so that the output can be unpacked later
for evaluation.

As a drop in replacement, encoding the same struct with `phew.Sdump`
produces the following output, which can be easily embedded in a single
log line:
```
1f8b08000000000000ffd2c84dccccd30b4b2daaf4492c4a4f0dc9c8cc4bd754a8e652f02bcd4d4a2db252d028cdcc2b3136d2543031d2e152f0c92f29f64f73aa2c492db652d0888e05495a682a68e4a4e6d95a28242716d85a80752b184081828299918299a58299b98291818299898299a1823998c40d6a9232d31552124b126bb8146ab96a01010000ffff2ee62bb7a3000000
```

To inspect a phewed object at a later point, one can use the `phew decode` CLI
command, which yields the original, spewed output:
```
$ phew decode 1f8b08000000000000ffd2c84dccccd30b4b2daaf4492c4a4f0dc9c8cc4bd754a8e652f02bcd4d4a2db252d028cdcc2b3136d2543031d2e152f0c92f29f64f73aa2c492db652d0888e05495a682a68e4a4e6d95a28242716d85a80752b184081828299918299a58299b98291818299898299a1823998c40d6a9232d31552124b126bb8146ab96a01010000ffff2ee62bb7a3000000
(main.VeryLargeThing) {
 Number: (uint32) 42,
 LotsOfBytes: ([]uint8) (len=8 cap=8) {
  00000000  62 69 67 20 64 61 74 61                           |big data|
 }
}
```

## CLI Installation
Using `go get`:
```
go get github.com/cfromknecht/phew
```
From source:
```
git clone git@github.com:cfromknecht/phew.git
cd phew
go install ./cmd/phew
```
