package krakend

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	cmd "github.com/krakend/krakend-cobra/v2"
	"github.com/luraproject/lura/v2/logging"
	proxy "github.com/luraproject/lura/v2/proxy/plugin"
	client "github.com/luraproject/lura/v2/transport/http/client/plugin"
	server "github.com/luraproject/lura/v2/transport/http/server/plugin"
	"github.com/spf13/cobra"
)

// LoadPlugins loads and registers the plugins so they can be used if enabled at the configuration
func LoadPlugins(folder, pattern string, logger logging.Logger) {
	LoadPluginsWithContext(context.Background(), folder, pattern, logger)
}

func LoadPluginsWithContext(ctx context.Context, folder, pattern string, logger logging.Logger) {
	logger.Debug("[SERVICE: Plugin Loader] Starting loading process")

	n, err := client.LoadWithLogger(
		folder,
		pattern,
		client.RegisterClient,
		logger,
	)
	logPluginLoaderErrors(logger, "[SERVICE: Executor Plugin]", n, err)

	n, err = server.LoadWithLogger(
		folder,
		pattern,
		server.RegisterHandler,
		logger,
	)
	logPluginLoaderErrors(logger, "[SERVICE: Handler Plugin]", n, err)

	n, err = proxy.LoadWithLoggerAndContext(
		ctx,
		folder,
		pattern,
		proxy.RegisterModifier,
		logger,
	)
	logPluginLoaderErrors(logger, "[SERVICE: Modifier Plugin]", n, err)

	logger.Debug("[SERVICE: Plugin Loader] Loading process completed")
}

func logPluginLoaderErrors(logger logging.Logger, tag string, n int, err error) {
	if err != nil {
		if mErrs, ok := err.(pluginLoaderErr); ok {
			for _, err := range mErrs.Errs() {
				logger.Debug(tag, err.Error())
			}
		} else {
			logger.Debug(tag, err.Error())
		}
	}
	if n > 0 {
		logger.Info(tag, "Total plugins loaded:", n)
	}
}

type pluginLoader struct{}

func (pluginLoader) Load(folder, pattern string, logger logging.Logger) {
	LoadPlugins(folder, pattern, logger)
}

func (pluginLoader) LoadWithContext(ctx context.Context, folder, pattern string, logger logging.Logger) {
	LoadPluginsWithContext(ctx, folder, pattern, logger)
}

type pluginLoaderErr interface {
	Errs() []error
}

var (
	serverExpected   bool
	clientExpected   bool
	modifierExpected bool

	testPluginCmd = &cobra.Command{
		Use:     "test-plugin [flags] [artifacts]",
		Short:   "Tests that one or more plugins are loadable into KrakenD.",
		Run:     testPluginFunc,
		Example: "krakend test-plugin -scm ./plugins/my_plugin.so ./plugins/my_other_plugin.so",
	}

	serverExpectedFlag   cmd.FlagBuilder
	clientExpectedFlag   cmd.FlagBuilder
	modifierExpectedFlag cmd.FlagBuilder

	reLogErrorPlugins = regexp.MustCompile(`(?m)plugin \#\d+ \(.*\): (.*)`)
)

func init() {
	serverExpectedFlag = cmd.BoolFlagBuilder(&serverExpected, "server", "s", false, "The artifact should contain a Server Plugin.")
	clientExpectedFlag = cmd.BoolFlagBuilder(&clientExpected, "client", "c", false, "The artifact should contain a Client Plugin.")
	modifierExpectedFlag = cmd.BoolFlagBuilder(&modifierExpected, "modifier", "m", false, "The artifact should contain a Req/Resp Modifier Plugin.")
}

func NewTestPluginCmd() cmd.Command {
	return cmd.NewCommand(testPluginCmd, serverExpectedFlag, clientExpectedFlag, modifierExpectedFlag)
}

func testPluginFunc(ccmd *cobra.Command, args []string) {
	if len(args) == 0 {
		ccmd.Println("At least one plugin is required.")
		os.Exit(1)
	}
	if !serverExpected && !clientExpected && !modifierExpected {
		ccmd.Println("You must declare the expected type of the plugin.")
		os.Exit(1)
	}

	start := time.Now()

	ctx, cancel := context.WithCancel(context.Background())

	var failed int
	globalOK := true
	for _, pluginPath := range args {
		f, err := os.Open(pluginPath)
		if os.IsNotExist(err) {
			ccmd.Println(fmt.Sprintf("[KO] Unable to open the plugin %s.", pluginPath))
			failed++
			globalOK = false
			continue
		}
		f.Close()

		name := filepath.Base(pluginPath)
		folder := filepath.Dir(pluginPath)
		ok := true

		if serverExpected {
			ok = checkHandlerPlugin(ccmd, folder, name) && ok
		}

		if modifierExpected {
			ok = checkModifierPlugin(ctx, ccmd, folder, name) && ok
		}

		if clientExpected {
			ok = checkClientPlugin(ccmd, folder, name) && ok
		}

		if !ok {
			failed++
		}

		globalOK = globalOK && ok
	}

	cancel()

	if !globalOK {
		ccmd.Println(fmt.Sprintf("[KO] %d tested plugin(s) in %s.\n%d plugin(s) failed.", len(args), time.Since(start), failed))
		os.Exit(1)
	}

	ccmd.Println(fmt.Sprintf("[OK] %d tested plugin(s) in %s", len(args), time.Since(start)))
}

func checkClientPlugin(ccmd *cobra.Command, folder, name string) bool {
	_, err := client.LoadWithLogger(
		folder,
		name,
		client.RegisterClient,
		logging.NoOp,
	)
	if err == nil {
		ccmd.Println(fmt.Sprintf("[OK] CLIENT\t%s", name))
		return true
	}

	var msg string
	if mErrs, ok := err.(pluginLoaderErr); ok {
		for _, err := range mErrs.Errs() {
			msg += err.Error()
		}
	} else {
		msg = err.Error()
	}

	if strings.Contains(msg, "symbol ClientRegisterer not found") {
		ccmd.Println(fmt.Sprintf("[KO] CLIENT\t%s: The plugin does not contain a ClientRegisterer.", name))
		return false
	}

	for _, match := range reLogErrorPlugins.FindAllStringSubmatch(msg, -1) {
		msg = match[1]
	}

	ccmd.Println(fmt.Sprintf("[KO] CLIENT\t%s: %s", name, msg))
	return false
}

func checkHandlerPlugin(ccmd *cobra.Command, folder, name string) bool {
	_, err := server.LoadWithLogger(
		folder,
		name,
		server.RegisterHandler,
		logging.NoOp,
	)
	if err == nil {
		ccmd.Println(fmt.Sprintf("[OK] SERVER\t%s", name))
		return true
	}

	var msg string
	if mErrs, ok := err.(pluginLoaderErr); ok {
		for _, err := range mErrs.Errs() {
			msg += err.Error()
		}
	} else {
		msg = err.Error()
	}

	if strings.Contains(msg, "symbol HandlerRegisterer not found") {
		ccmd.Println(fmt.Sprintf("[KO] SERVER\t%s: The plugin does not contain a HandlerRegisterer.", name))
		return false
	}

	for _, match := range reLogErrorPlugins.FindAllStringSubmatch(msg, -1) {
		msg = match[1]
	}

	ccmd.Println(fmt.Sprintf("[KO] SERVER\t%s: %s", name, msg))
	return false
}

func checkModifierPlugin(ctx context.Context, ccmd *cobra.Command, folder, name string) bool {
	_, err := proxy.LoadWithLoggerAndContext(
		ctx,
		folder,
		name,
		proxy.RegisterModifier,
		logging.NoOp,
	)
	if err == nil {
		ccmd.Println(fmt.Sprintf("[OK] MODIFIER\t%s", name))
		return true
	}

	var msg string
	if mErrs, ok := err.(pluginLoaderErr); ok {
		for _, err := range mErrs.Errs() {
			msg += err.Error()
		}
	} else {
		msg = err.Error()
	}

	if strings.Contains(msg, "symbol ModifierRegisterer not found") {
		ccmd.Println(fmt.Sprintf("[KO] MODIFIER\t%s: The plugin does not contain a ModifierRegisterer.", name))
		return false
	}

	for _, match := range reLogErrorPlugins.FindAllStringSubmatch(msg, -1) {
		msg = match[1]
	}

	ccmd.Println(fmt.Sprintf("[KO] MODIFIER\t%s: %s", name, msg))
	return false
}
