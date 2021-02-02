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

type CmdArgs struct {
	Filename string           `arg:"-f,--filename" help:"the name and location of the vendor.json file" default:"vendor/vendor.json" placeholder:"VENDOR_JSON"`
	GoPath   string           `arg:"env" help:"the GOPATH location" placeholder:"GOPATH"`
	DryRun   bool             `arg:"--dryrun" help:"perform a dry-run"`
	Init     *InitArgs        `arg:"subcommand:init" help:"initialized the vendor.json file"`
	List     *PackageArgs     `arg:"subcommand:list" help:"list packages along with git timestamps"`
	Add      *AddOrDeleteArgs `arg:"subcommand:add" help:"add one package(s) to vendoring"`
	Delete   *AddOrDeleteArgs `arg:"subcommand:delete" help:"delete package(s) from vendoring"`
	Get      *AddOrDeleteArgs `arg:"subcommand:get" help:"add and update package(s) to vendoring"`
	Update   *PackageArgs     `arg:"subcommand:update" help:"update package(s) from GOPATH"`
	//Check    *PackageArgs     `arg:"subcommand:check"`
}

func (CmdArgs) Version() string {
	return "govendor 0.1.0"
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
	console.w = w
	console.h = h
}
