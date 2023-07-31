cmds = kwt kwtd
vfile = pkg/version/version.go
vdata = `git describe --tags`-`git log -1 --format=%cd --date=format:"%Y%m%d%H%M%S"`

.PHONY: $(cmds) all clean version tidy assets third package

all: $(cmds)

clean:
	rm -f $(cmds) $(addsuffix .exe, $(cmds))

$(cmds): version tidy
	go build -mod=mod ./cmd/$@

version:
	rm -f $(vfile)
	@echo "package version" > $(vfile)
	@echo "const (" >> $(vfile)
	@echo "  Version = \"$(vdata)\"" >> $(vfile)
	@echo ")" >> $(vfile)

tidy:
	go mod tidy && go mod vendor && go fmt ./pkg/* ./cmd/*

kwtd: third assets

assets:
	go install -mod=mod github.com/jessevdk/go-assets-builder && \
	go mod vendor && \
	cd ./cmd/kwtd && \
	go-assets-builder -s=/assets/ -o assets.go assets

third:
	GOOS= GOARCH= go run cmd/prebuild/main.go

package: $(cmds)
	zip -q dist/kwt-$$GOOS-$$GOARCH-$(vdata).zip kwt* kwtd* LICENSE README.md

test: version
	go test ./pkg/*
