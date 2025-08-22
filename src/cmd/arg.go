package cmd

import (
	"fmt"
	"reflect"
	"strings"
)

/*
this command do things :
take an argument and display it
*/
// type CommandArgs struct {
// 	FieldI    int     `clap:"short:i,long:integer" doc:"get an integer as an argument"`
// 	FieldS    string  `clap:"short:s,long:string" doc:"get a string as an argument"`
// 	FieldF    float32 `clap:"short:f,long:float" doc:"get a float as an argument"`
// 	FieldVoid uint
// }

type arg struct {
	field     *reflect.StructField
	fieldName string
	doc       string
	short     string
	long      string
}

func Parse[T any](i *T) error {
	var numOfFields = reflect.TypeOf(*i).NumField()
	var typeOf = reflect.TypeOf(*i)
	var args []arg
	for i := range numOfFields {
		var field = typeOf.Field(i)
		var doc string
		if tag_v := field.Tag.Get("doc"); tag_v != "" {
			doc = tag_v
		}
		if tag := field.Tag.Get("clap"); tag != "" {
			arg, err := newArg(field.Name, tag, doc)
			if err != nil {
				return err
			}
			arg.field = &field
			args = append(args, *arg)
		}
	}
	// fmt.Println(args)
	// ParseArgs(args)
	return nil
}

func newArg(field_name string, opts string, doc string) (*arg, error) {
	var arg arg
	for part := range strings.SplitSeq(opts, ",") {
		if option, value, has_value := strings.Cut(part, ":"); has_value {
			switch option {
			case "short":
				arg.short = value
			case "long":
				arg.long = value
			default:
				return nil, fmt.Errorf("tag field \"%s\" does not exist", option)
			}
			if arg.long == "" {
				arg.long = field_name
			}
		} else {
			switch option {
			case "short":
				arg.short = string(field_name[0])
			case "long":
				arg.long = field_name
			default:
				return nil, fmt.Errorf("tag field \"%s\" does not exist", option)
			}
		}
	}
	if arg.short == "" && arg.long == "" {
		return nil, fmt.Errorf("tags short and long can't be both empty")
	}
	arg.doc = doc
	arg.fieldName = field_name
	return &arg, nil
}
