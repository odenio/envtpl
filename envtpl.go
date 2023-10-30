package envtpl

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"github.com/spf13/cobra"
)

// options
var (
	missingKey = "default"
	output     string
	version    bool
)

func die(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func Template(cmd *cobra.Command, args []string) {

	if cmd.Flag("version").Value.String() == "true" {
		fmt.Println(Version())
		os.Exit(0)
	}
	missingKey = cmd.Flag("missingkey").Value.String()
	output = cmd.Flag("output").Value.String()

	var t *template.Template
	var err error

	// load template; if an argument is not specified, default to stdin
	if len(args) > 0 {
		t, err = parseFiles(args...)
		die(err)
	} else {
		bytes, err := ioutil.ReadAll(os.Stdin)
		die(err)
		t, err = parse(string(bytes))
		die(err)
	}

	// get environment variables to supply to the template
	env := readEnv()

	// get writer for rendered output; if an output file is not
	// specified, default to stdout
	var w io.Writer
	if output == "" {
		w = os.Stdout
	} else {
		f, err := os.Create(output)
		die(err)
		defer f.Close()
		w = io.Writer(f)
	}

	// set error handling strategy for missing keys
	if missingKey != "default" {
		t = t.Option("missingkey=" + missingKey)
	}

	// render the template
	err = t.Execute(w, env)
	die(err)
}

// returns a new template with custom function maps
// (sprig and environment key prefix matcher) applied
func parse(s string) (*template.Template, error) {
	return template.New("").Funcs(sprig.TxtFuncMap()).Funcs(customFuncMap()).Parse(s)
}

// returns a new template with custom function maps
// (sprig and environment key prefix matcher) applied
func parseFiles(files ...string) (*template.Template, error) {
	return template.New(filepath.Base(files[0])).Funcs(sprig.TxtFuncMap()).Funcs(customFuncMap()).ParseFiles(files...)
}

// returns map of environment variables
func readEnv() (env map[string]string) {
	env = make(map[string]string)
	for _, setting := range os.Environ() {
		pair := strings.SplitN(setting, "=", 2)
		env[pair[0]] = pair[1]
	}
	return
}

// custom function that returns key, value for all environment variable keys matching prefix
// (see original envtpl: https://pypi.org/project/envtpl/)
func environment(prefix string) map[string]string {
	env := make(map[string]string)
	for _, setting := range os.Environ() {
		pair := strings.SplitN(setting, "=", 2)
		if strings.HasPrefix(pair[0], prefix) {
			env[pair[0]] = pair[1]
		}
	}
	return env
}

// returns custom template functions map
func customFuncMap() template.FuncMap {
	var functionMap = map[string]interface{}{
		"environment": environment,
	}
	return template.FuncMap(functionMap)
}
