main.wasm:
	GOOS=js GOARCH=wasm go build -o $@
.PHONY: main.wasm

push-site: main.wasm
	rm -fr site-repo
	mkdir site-repo
	cd site-repo && git init
	cp index.html wasm_exec.js main.wasm favicon.ico site-repo
	cd site-repo && git add . \
		&& git commit -m "first commit" \
		&& git branch -M main \
		&& git remote add origin git@github.com:clls-dev/clls-dev.github.io.git
	cd site-repo && git push -fu origin main
	rm -fr site-repo
.PHony: push-site