debug:
	go build -v
	scriptprox.exe

build_release:
	go.exe build -v -ldflags "-s -w -H=windowsgui"

copy:
	cp.exe "resource.rpf" "Z:\home\files\web\scriptprox"
	cp.exe "scriptprox.exe" "Z:\home\files\web\scriptprox"

release: build_release copy