package siesta

import (
	"flag"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// Params represents a set of URL parameters from a request's query string.
// The interface is similar to a flag.FlagSet, but a) there is no usage string,
// b) there are no custom Var()s, and c) there are SliceXXX types. Sliced types
// support two ways of generating a multi-valued parameter: setting the parameter
// multiple times, and using a comma-delimited string. This adds the limitation
// that you can't have a value with a comma if in a Sliced type.
// Under the covers, Params uses flag.FlagSet.
type Params struct {
	fset *flag.FlagSet
}

// Parse parses URL parameters from a http.Request.URL.Query(), which is a
// url.Values, which is just a map[string][string].
func (rp *Params) Parse(args url.Values) error {
	if rp.fset == nil {
		rp.fset = flag.NewFlagSet("anonymous", flag.ExitOnError) // both args are unused.
	}

	// Parse items from URL query string
FLAG_LOOP:
	for name, vals := range args {
		for _, v := range vals {

			f := rp.fset.Lookup(name)
			if f == nil {
				// Flag wasn't found.
				continue FLAG_LOOP
			}

			// Check if the value is empty
			if v == "" {
				if bv, ok := f.Value.(boolFlag); ok && bv.IsBoolFlag() {
					bv.Set("true")

					continue FLAG_LOOP
				}
			}

			err := rp.fset.Set(name, v)
			if err != nil {
				// Remove the "flag" error message and make a "params" one.
				// TODO: optionally allow undefined params to be given, but ignored?
				if !strings.Contains(err.Error(), "no such flag -") {
					// Give a helpful message about which param caused the error
					err = fmt.Errorf("bad param '%s': %s", name, err.Error())
					return err
				}
			}
		}
	}

	return nil
}

// Usage returns a map keyed on parameter names. The map values are an array of
// name, type, and usage information for each parameter.
func (rp *Params) Usage() map[string][3]string {
	docs := make(map[string][3]string)
	var translations map[string]string = map[string]string{
		"*flag.stringValue":   "string",
		"*flag.durationValue": "duration",
		"*flag.intValue":      "int",
		"*flag.boolValue":     "bool",
		"*flag.float64Value":  "float64",
		"*flag.int64Value":    "int64",
		"*flag.uintValue":     "uint",
		"*flag.uint64Value":   "uint64",
		"*siesta.SString":     "[]string",
		"*siesta.SDuration":   "[]duration",
		"*siesta.SInt":        "[]int",
		"*siesta.SBool":       "[]bool",
		"*siesta.SFloat64":    "[]float64",
		"*siesta.SInt64":      "[]int64",
		"*siesta.SUint":       "[]uint",
		"*siesta.SUint64":     "[]uint64",
	}
	rp.fset.VisitAll(func(flag *flag.Flag) {
		niceName := translations[fmt.Sprintf("%T", flag.Value)]
		if niceName == "" {
			niceName = fmt.Sprintf("%T", flag.Value)
		}
		docs[flag.Name] = [...]string{flag.Name, niceName, flag.Usage}
	})
	return docs
}

type boolFlag interface {
	flag.Value
	IsBoolFlag() bool
}

// Bool defines a bool param with specified name and default value.
// The return value is the address of a bool variable that stores the value of the param.
func (rp *Params) Bool(name string, value bool, usage string) *bool {
	if rp.fset == nil {
		rp.fset = flag.NewFlagSet("anonymous", flag.ExitOnError) // both args are unused.
	}
	p := new(bool)
	rp.fset.BoolVar(p, name, value, usage)
	return p
}

// SBool is a slice of bool.
type SBool []bool

// String is the method to format the param's value, part of the flag.Value interface.
// The String method's output will be used in diagnostics.
func (s *SBool) String() string {
	return fmt.Sprint(*s)
}

// Set is the method to set the param value, part of the flag.Value interface.
// Set's argument is a string to be parsed to set the param.
// It's a comma-separated list, so we split it.
func (s *SBool) Set(value string) error {
	for _, dt := range strings.Split(value, ",") {
		parsed, err := strconv.ParseBool(dt)
		if err != nil {
			return err
		}
		*s = append(*s, parsed)
	}
	return nil
}

// SliceBool defines a multi-value bool param with specified name and default value.
// The return value is the address of a SBool variable that stores the values of the param.
func (rp *Params) SliceBool(name string, value bool, usage string) *SBool {
	if rp.fset == nil {
		rp.fset = flag.NewFlagSet("anonymous", flag.ExitOnError) // both args are unused.
	}
	p := new(SBool)
	rp.fset.Var(p, name, usage)
	return p
}

// Int defines an int param with specified name and default value.
// The return value is the address of an int variable that stores the value of the param.
func (rp *Params) Int(name string, value int, usage string) *int {
	if rp.fset == nil {
		rp.fset = flag.NewFlagSet("anonymous", flag.ExitOnError) // both args are unused.
	}
	p := new(int)
	rp.fset.IntVar(p, name, value, usage)
	return p
}

// SInt is a slice of int.
type SInt []int

// String is the method to format the param's value, part of the flag.Value interface.
// The String method's output will be used in diagnostics.
func (s *SInt) String() string {
	return fmt.Sprint(*s)
}

// Set is the method to set the param value, part of the flag.Value interface.
// Set's argument is a string to be parsed to set the param.
// It's a comma-separated list, so we split it.
func (s *SInt) Set(value string) error {
	for _, dt := range strings.Split(value, ",") {
		parsed, err := strconv.ParseInt(dt, 0, 64)
		if err != nil {
			return err
		}
		*s = append(*s, int(parsed))
	}
	return nil
}

// SliceInt defines a multi-value int param with specified name and default value.
// The return value is the address of a SInt variable that stores the values of the param.
func (rp *Params) SliceInt(name string, value int, usage string) *SInt {
	if rp.fset == nil {
		rp.fset = flag.NewFlagSet("anonymous", flag.ExitOnError) // both args are unused.
	}
	p := new(SInt)
	rp.fset.Var(p, name, usage)
	return p
}

// Int64 defines an int64 param with specified name and default value.
// The return value is the address of an int64 variable that stores the value of the param.
func (rp *Params) Int64(name string, value int64, usage string) *int64 {
	if rp.fset == nil {
		rp.fset = flag.NewFlagSet("anonymous", flag.ExitOnError) // both args are unused.
	}
	p := new(int64)
	rp.fset.Int64Var(p, name, value, usage)
	return p
}

// SInt64 is a slice of int64.
type SInt64 []int64

// String is the method to format the param's value, part of the flag.Value interface.
// The String method's output will be used in diagnostics.
func (s *SInt64) String() string {
	return fmt.Sprint(*s)
}

// Set is the method to set the param value, part of the flag.Value interface.
// Set's argument is a string to be parsed to set the param.
// It's a comma-separated list, so we split it.
func (s *SInt64) Set(value string) error {
	for _, dt := range strings.Split(value, ",") {
		parsed, err := strconv.ParseInt(dt, 0, 64)
		if err != nil {
			return err
		}
		*s = append(*s, int64(parsed))
	}
	return nil
}

// SliceInt64 defines a multi-value int64 param with specified name and default value.
// The return value is the address of a SInt64 variable that stores the values of the param.
func (rp *Params) SliceInt64(name string, value int64, usage string) *SInt64 {
	if rp.fset == nil {
		rp.fset = flag.NewFlagSet("anonymous", flag.ExitOnError) // both args are unused.
	}
	p := new(SInt64)
	rp.fset.Var(p, name, usage)
	return p
}

// Uint defines a uint param with specified name and default value.
// The return value is the address of a uint variable that stores the value of the param.
func (rp *Params) Uint(name string, value uint, usage string) *uint {
	if rp.fset == nil {
		rp.fset = flag.NewFlagSet("anonymous", flag.ExitOnError) // both args are unused.
	}
	p := new(uint)
	rp.fset.UintVar(p, name, value, usage)
	return p
}

// SUint is a slice of uint.
type SUint []uint

// String is the method to format the param's value, part of the flag.Value interface.
// The String method's output will be used in diagnostics.
func (s *SUint) String() string {
	return fmt.Sprint(*s)
}

// Set is the method to set the param value, part of the flag.Value interface.
// Set's argument is a string to be parsed to set the param.
// It's a comma-separated list, so we split it.
func (s *SUint) Set(value string) error {
	for _, dt := range strings.Split(value, ",") {
		parsed, err := strconv.ParseFloat(dt, 64)
		if err != nil {
			return err
		}
		*s = append(*s, uint(parsed))
	}
	return nil
}

// SliceUint defines a multi-value uint param with specified name and default value.
// The return value is the address of a SUint variable that stores the values of the param.
func (rp *Params) SliceUint(name string, value uint, usage string) *SUint {
	if rp.fset == nil {
		rp.fset = flag.NewFlagSet("anonymous", flag.ExitOnError) // both args are unused.
	}
	p := new(SUint)
	rp.fset.Var(p, name, usage)
	return p
}

// Uint64 defines a uint64 param with specified name and default value.
// The return value is the address of a uint64 variable that stores the value of the param.
func (rp *Params) Uint64(name string, value uint64, usage string) *uint64 {
	if rp.fset == nil {
		rp.fset = flag.NewFlagSet("anonymous", flag.ExitOnError) // both args are unused.
	}
	p := new(uint64)
	rp.fset.Uint64Var(p, name, value, usage)
	return p
}

// SUint64 is a slice of uint64.
type SUint64 []uint64

// String is the method to format the param's value, part of the flag.Value interface.
// The String method's output will be used in diagnostics.
func (s *SUint64) String() string {
	return fmt.Sprint(*s)
}

// Set is the method to set the param value, part of the flag.Value interface.
// Set's argument is a string to be parsed to set the param.
// It's a comma-separated list, so we split it.
func (s *SUint64) Set(value string) error {
	for _, dt := range strings.Split(value, ",") {
		parsed, err := strconv.ParseFloat(dt, 64)
		if err != nil {
			return err
		}
		*s = append(*s, uint64(parsed))
	}
	return nil
}

// SliceUint64 defines a multi-value uint64 param with specified name and default value.
// The return value is the address of a SUint64 variable that stores the values of the param.
func (rp *Params) SliceUint64(name string, value uint64, usage string) *SUint64 {
	if rp.fset == nil {
		rp.fset = flag.NewFlagSet("anonymous", flag.ExitOnError) // both args are unused.
	}
	p := new(SUint64)
	rp.fset.Var(p, name, usage)
	return p
}

// String defines a string param with specified name and default value.
// The return value is the address of a string variable that stores the value of the param.
func (rp *Params) String(name string, value string, usage string) *string {
	if rp.fset == nil {
		rp.fset = flag.NewFlagSet("anonymous", flag.ExitOnError) // both args are unused.
	}
	p := new(string)
	rp.fset.StringVar(p, name, value, usage)
	return p
}

// SString is a slice of string.
type SString []string

// String is the method to format the param's value, part of the flag.Value interface.
// The String method's output will be used in diagnostics.
func (s *SString) String() string {
	return strings.Join(*s, ",")
}

// Set is the method to set the param value, part of the flag.Value interface.
// Set's argument is a string to be parsed to set the param.
// It's a comma-separated list, so we split it.
func (s *SString) Set(value string) error {
	for _, dt := range strings.Split(value, ",") {
		*s = append(*s, dt)
	}
	return nil
}

// SliceString defines a multi-value string param with specified name and default value.
// The return value is the address of a SString variable that stores the values of the param.
func (rp *Params) SliceString(name string, value string, usage string) *SString {
	if rp.fset == nil {
		rp.fset = flag.NewFlagSet("anonymous", flag.ExitOnError) // both args are unused.
	}
	p := new(SString)
	rp.fset.Var(p, name, usage)
	return p
}

// Float64 defines a float64 param with specified name and default value.
// The return value is the address of a float64 variable that stores the value of the param.
func (rp *Params) Float64(name string, value float64, usage string) *float64 {
	if rp.fset == nil {
		rp.fset = flag.NewFlagSet("anonymous", flag.ExitOnError) // both args are unused.
	}
	p := new(float64)
	rp.fset.Float64Var(p, name, value, usage)
	return p
}

// SFloat64 is a slice of float64.
type SFloat64 []float64

// String is the method to format the param's value, part of the flag.Value interface.
// The String method's output will be used in diagnostics.
func (s *SFloat64) String() string {
	return fmt.Sprintf("%f", *s)
}

// Set is the method to set the param value, part of the flag.Value interface.
// Set's argument is a string to be parsed to set the param.
// It's a comma-separated list, so we split it.
func (s *SFloat64) Set(value string) error {
	for _, dt := range strings.Split(value, ",") {
		parsed, err := strconv.ParseFloat(dt, 64)
		if err != nil {
			return err
		}
		*s = append(*s, parsed)
	}
	return nil
}

// SliceFloat64 defines a multi-value float64 param with specified name and default value.
// The return value is the address of a SFloat64 variable that stores the values of the param.
func (rp *Params) SliceFloat64(name string, value float64, usage string) *SFloat64 {
	if rp.fset == nil {
		rp.fset = flag.NewFlagSet("anonymous", flag.ExitOnError) // both args are unused.
	}
	p := new(SFloat64)
	rp.fset.Var(p, name, usage)
	return p
}

// Duration defines a time.Duration param with specified name and default value.
// The return value is the address of a time.Duration variable that stores the value of the param.
func (rp *Params) Duration(name string, value time.Duration, usage string) *time.Duration {
	if rp.fset == nil {
		rp.fset = flag.NewFlagSet("anonymous", flag.ExitOnError) // both args are unused.
	}
	p := new(time.Duration)
	rp.fset.DurationVar(p, name, value, usage)
	return p
}

// SDuration is a slice of time.Duration.
type SDuration []time.Duration

// String is the method to format the param's value, part of the flag.Value interface.
// The String method's output will be used in diagnostics.
func (s *SDuration) String() string {
	return fmt.Sprint(*s)
}

// Set is the method to set the param value, part of the flag.Value interface.
// Set's argument is a string to be parsed to set the param.
// It's a comma-separated list, so we split it.
func (s *SDuration) Set(value string) error {
	for _, dt := range strings.Split(value, ",") {
		parsed, err := time.ParseDuration(dt)
		if err != nil {
			return err
		}
		*s = append(*s, parsed)
	}
	return nil
}

// SliceDuration defines a multi-value time.Duration param with specified name and default value.
// The return value is the address of a SDuration variable that stores the values of the param.
func (rp *Params) SliceDuration(name string, value time.Duration, usage string) *SDuration {
	if rp.fset == nil {
		rp.fset = flag.NewFlagSet("anonymous", flag.ExitOnError) // both args are unused.
	}
	p := new(SDuration)
	rp.fset.Var(p, name, usage)
	return p
}
