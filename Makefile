cmds = kut kutd
vfile = pkg/version/version.go
vdata = `git describe --tags`-`date -u +%Y%m%d%H%M%S`

all: $(cmds)

clean:
	rm $(cmds)
	rm *.zip

.PHONY: $(cmds) assets

$(cmds): version
	go build ./cmd/$@

kutd: assets third

assets:
	cd ./cmd/kutd && \
	go-assets-builder -s=/assets/ -o assets.go assets

third:
	go fmt cmd/prebuild/main.go

version:
	rm -f $(vfile)
	@echo "package version" > $(vfile)
	@echo "const (" >> $(vfile)
	@echo "  Version = \"$(vdata)\"" >> $(vfile)
	@echo ")" >> $(vfile)

package: $(cmds)
	zip -q dist/kut-$$PLATFORM-$(vdata).zip kut* kutd* LICENSE README.md
