run: main.wasm node_modules
	npx parcel index.html

build: site_repo main.wasm node_modules
	npx parcel build index.html
.PHONY: build

site_repo:
	mkdir site_repo

main.wasm:
	GOOS=js GOARCH=wasm go build -o $@
.PHONY: main.wasm

publish:
	rm -fr site_repo
	$(MAKE) build
	cd site_repo \
		&& git init \
		&& git add . \
		&& git commit -m "first commit" \
		&& git branch -M main \
		&& git remote add origin git@github.com:clls-dev/clls-dev.github.io.git
	cd site_repo && git push -fu origin main
.PHONY: publish

node_modules:
	yarn
.PHONY: node_modules

clean:
	rm -fr .parcel-cache dist node_modules site_repo
.PHONY: clean