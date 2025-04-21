debug:
	go build -v
	scriptprox.exe

build_release:
	go.exe build -v -ldflags "-s -w -H=windowsgui"

generate_metadata:
	git log --date=short --format="%ad %B%-C()" -n 10 HEAD > changelog.txt

release: build_release generate_metadata