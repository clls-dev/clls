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
- Semantic tokens (syntax highlighting)
- Formatting (very rough for now)