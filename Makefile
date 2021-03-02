build: longshore

longshore: *.go
	CGO_ENABLED=0 go build -ldflags="-s -w" -v -o longshore

pack: upx

upx: longshore
	upx -q -q -9 longshore
