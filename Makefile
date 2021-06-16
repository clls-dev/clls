pkg/lspsrv/server.gen.go:
	go run ./pkg/lspsrv/gen > $@
	goimports -w $@
.PHONY: pkg/lspsrv/server.gen.go
