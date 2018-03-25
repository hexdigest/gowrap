package gowrap

// CommandLineError is returned from the commands when invalid command line parameters are passed
type CommandLineError string

func (e CommandLineError) Error() string {
	return string(e)
}
