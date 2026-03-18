package forgemsg

type CmdSuccessMsg struct {
	Output string
}

type CmdErrorMsg struct {
	Error error
	Debug []string
}
