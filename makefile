
#make verison=123 buildplugin
version?=0

buildPlugin:
	cd logic && go build -buildmode=plugin -o plugin.${version}.so ./user.go

build: buildPlugin
	go build -o test main.go