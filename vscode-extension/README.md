# Chialisp language server

More info at github.com/clls-dev/clls

## Requirements

### clls

You need the `clls` program in your `$PATH`
- Install [go](https://golang.org/)
- Put the go binaries in your path, for example
  ```shell
  export PATH="${PATH}:`go env GOPATH`/bin"
  ```
- Get clls, make sure you are not `cd`ed in a go module
  ```shell
  cd && go get -u github.com/clls-dev/clls/cmd/clls
  ```

## Functionality

This Language Server works for .clvm files. It has the following language features:
- Semantic tokens (syntax coloring)
- Formatting (very rough for now)
- Rename (not across includes yet as it would require to parse all .clvm files in the project and it's not practical for now)
- Document highlight (highlights the symbol under the cursor throughout the document)

## Donate

xch17hh3c0kjtrrkvnjsvqu3m2wm94yavztdpdr3g8y9gncsv3t9pz2qckyfvx
