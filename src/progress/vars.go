package progress

import "os"

var output *os.File = os.Stdout

/*
set the output of the bar following the choosen path
on stdout by default
can be put on "/dev/null" or "NUL" to silence
*/
func SetOutput(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	output = file
	return nil
}

func GetOuptut() string {
	return output.Name()
}
