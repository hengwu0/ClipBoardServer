all:
	GOPATH=`pwd` GOARCH=amd64 go build -o ClipBoardServer Clip
test:
	GOPATH=`pwd` GOARCH=amd64 go test -count=1 -v Algorithm
