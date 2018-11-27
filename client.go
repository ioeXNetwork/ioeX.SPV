package main

import (
	"os"

	"github.com/ioeXNetwork/ioeX.SPV/log"
	"github.com/ioeXNetwork/ioeX.SPV/spvwallet/cli/account"
	"github.com/ioeXNetwork/ioeX.SPV/spvwallet/cli/transaction"
	"github.com/ioeXNetwork/ioeX.SPV/spvwallet/cli/wallet"
	"github.com/ioeXNetwork/ioeX.SPV/spvwallet/config"

	"github.com/urfave/cli"
)

var Version string

func init() {
	log.Init()
}

func main() {
	app := cli.NewApp()
	app.Name = "ELASTOS SPV WALLET"
	app.Version = Version
	app.HelpName = "ELASTOS SPV WALLET HELP"
	app.Usage = "command line user interface"
	app.UsageText = "[global option] command [command options] [args]"
	app.HideHelp = false
	app.HideVersion = false
	//commands
	app.Commands = []cli.Command{
		wallet.NewCreateCommand(),
		wallet.NewChangePasswordCommand(),
		wallet.NewResetCommand(),
		account.NewCommand(),
		transaction.NewCommand(),
	}

	app.Run(os.Args)
}
