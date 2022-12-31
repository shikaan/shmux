<p align="center">
  <img width="96" height="96" src="./docs/96x96.png" alt="logo">
</p>

<h1 align="center">shmux</h1>

<p align="center">
Shell script multiplexer.
<b>Write</b> and <b>run</b> multiple scripts from the same file. In (almost) any language.
</p>

## ⚡️ Quick start

### Installation

_Unix_
```sh
OS=$(uname | tr '[:upper:]' '[:lower:]')
ARCH=$([[ $(uname -m) == "x86_64" ]] && echo "amd64" || echo "386")

wget -O /usr/local/bin/shmux https://github.com/shikaan/shmux/releases/latest/download/shmux-${OS}-${ARCH}
chmod u+x /usr/local/bin/shmux
```

_Windows and manual instructions_

Head to the [releases](https://github.com/shikaan/shmux/releases) page and download the executable for your system and architecture.

### Usage

A common use case for `shmux` is running simple scripts in a standardized and language-agnostic way. These scripts are to be found in the _configuration_ file, also known as _shmuxfile_.

For example, a `shmuxfile.sh` for a Go project might look like: 

```sh
build:
  go generate
  GOOS=$1 go build

greet:
  echo "Hello $1, my old friend"
```

Which can then be utilized as

```bash
# Runs the test command
$ shmux test

# Runs the build command with "linux" as $1
$ shmux build -- "linux"

# Runs the greet command with "darkness" as $1
$ shmux greet -- "darkness" 
# => Hello darkness, my old friend
```

### More Usage

What if we wanted to write the scripts in JavaScript? Well, you then just need a `shmuxfile.js` which reads something like

```js
greet:
  const friend = "$1"
  const author = "$@"
  const message = friend === "darkness" 
    ? "Hello darkness, my old friend"
    : `Hello ${friend}, from ${author}`
  
  console.log(message)
```

and run it like

```bash
# As flags
$ shmux -shell=$(which node) greet -- "Manuel"
# => Hello Manuel, from greet

# or from environment
export SHMUX_SHELL=$(which node)

shmux greet -- "Manuel"
```

## 📄 Documentation

More detailed documentation can be found [here](./docs/docs.md).

## ❓ FAQs

* _Which languages are supported?_
  
  `shmux` makes no assumptions about the underlying scripting language to utilize, because it always requires you to specify the shell. Any language whose syntax is compatible with shmuxfiles' requirements is supported.

* _Does it have editor support?_

  As long as the language you choose is fine with having strings like `script:` in its syntax, you can just piggy-back on the existing editor support. 
  
  For example, if your _shmuxfile_ hosts JavaScript code, calling it `shmuxfile.js` will give you decent syntax highlighting out of the box in most editors.

  More sophisticated editor support may be coming soon. If you are interested, feel free to open an issue.

## 🤓 Contributing

Have a look through existing [Issues](https://github.com/shikaan/shmux/issues) and [Pull Requests](https://github.com/shikaan/shmux/pulls) that you could help with. If you'd like to request a feature or report a bug, please create a [GitHub Issue](https://github.com/shikaan/shmux/issues).

## License

[MIT](./LICENSE)