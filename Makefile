DEPLOY_PATH = "Z:\home\files\web\scriptprox"

debug:
	go build -v
	scriptprox.exe

build_release:
	go.exe build -v -ldflags "-s -w -H=windowsgui"

generate_metadata:
	git log --date=short --format="%ad %B%-C()" -n 5 HEAD > changelog.txt

copy:
	cp.exe "resource.rpf" $(DEPLOY_PATH)
	cp.exe "scriptprox.exe" $(DEPLOY_PATH)
	cp.exe "changelog.txt" $(DEPLOY_PATH)

release: build_release generate_metadata copy