package main

import (
	"os"
)

func IsBackground() bool {
	return os.Getenv("GOGET_IS_BACKGROUND") == "true"
}

func IntoBackground() (int, error) {
	arg := os.Args
	files := []*os.File{
		os.Stdin,
		os.Stdout,
		os.Stderr,
	}
	logfile, err := os.Create("goget-log")
	if err != nil {
		return -1, err
	}
	files[1] = logfile
	files[2] = logfile

	attr := &os.ProcAttr{
		Files: files,
		Env:   append(os.Environ(), "GOGET_IS_BACKGROUND=true"),
	}

	proc, err := os.StartProcess(os.Args[0], arg, attr)
	if err != nil {
		panic(err)
	}
	return proc.Pid, nil
}
