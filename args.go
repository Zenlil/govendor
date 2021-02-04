package main

import (
	"fmt"

	arg "github.com/alexflint/go-arg"
	"golang.org/x/term"
)

type InitArgs struct {
	Name string `arg:"positional"`
}

type AddOrDeleteArgs struct {
	Names []string `arg:"positional"`
}

type PackageArgs struct {
	Names []string `arg:"positional"`
}

type DepsArgs struct {
	Dump bool `arg:"--dump" help:"creates a 'deps.json' of relevant packages"`
}

type CmdArgs struct {
	Source   string           `arg:"-s,--source" help:"source location of .go-files" default:"."`
	Filename string           `arg:"-f,--filename" help:"the name and location of the vendor.json file" default:"vendor/vendor.json" placeholder:"VENDOR_JSON"`
	GoPath   string           `arg:"env" help:"the GOPATH location" placeholder:"GOPATH"`
	DryRun   bool             `arg:"--dryrun" help:"perform a dry-run"`
	Init     *InitArgs        `arg:"subcommand:init" help:"initialized the vendor.json file"`
	List     *PackageArgs     `arg:"subcommand:list" help:"list packages along with git timestamps"`
	Add      *AddOrDeleteArgs `arg:"subcommand:add" help:"add one package(s) to vendoring"`
	Delete   *AddOrDeleteArgs `arg:"subcommand:delete" help:"delete package(s) from vendoring"`
	Get      *AddOrDeleteArgs `arg:"subcommand:get" help:"add and update package(s) to vendoring"`
	Update   *PackageArgs     `arg:"subcommand:update" help:"update package(s) from GOPATH"`
	Deps     *DepsArgs        `arg:"subcommand:deps" help:"find dependent packages"`
	Tidy     *DepsArgs        `arg:"subcommand:tidy" help:"add, update and removes according to 'deps'"`
}

func (CmdArgs) Version() string {
	return "govendor 0.2.0"
}

var args CmdArgs

func init() {
	pa := arg.MustParse(&args)
	if pa.Subcommand() == nil {
		pa.Fail("missing command: init, list, add, delete, get, update")
	}

	w, h, err := term.GetSize(0)
	if err != nil {
		fmt.Printf("term.GetSize-error: %s", err)
	}
	if w < 5 {
		w = 80
	}
	console.w = w
	console.h = h
}
