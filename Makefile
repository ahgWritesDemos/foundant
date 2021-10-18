run:
	go run src/main.go

debug:
	dlv debug --check-go-version=false src/main.go
