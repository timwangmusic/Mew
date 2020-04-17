package commands

type CommandBase interface {
	Validate() error
	Prepare() error
	Execute() error
}
