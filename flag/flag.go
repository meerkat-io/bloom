package flag

import (
	"fmt"
	"os"
	"path"
	"reflect"
	"strconv"
	"unsafe"
)

/**********************************************************
type Flags struct {
	Help    bool   `flag:"h" usage:"help"`
	Lang    string `flag:"t" usage:"language template" tip:"go|cs" required:"true"`
	Input   string `flag:"i" usage:"input folder" tip:"input" default:"."`
	Output  string `flag:"o" usage:"output folder" tip:"output" default:"."`
}
// Supported types: int, string, float64, bool
***********************************************************/

// Parse command line arguments
func ParseCommandLine(config interface{}) error {
	flags, err := parseFlags(config)
	if err != nil {
		return nil
	}
	return parseArgs(flags)
}

func ParseEnvironment(config interface{}) error {
	flags, err := parseFlags(config)
	if err != nil {
		return nil
	}
	return parseEnv(flags)
}

// Usage print usage
func Usage() {
	file := path.Base(os.Args[0])
	fmt.Printf("\nUsage: %s ", file)
	for _, f := range flags {
		fmt.Printf("[-%s %s]", f.name, f.tip)
	}
	fmt.Print("\n\nOptions:\n")
	for _, f := range flags {
		fmt.Printf("    -%s:    %s\n", f.name, f.usage)
	}
	fmt.Println("")
}

var flags map[string]*flag

type flag struct {
	name         string // name as it appears on command line
	usage        string // help message
	tip          string // short help message
	defaultValue string // default value (as text); for usage message
	required     bool   // if required value

	value   value
	visited bool
}

type value interface {
	get() interface{}
	set(string) error
	isBool() bool
}

type boolValue bool

func newBoolValue(val bool, p *bool) *boolValue {
	*p = val
	return (*boolValue)(p)
}

func (b *boolValue) get() interface{} { return bool(*b) }

func (b *boolValue) set(s string) error {
	v, err := strconv.ParseBool(s)
	*b = boolValue(v)
	return err
}

func (b *boolValue) isBool() bool { return true }

type stringValue string

func newStringValue(val string, p *string) *stringValue {
	*p = val
	return (*stringValue)(p)
}

func (s *stringValue) get() interface{} { return string(*s) }

func (s *stringValue) set(val string) error {
	*s = stringValue(val)
	return nil
}

func (s *stringValue) isBool() bool { return false }

type intValue int

func newIntValue(val int, p *int) *intValue {
	*p = val
	return (*intValue)(p)
}

func (i *intValue) get() interface{} { return int(*i) }

func (i *intValue) set(s string) error {
	v, err := strconv.ParseInt(s, 0, strconv.IntSize)
	*i = intValue(v)
	return err
}

func (i *intValue) isBool() bool { return false }

type floatValue float64

func newFloatValue(val float64, p *float64) *floatValue {
	*p = val
	return (*floatValue)(p)
}

func (f *floatValue) get() interface{} { return float64(*f) }

func (f *floatValue) set(s string) error {
	v, err := strconv.ParseFloat(s, 64)
	*f = floatValue(v)
	return err
}

func (f *floatValue) isBool() bool { return false }

func parseFlags(config interface{}) (map[string]*flag, error) {
	flags = make(map[string]*flag)

	info := reflect.TypeOf(config)
	if info.Kind() != reflect.Ptr {
		return nil, fmt.Errorf("must use pointer of config struct")
	}
	info = info.Elem()
	value := reflect.ValueOf(config)
	value = value.Elem()

	var err error
	num := info.NumField()
	for i := 0; i < num; i++ {
		f := &flag{}
		field := info.Field(i)

		f.name = field.Tag.Get("flag")
		if f.name == "" {
			return nil, fmt.Errorf("flag is empty: %s", field.Name)
		}
		f.defaultValue = field.Tag.Get("default")
		f.usage = field.Tag.Get("usage")
		f.tip = field.Tag.Get("tip")
		f.required = field.Tag.Get("required") == "true"

		switch value.Field(i).Interface().(type) {
		case bool:
			f.value = newBoolValue(f.defaultValue == "true", (*bool)(unsafe.Pointer(value.Field(i).Addr().Pointer())))
		case string:
			f.value = newStringValue(f.defaultValue, (*string)(unsafe.Pointer(value.Field(i).Addr().Pointer())))
		case int:
			var v int64
			if f.defaultValue != "" {
				v, err = strconv.ParseInt(f.defaultValue, 10, 32)
				if err != nil {
					return nil, fmt.Errorf("invalid flag default value: [name:%s], [value:%s]", field.Name, f.defaultValue)
				}
			}
			f.value = newIntValue(int(v), (*int)(unsafe.Pointer(value.Field(i).Addr().Pointer())))
		case float64:
			var v float64
			if f.defaultValue != "" {
				v, err = strconv.ParseFloat(f.defaultValue, 64)
				if err != nil {
					return nil, fmt.Errorf("invalid flag default value: [name:%s], [value:%s]", field.Name, f.defaultValue)
				}
			}
			f.value = newFloatValue(v, (*float64)(unsafe.Pointer(value.Field(i).Addr().Pointer())))
		default:
			return nil, fmt.Errorf("invalid flag type: [name:%s], [value:%s]", field.Name, field.Type.String())
		}
		flags[f.name] = f
	}
	return flags, nil
}

func parseArgs(flags map[string]*flag) error {
	args := os.Args[1:]
	for len(args) > 0 {
		s := args[0]
		if len(s) < 2 || s[0] != '-' {
			return fmt.Errorf("invalid flag %s", s)
		}
		name := s[1:]
		if name[0] == '-' {
			return fmt.Errorf("bad flag syntax: %s", s)
		}
		f, ok := flags[name]
		if !ok {
			return fmt.Errorf("flag provided but not defined: -%s", name)
		}
		f.visited = true
		args = args[1:]
		if f.value.isBool() {
			_ = f.value.set("true")
		} else {
			if len(args) == 0 || args[0][0] == '-' {
				return fmt.Errorf("flag -%s has no value", name)
			}
			err := f.value.set(args[0])
			if err != nil {
				return err
			}
			args = args[1:]
		}
	}

	for _, f := range flags {
		if !f.visited && f.required {
			return fmt.Errorf("flag -%s is required but not set", f.name)
		}
	}
	return nil
}

func parseEnv(flags map[string]*flag) error {
	for name, f := range flags {
		value, exists := os.LookupEnv(name)
		if !exists && f.required {
			return fmt.Errorf("flag -%s has no value", name)
		}
		if f.value.isBool() {
			if exists {
				_ = f.value.set("true")
			}
		} else {
			if value == "" {
				return fmt.Errorf("flag -%s has no value", name)
			}
			err := f.value.set(value)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
