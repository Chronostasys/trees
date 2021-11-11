# Bitree

好吧，这其实是颗红黑树  

一开始想先写二叉树的，写好的太快了，然后就查了下红黑树直接在原本基础上改了

## 接口
只要实现`Hasher`接口就能往里存。搜搜或者删除的时候直接用hash值来删
```golang
type Hasher interface {
	Hash() int
}
```



