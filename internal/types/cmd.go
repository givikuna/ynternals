package types

type CommandName string

const (
	CmdAdd          CommandName = "add"
	CmdRemove       CommandName = "remove"
	CmdValidate     CommandName = "validate"
	CmdList         CommandName = "list"
	CmdPeek         CommandName = "peek"
	CmdVersionShort CommandName = "-v"
	CmdVersionLong  CommandName = "--version"
)

type ExitOutcome struct {
	Stdout   string
	Stderr   string
	ExitCode int
}

type Invocation struct {
	Env     EnvConfig
	Command CommandName
	Args    []string
}

type CommandHandler func(Invocation) ExitOutcome
