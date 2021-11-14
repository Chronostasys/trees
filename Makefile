test: test-rbtree
bench: bench-rbtree
test-rbtree:
	cd bitree && go test -v
bench-rbtree:
	cd bitree && go test -bench=. -run=Bench -v