# B+ 树

实现了b+树  
目前还没做持久化

## Benchmark

插入比golang的map稍快，删除、搜索稍慢

```
goos: linux
goarch: amd64
pkg: github.com/Chronostasys/trees/btree
cpu: AMD Ryzen 7 5700U with Radeon Graphics         
BenchmarkInsert
BenchmarkInsert-16               7319827               184.5 ns/op
BenchmarkMap
BenchmarkMap-16                  6172490               212.1 ns/op
BenchmarkDelete
BenchmarkDelete-16               8284734               157.9 ns/op
BenchmarkMapDelete
BenchmarkMapDelete-16           11897185               117.9 ns/op
BenchmarkSearch
BenchmarkSearch-16              10071132               128.4 ns/op
BenchmarkMapSearch
BenchmarkMapSearch-16           18661098                82.89 ns/op
```