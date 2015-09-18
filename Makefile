default: build

build:
	gox -os="darwin linux" ./...
	test -d release || mkdir release
	rm -f release/*
	mv sapt_* release
	cd release && for b in `ls`; do zip $$b.zip $$b; done
