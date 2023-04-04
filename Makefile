all:
	go build -v

releases:
	GOOS=windows GOARCH=amd64 go build -o builds/wordcrack-amd64.exe .
	GOOS=darwin GOARCH=amd64 go build -o builds/wordcrack-amd64-darwin . 
	GOOS=linux GOARCH=amd64 go build -o builds/wordcrack-amd64-linux .

clean:
	rm -rf builds
	go clean