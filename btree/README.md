# B+ 树

实现了b+树  
目前还没做持久化

## Benchmark

插入删除比[谷歌的b树](https://github.com/google/btree)稍快，搜索差不多

```
goos: linux
goarch: amd64
pkg: github.com/Chronostasys/trees/btree
cpu: AMD Ryzen 7 5700U with Radeon Graphics         
BenchmarkInsert
BenchmarkInsert-16               1000000              1494 ns/op              52 B/op          1 allocs/op
BenchmarkGoogleInsert
BenchmarkGoogleInsert-16         1000000              2208 ns/op              65 B/op          1 allocs/op
BenchmarkDelete
BenchmarkDelete-16               1000000              1814 ns/op              11 B/op          0 allocs/op
BenchmarkGoogleDelete
BenchmarkGoogleDelete-16         1000000              2223 ns/op              10 B/op          0 allocs/op
BenchmarkSearch
BenchmarkSearch-16               1000000              1135 ns/op               0 B/op          0 allocs/op
BenchmarkGoogleSearch
BenchmarkGoogleSearch-16         1000000              1159 ns/op               7 B/op          0 allocs/op
```