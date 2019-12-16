cmds = kut kutd

all: $(cmds)

clean:
	rm $(cmds)

.PHONY: $(cmds) assets

$(cmds):
	go build ./cmd/$@

kutd: assets third

assets:
	cd ./cmd/kutd && \
	go-assets-builder -s=/assets/ -o assets.go assets

third:
	cd ./cmd/kutd/assets/third && \
	python3 download.py
