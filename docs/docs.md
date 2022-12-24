# Documentation

<details>
<summary>Table of contents</summary>

* [Configuration](#configuration)
* [Runtime](#runtime)

</details>
<br>

## Configuration

The scripts run by `shmux` live in files called _shmuxfiles_. These configuration files follow the following pattern:

* lines starting with a non-white space and ending with a `:` will be interpreted as a _script definition_
* lines startung with a `\t` are considered _script lines_ 
* the other lines are ignored

A script is composed of all the lines in between two script definitions or last script definition and EOF. 

In a nutshell, `shmux` is not opinionated about which languages the script are written in and - so long as the syntax allows[^1] - editor support comes out of the box. Calling the shmuxfile with the most common extension of your language of choice, will make it significantly easier.

For example, a shmuxfile with bash scripts can be called `shmux.bash`. If it was with JavaScript scripts, it can be called `shmux.js`. This will yield pretty decent syntax highlighting.

If you need more sophisticated tooling, please [open an Issue](https://github.com/shikaan/shmux/issues).

[^1]: Namely, permits intendations and presence of `script:`.

## Runtime

All the scripts are executed in isolation. Under the hood, `shmux` parses the file, creates a temporary file with it's content and runs it with the specified shell.

This means that all the line in the same script share scope, as if they were on a single file.

In the runtime, scripts have the following variables available

| Variable    | Description   |
|---          |--- |
| `$1`..`$9`  | Respectively the first 9 arguments passed after the `--` separator
| `$@`        | Holds the name of the current running script
