# Scriptprox2
Rewrite von [Scriptprox](http://192.168.178.80/git/patrick/Scriptprox) in Golang

## Building
```batch
make debug
make build_release
make release
```

Die `debug`-Config erstellt das Programm mit Console-Output eingeschalten und öffnet es anschließend  
Die `build_release`-Config erstellt das Programm ohne Console-Output eingeschalten und ohne Symbols  
Die `release`-Config macht alles, was die `build_release` Config macht, aber kopiert die Dateien anschließend zu `Z:\home\files\web\scriptprox`