build: 
	go build -o server main.go

run: 
	build ./server
	
watch:
	reflex -g '*.go' go run main.go --start-service