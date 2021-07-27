package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"local.go/imgconv"
)

type CLI struct {
	in       io.Reader
	out, err io.Writer
}

func main() {
	cli := &CLI{os.Stdin, os.Stdout, os.Stderr}
	os.Exit(cli.Run())
}

func (c *CLI) Run() int {
	f_in := flag.String("i", "jpg", "file format to convert from")
	f_out := flag.String("o", "png", "file format to convert to")
	f_v := flag.Bool("v", false, "show debug messages")
	flag.Parse()
	if *f_v {
		fmt.Fprintln(c.out, "in=", *f_in, ", out=", *f_out, flag.Args())
	}
	if err := imgconv.IsValidExt(*f_in); err != nil {
		fmt.Fprintln(c.err, err)
		return 1
	}
	if err := imgconv.IsValidExt(*f_out); err != nil {
		fmt.Fprintln(c.err, err)
		return 1
	}
	if flag.NArg() == 0 {
		fmt.Fprintln(c.err, "invalid argument: missing target dir")
		return 1
	}
	for _, d := range flag.Args() {
		if info, err := os.Stat(d); err != nil {
			fmt.Fprintln(c.err, "error: "+d+": no such file or directory")
			return 1
		} else if *f_v {
			fmt.Fprintln(c.err, info.Name())
		}
	}
	n := 0
	for _, d := range flag.Args() {
		err := filepath.Walk(d,
			func(p string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}
				if info.IsDir() || filepath.Ext(p) != "."+*f_in {
					return nil
				}
				if err := imgconv.Convert(p, *f_in, *f_out); err != nil {
					fmt.Fprintln(c.err, p, err)
					return nil
				}
				if *f_v {
					fmt.Fprintln(c.out, p, "\t=>\n", imgconv.ConvertExt(p, *f_out))
				}
				n += 1
				return nil
			})
		if err != nil {
			fmt.Fprintln(c.err, err)
			return 1
		}
	}
	fmt.Fprintln(c.out, n, "files converted")
	return 0
}
