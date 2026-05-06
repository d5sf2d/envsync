package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/yourorg/envsync/internal/envfile"
)

func main() {
	diffCmd := flag.NewFlagSet("diff", flag.ExitOnError)
	syncCmd := flag.NewFlagSet("sync", flag.ExitOnError)

	// diff flags
	diffMask := diffCmd.Bool("mask", true, "mask secret values")
	diffColor := diffCmd.Bool("color", true, "colorize output")

	// sync flags
	syncDry := syncCmd.Bool("dry-run", false, "print changes without writing")
	syncOverwrite := syncCmd.Bool("overwrite", false, "overwrite changed keys")
	syncAdd := syncCmd.Bool("add", true, "add missing keys")

	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: envsync <diff|sync> [flags] <source> <destination>")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "diff":
		_ = diffCmd.Parse(os.Args[2:])
		args := diffCmd.Args()
		if len(args) != 2 {
			fmt.Fprintln(os.Stderr, "usage: envsync diff [flags] <source> <destination>")
			os.Exit(1)
		}
		runDiff(args[0], args[1], *diffMask, *diffColor)
	case "sync":
		_ = syncCmd.Parse(os.Args[2:])
		args := syncCmd.Args()
		if len(args) != 2 {
			fmt.Fprintln(os.Stderr, "usage: envsync sync [flags] <source> <destination>")
			os.Exit(1)
		}
		runSync(args[0], args[1], *syncDry, *syncOverwrite, *syncAdd)
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n", os.Args[1])
		os.Exit(1)
	}
}

func runDiff(srcPath, dstPath string, masked, color bool) {
	src, err := envfile.Parse(srcPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading %s: %v\n", srcPath, err)
		os.Exit(1)
	}
	dst, err := envfile.Parse(dstPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading %s: %v\n", dstPath, err)
		os.Exit(1)
	}
	diffs := envfile.Diff(dst, src)
	envfile.WriteDiffReport(os.Stdout, diffs, envfile.ReportOptions{Masked: masked, Color: color})
}

func runSync(srcPath, dstPath string, dry, overwrite, add bool) {
	src, err := envfile.Parse(srcPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading %s: %v\n", srcPath, err)
		os.Exit(1)
	}
	dst, err := envfile.Parse(dstPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading %s: %v\n", dstPath, err)
		os.Exit(1)
	}
	opts := envfile.SyncOptions{DryRun: dry, Overwrite: overwrite, AddMissing: add}
	res, err := envfile.Sync(dst, src, dstPath, opts)
	if err != nil {
		fmt.Fprintf(os.Stderr, "sync error: %v\n", err)
		os.Exit(1)
	}
	envfile.WriteSyncReport(os.Stdout, res)
}
