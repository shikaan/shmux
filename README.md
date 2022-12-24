
<center>

![logo](./docs/96x96.png)

# shmux

Shell script multiplexer. 

**Write** and **run** multiple scripts from the same file. In (almost) any language.
</center>

## ⚡️ Quick start

### Installation

_Unix_
```sh
OS=$(uname | tr '[:upper:]' '[:lower:]')
ARCH=$([[ $(uname -m) == "x86_64" ]] && echo "amd64" || echo "386")

wget -O /usr/local/bin/shmux https://github.com/shikaan/shmux/releases/download/latest/shmux-${OS}-${ARCH}
```

_Windows and manual instructions_

Head to the [releases](https://github.com/shikaan/shmux/releases) page and download the executable for your system and architecture.

### Running scripts

A common use case for `shmux` is running simple scripts for your app in a standardised and language agnostic way. These scripts are to be found in the _configuraiton_ file, also known as _shmuxfile_.

For exmaple, a `shmux.sh` for a Go project might look like: 

```sh
build:
  go generate
  GOOS=$1 go build

greet:
  echo "Hello $1, my old friend"
```

Which can then be utilized from within the same folder as

```bash
# Runs the build command
$ shmux build -- "linux"

# Runs the greet command
$ shmux greet -- "darkness" # => Hello Darkness, my old friend
```

### Running more scripts

What if we wanted to write the scripts in JavaScript? Well, you then just need a `shmux.js` which reads something like

```js
greet:
  const friend = "$1"
  const message = friend === "darkness" 
    ? "Hello darkness, my old friend"
    : `Hello ${friend}`
  
  console.log(message)
```

and run it like

```bash
# As flags
$ shmux -c="shmux.js" -s=$(which node) greet -- "darkness"

# or from environment
export SHMUX_CONFIG="shmux.js"
export SHMUX_SHELL=$(which node)

shmux greet -- "Manuel"
```

## Documentation

More detailed documentation can be found [here](./docs/docs.md).

## FAQs

* _Which languages are supported?_
  
  `shmux` makes no assumptions about the underlying scripting language to utilize, because it always requires you to specify the shell

* _Does it have editor support?_

  As long as the language you choose is fine with having strings like `script:` in its syntax, you can just piggy-back on the existing editor support. 
  
  For example, if your _shmuxfile_ hosts JavaScript code, calling it `shmux.js` will give you decent syntax highlighting out of the box in most editors.

  More sophisticated editor support may be coming soon. If you are interested, feel free to open an issue.

## Contributing

Have a look through existing [Issues](https://github.com/shikaan/shmux/issues) and [Pull Requests](https://github.com/shikaan/shmux/pulls) that you could help with. If you'd like to request a feature or report a bug, please create a [GitHub Issue](https://github.com/shikaan/shmux/issues).

## License

[MIT](./LICENSE)