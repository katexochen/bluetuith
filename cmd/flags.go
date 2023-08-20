package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/darkhz/bluetuith/bluez"
	"github.com/darkhz/bluetuith/theme"
	"github.com/knadh/koanf/parsers/hjson"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/providers/posflag"
	flag "github.com/spf13/pflag"
)

// Option describes a command-line option.
type Option struct {
	Name, Description, Value string
	IsBoolean                bool
}

var options = []Option{
	{
		Name:        "list-adapters",
		Description: "List available adapters.",
		IsBoolean:   true,
	},
	{
		Name:        "adapter",
		Description: "Specify an adapter to use. (For example, hci0)",
	},
	{
		Name:        "receive-dir",
		Description: "Specify a directory to store received files.",
	},
	{
		Name:        "gsm-apn",
		Description: "Specify GSM APN to connect to. (Required for DUN)",
	},
	{
		Name:        "gsm-number",
		Description: "Specify GSM number to dial. (Required for DUN)",
	},
	{
		Name:        "set-theme",
		Description: "Specify a theme." + theme.GetThemes(),
	},
}

func parse() {
	configFile, err := ConfigPath("bluetuith.conf")
	if err != nil {
		PrintError("Cannot get config directory")
	}

	fs := flag.NewFlagSet("bluetuith", flag.ContinueOnError)
	fs.Usage = func() {
		var usage string

		usage += fmt.Sprintf(
			"bluetuith [<flags>]\n\nConfig file is %s\n\nFlags:\n",
			configFile,
		)

		fs.VisitAll(func(f *flag.Flag) {
			s := fmt.Sprintf("  --%s", f.Name)

			switch f.Name {
			case "adapter":
				s += " <adapter>"

			case "receive-dir":
				s += " <dir>"

			case "gsm-apn":
				s += " <apn>"

			case "gsm-number":
				s += " <number>"

			case "set-theme":
				s += " <theme>"
			}

			if len(s) <= 4 {
				s += "\t"
			} else {
				s += "\n    \t"
			}

			s += strings.ReplaceAll(f.Usage, "\n", "\n    \t")

			usage += s + "\n"
		})

		Print(usage, 0)
	}

	for _, option := range options {
		if option.IsBoolean {
			fs.Bool(option.Name, false, option.Description)
			continue
		}

		fs.String(option.Name, option.Value, option.Description)
	}

	if err = fs.Parse(os.Args[1:]); err != nil {
		PrintError(err.Error())
	}

	if err := config.Load(file.Provider(configFile), hjson.Parser()); err != nil {
		PrintError(err.Error())
	}

	if err := config.Load(posflag.Provider(fs, ".", config.Koanf), nil); err != nil {
		PrintError(err.Error())
	}
}

func cmdOptionAdapter(b *bluez.Bluez) {
	optionAdapter := GetProperty("adapter")
	if optionAdapter == "" {
		b.SetCurrentAdapter()
		return
	}

	for _, adapter := range b.GetAdapters() {
		if optionAdapter == filepath.Base(adapter.Path) {
			b.SetCurrentAdapter(adapter)
			return
		}
	}

	PrintError(optionAdapter + ": The adapter does not exist.")
}

func cmdOptionListAdapters(b *bluez.Bluez) {
	var adapters string

	if !IsPropertyEnabled("list-adapters") {
		return
	}

	adapters += "List of adapters:\n"
	for _, adapter := range b.GetAdapters() {
		adapters += "- " + filepath.Base(adapter.Path) + "\n"
	}

	Print(adapters, 0)
}

func cmdOptionReceiveDir() {
	optionReceiveDir := GetProperty("receive-dir")
	if optionReceiveDir == "" {
		return
	}

	if statpath, err := os.Stat(optionReceiveDir); err == nil && statpath.IsDir() {
		AddProperty("receive-dir", optionReceiveDir)
		return
	}

	PrintError(optionReceiveDir + ": Directory is not accessible.")
}

func cmdOptionGsm() {
	optionGsmNumber := GetProperty("gsm-number")
	optionGsmApn := GetProperty("gsm-apn")

	if optionGsmNumber == "" && optionGsmApn != "" {
		PrintError("Specify GSM Number.")
	}

	number := "*99#"
	if optionGsmNumber != "" {
		number = optionGsmNumber
	}

	AddProperty("gsm-apn", optionGsmApn)
	AddProperty("gsm-number", number)
}

func cmdOptionTheme() {
	optionTheme := GetProperty("set-theme")
	if optionTheme == "" {
		return
	}

	if err := theme.ParseThemeFile(optionTheme); err != nil {
		PrintError(err.Error())
	}
}
