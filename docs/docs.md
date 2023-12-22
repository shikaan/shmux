# shmux

### Table of contents

* [Configuration](#configuration)
* [Runtime](#runtime)
* [Environment, flags, and defaults](#environment-flags-and-defaults)

## Configuration

The scripts run by `shmux` live in files called _shmuxfiles_. These configuration files follow the following pattern:

* lines starting with a non-white space and ending with a `:` will be interpreted as a _script definition_
    * if the colon is followed by space-separated words, they are treated as _script dependencies_ (i.e., scripts running before the invoked one)
* non-empty lines prepended with whitespaces are considered _script lines_ 
* the other lines are ignored

A script is composed of all the lines in between two script definitions or last script definition and EOF. 

In a nutshell, `shmux` is not opinionated about which languages the script are written in and - so long as the syntax allows[^1] - editor support comes out of the box. Calling the shmuxfile with the most common extension of your language of choice, will make it significantly easier.

For example, a shmuxfile with bash scripts can be called `shmuxfile.bash`. If it was with JavaScript scripts, it can be called `shmuxfile.js`. This will yield pretty decent syntax highlighting.

At this point, `shmux` is known to be working with:
* sh and derviatives (bash, dash, fish, zsh...)
* JavaScript / TypeScript (with ts-node)
* Perl
* Python
* Ruby

If you need more sophisticated tooling, please [open an Issue](https://github.com/shikaan/shmux/issues).

[^1]: Namely, permits intendations and presence of the `script:` labels.

## Runtime

All the scripts are executed in isolation. Under the hood, `shmux` parses the file, creates a temporary file with it's content and runs it with the specified shell.

This means that all the lines in the same script share scope, as if they were on a single file.

In the runtime, scripts have the following variables available

| Variable    | Description                                                         |
|---          |---                                                                  |
| `$1`..`$9`  | Respectively the first 9 arguments passed after the `--` separator  |
| `$@`        | Holds the name of the current running script                        |

## Environment, flags, and defaults

The general rule is that as little configuration as possible should be provided for `shmux` to run. It is in fact possible to provide no configuration and have `shmux` operating on sensible defaults most of the times. However, `shmux` also provides means to customise its behaviour, namely CLI flags and environment variables. 

Hierarachy for those configuration points goes as follows: inline configuration (when applicable) takes precedence over everything, CLI flags override environment variables, and lack of any of them will make `shmux` operate on defaults.

In short: `inline configuration > CLI flags > environment variables > defaults` where the `>` means "takes precedence over".

| CLI Flag          | Environment Variable  | Default               | Description                                                 |
|---                |---                    | ---                   | ---                                                         |
| `-configuration`  | `SHMUX_CONFIG`        | closest `shmuxfile.*` | Location of the _shmuxfile_.                                |
| `-shell`          | `SHMUX_SHELL`         | current `$SHELL`      | Interpreter to run the script. Overriden by inline shebang. |