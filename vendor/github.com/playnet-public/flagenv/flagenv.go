package flagenv

import (
	"flag"
	"os"
	"strconv"
	"strings"

	"time"

	"github.com/golang/glog"
)

//EnvPrefix is the prefix for environment variables
var EnvPrefix string

//String returns a flagenv variable pointer of type string
func String(name string, value string, usage string) *string {
	return flag.String(name, envString(parameterNameToEnvName(name), value), usage)
}

//StringVar sets a given pointer flagenv variable of type string
func StringVar(p *string, name string, value string, usage string) {
	flag.StringVar(p, name, envString(parameterNameToEnvName(name), value), usage)
}

//Duration returns a flagenv variable pointer of type time.Duration
func Duration(name string, value time.Duration, usage string) *time.Duration {
	return flag.Duration(name, envDuration(parameterNameToEnvName(name), value), usage)
}

//DurationVar sets a given pointer flagenv variable of type time.Duration
func DurationVar(p *time.Duration, name string, value time.Duration, usage string) {
	flag.DurationVar(p, name, envDuration(parameterNameToEnvName(name), value), usage)
}

//Int returns a flagenv variable pointer of type int
func Int(name string, value int, usage string) *int {
	return flag.Int(name, envInt(parameterNameToEnvName(name), value), usage)
}

//IntVar sets a given pointer flagenv variable of type int
func IntVar(p *int, name string, value int, usage string) {
	flag.IntVar(p, name, envInt(parameterNameToEnvName(name), value), usage)
}

//Bool returns a flagenv variable pointer of type bool
func Bool(name string, value bool, usage string) *bool {
	return flag.Bool(name, envBool(parameterNameToEnvName(name), value), usage)
}

//BoolVar sets a given pointer flagenv variable of type bool
func BoolVar(p *bool, name string, value bool, usage string) {
	flag.BoolVar(p, name, envBool(parameterNameToEnvName(name), value), usage)
}

func envString(key, def string) string {
	if env := os.Getenv(key); env != "" {
		return env
	}
	return def
}

func envBool(key string, def bool) bool {
	if env := os.Getenv(key); env != "" {
		res, err := strconv.ParseBool(env)
		if err != nil {
			return def
		}

		return res
	}
	return def
}

func envDuration(key string, def time.Duration) time.Duration {
	if env := os.Getenv(key); env != "" {
		val, err := time.ParseDuration(env)
		if err != nil {
			glog.V(2).Infof("invalid value for %q: using default: %q", key, def)
			return def
		}
		return val
	}
	return def
}

func envInt(key string, def int) int {
	if env := os.Getenv(key); env != "" {
		val, err := strconv.Atoi(env)
		if err != nil {
			glog.V(2).Infof("invalid value for %q: using default: %q", key, def)
			return def
		}
		return val
	}
	return def
}

func parameterNameToEnvName(name string) string {
	return EnvPrefix + strings.Replace(strings.ToUpper(name), "-", "_", -1)
}

//Parse is an alias to flag.Parse()
func Parse() {
	flag.Parse()
}
