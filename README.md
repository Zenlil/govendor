# govendor
Basic re-implementation of govendor

```
$ govendor --help

govendor 0.3.0
Usage: govendor [--filename VENDOR_JSON] [--gopath GOPATH] [--ignore IGNORE] [--dryrun] [--verbose] <command> [<args>]

Options:
  --filename VENDOR_JSON, -f VENDOR_JSON
                         the name and location of the vendor.json file [default: vendor/vendor.json]
  --gopath GOPATH        the GOPATH location
  --ignore IGNORE        ignore these packages prefixes
  --dryrun               perform a dry-run
  --verbose, -v          show more information regarding operations
  --help, -h             display this help and exit
  --version              display version and exit

Commands:
  init                   initialized the vendor.json file
  list                   list packages along with git timestamps
  add                    add one package(s) to vendoring
  delete                 delete package(s) from vendoring
  get                    add and update package(s) to vendoring
  update                 update package(s) from GOPATH
  deps                   find dependent packages
  tidy                   add, update and removes according to 'deps'
```

## init
Create the basic vendor.json in the vendor folder

## Package operations
The `list`, `add`, `delete`, `get` and `update` commands take a list of package-names

Example:
```
$ govendor add github.com/alexflint/go-arg
...

$ govendor update github.com/dustin/go-humanize
...
```

## Dependency operations
These operations are scanning the go-source to show (and perform add, update and delete operations)

You can supply the name of multiple source-folders to allow them to share the same vendor-folder. (see example below)

The `deps` and `tidy` both scan the desired source-folders for their dependencies and present a list on what to do. The `tidy` also performs those actions.

Use the `--no-remove` argument on `tidy` to not remove unused packages

## Example: multi-source-folders
```
.../go/src/mypackage
           |-- main
           |   |-- main.go
           |   |-- ...
           |
           |-- logic
           |   |-- a.go
           |   |-- b.go
           |   |-- ...
           |
           |-- vendor
               |-- vendor.json
               |-- ...(vendored packages)
```

To scan for dependencies in the following situation you can do (from the `mypackage` folder):
```
govendor deps *
```
this will then scan both the `main` and the `logic` folders for shared their dependencies.

The `tidy *` command will also add, update or remove using the same principle.

## Example: always-present side repo
If you have a shared package thats always available during build that you're not interested in vendoring to each and every one of you other packages you can use the env-var `GOVENDOR_IGNORE`

```
...go/src
      |-- mypackage
      |   |---
      |
      |-- otherpackage
      |   |---
      |
      |-- frameworkpkg
          |---
```

Now, just set the `GOVENDOR_IGNORE=.,frameworkpkg` and any import from mypackage or otherpackage to anything inside 'frameworkpkg' will be ignored and not vendored.

The '.' in the option will also ignore vendoring between packages sharing the same git-repo.