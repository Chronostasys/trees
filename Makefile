test: test-rbtree test-btree
bench: bench-rbtree bench-btree
test-rbtree:
	cd bitree && go test -v
bench-rbtree:
	cd bitree && go test -bench=. -run=Bench -v

test-btree:
	cd btree && go test -v
bench-btree:
	cd btree && go test -bench=. -run=Bench -v

