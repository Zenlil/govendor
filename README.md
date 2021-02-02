# govendor
Basic re-implementation of govendor

```
$ govendor --help

govendor 0.1.0
Usage: govendor [--filename VENDOR_JSON] [--gopath GOPATH] [--dryrun] <command> [<args>]

Options:
  --filename VENDOR_JSON, -f VENDOR_JSON
                         the name and location of the vendor.json file [default: vendor/vendor.json]
  --gopath GOPATH        the GOPATH location
  --dryrun               perform a dry-run
  --help, -h             display this help and exit
  --version              display version and exit

Commands:
  init                   initialized the vendor.json file
  list                   list packages along with git timestamps
  add                    add one package(s) to vendoring
  delete                 delete package(s) from vendoring
  get                    add and update package(s) to vendoring
  update                 update package(s) from GOPATH
```

## init
Create the basic vendor.json in the vendor folder

## list [packages]
Lists the packages in the vendor.json file, along with git timestamps

## add [packages]
Add packages to vendoring, does not copy any files

## delete [packages]
Delete packages from vendor.json and the vendor-folder

## get [packages]
Adds missing packages to vendor.json and updates all selected packages

## update [packages]
Copy the content of each package to the vendor folder.

Does not download anything.