cmds = kut kutd
vfile = pkg/version/version.go
vdata = `git describe --tags`-`date -u +%Y%m%d%H%M%S`

.PHONY: $(cmds) all clean version tidy assets third package

all: $(cmds)

clean:
	rm -f $(cmds)

$(cmds): version tidy
	go build ./cmd/$@

version:
	rm -f $(vfile)
	@echo "package version" > $(vfile)
	@echo "const (" >> $(vfile)
	@echo "  Version = \"$(vdata)\"" >> $(vfile)
	@echo ")" >> $(vfile)

tidy:
	go mod tidy && go mod vendor

kutd: assets third

assets:
	cd ./cmd/kutd && \
	go-assets-builder -s=/assets/ -o assets.go assets

third:
	GOOS= GOARCH= go run cmd/prebuild/main.go

package: $(cmds)
	zip -q dist/kut-$$GOOS-$$GOARCH-$(vdata).zip kut* kutd* LICENSE README.md

test:
	go test ./pkg/*
