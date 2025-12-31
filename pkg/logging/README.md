# krakend-gologging

An improved logger for the KrakenD framework

## How to use it

Import the package

	import "github.com/krakend/krakend-gologging"

After parsing the config file (here `cfg`), call the logger factory with the extra config of the system

	logger, _ := logging.NewLogger(cfg.ExtraConfig)

And the logger is ready to be injected

## Configuration

Add the `github_com/devopsfaith/krakend-gologging` section to the service extra config.

Example:

	"extra": {
		"github_com/devopsfaith/krakend-gologging": {
			"level":  "INFO",
			"prefix": "[KRAKEND]",
			"syslog": false,
			"stdout": true
		}
	}

1. `level`: name of the min log level to display
2. `prefix`: prefix to use with the application logs
3. `syslog`: enable logging over syslog
4. `stdout`: enable logging over the stdout
