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
BenchmarkInsert-16               1000000              1201 ns/op
BenchmarkGoogleInsert
BenchmarkGoogleInsert-16         1000000              1300 ns/op
BenchmarkDelete
BenchmarkDelete-16               1000000              1124 ns/op
BenchmarkGoogleDelete
BenchmarkGoogleDelete-16         1000000              1361 ns/op
BenchmarkSearch
BenchmarkSearch-16               1237681              1035 ns/op
BenchmarkGoogleSearch
BenchmarkGoogleSearch-16         1000000              1042 ns/op
```