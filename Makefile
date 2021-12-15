test: test-rbtree test-btree
bench: bench-rbtree bench-btree
test-rbtree:
	cd bitree && go test -v -cover
bench-rbtree:
	cd bitree && go test -bench=. -run=Bench -v

test-btree:
	cd btree && go test -v -cover
bench-btree:
	cd btree && go test -bench=. -run=Bench -v -benchtime=1000000x -benchmem

test-sql:
	cd sql && go test -v -cover
bench-sql:
	cd sql && go test -bench=. -run=Bench -v -benchtime=1000000x -benchmem