prepare:
	go mod tidy
run: prepare
	echo "Starting server..."
	go run main.go