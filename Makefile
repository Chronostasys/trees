test: test-rbtree test-btree
bench: bench-rbtree bench-btree
test-rbtree:
	cd bitree && go test -v -cover -coverprofile=c.out && go tool cover -html=c.out -o coverage.html
bench-rbtree:
	cd bitree && go test -bench=. -run=Bench -v

test-btree:
	cd btree && go test -v -cover -coverprofile=c.out && go tool cover -html=c.out -o coverage.html
bench-btree:
	cd btree && go test -bench=. -run=Bench -v -benchtime=1000000x -benchmem

test-sql:
	cd sql && go test -v -cover -coverprofile=c.out && go tool cover -html=c.out -o coverage.html
bench-sql:
	cd sql && go test -bench=. -run=Bench -v -benchtime=1000000x -benchmem