<p align="center">
  <img width="96" height="96" src="./docs/96x96.png" alt="logo">
</p>

<h1 align="center">shmux</h1>

<p align="center">
Run multiple scripts from one file. In (almost) any language.
</p>

<p align="center">
  <a href="https://asciinema.org/a/548928" target="_blank">
    <img src="https://asciinema.org/a/548928.svg" height="288"/>
  </a>
</p>

## ⚡️ Quick start

### Installation

_MacOS and Linux_
```sh
sudo sh -c "curl -s https://shikaan.github.io/sup/install | REPO=shikaan/shmux sh -"

# or

sudo sh -c "wget -q https://shikaan.github.io/sup/install -O- | REPO=shikaan/shmux sh -"
```

_Windows and manual instructions_

Head to the [releases](https://github.com/shikaan/shmux/releases) page and download the executable for your system and architecture.

### Usage

A common use case for `shmux` is running simple scripts in a standardized and language-agnostic way. These scripts are to be found in the _configuration_ file, also known as _shmuxfile_.

For example, a `shmuxfile.sh` for a Go project might look like: 

```sh
test:
  go test ./...

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

What if we wanted to write the scripts in JavaScript? Well, you then just need a `shmuxfile.js` with a [shebang](https://en.wikipedia.org/wiki/Shebang_(Unix)) defining the interpreter to be used and you're set.

```js
greet:
  #!/usr/bin/env node

  const friend = "$1"
  const author = "$@"
  const message = friend === "darkness" 
    ? "Hello darkness, my old friend"
    : `Hello ${friend}, from ${author}`
  
  console.log(message)
```

and run it like

```bash
$ shmux greet -- "Manuel"
# => Hello Manuel, from greet
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
