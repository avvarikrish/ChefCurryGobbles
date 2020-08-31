package input

import (
	"bufio"
	"fmt"
	"os"
	"reflect"
	"strings"
)

func ReadInput(obj reflect.Type, o reflect.Value) {
	reader := bufio.NewReader(os.Stdin)

	for i := 0; i < obj.NumField(); i++ {
		f := obj.Field(i)
		ft := o.Field(i)
		switch f.Type.Kind() {
		case reflect.Slice:
			fmt.Print("How many ", f.Name, "? ")
			var n int
			fmt.Scan(&n)
			s := reflect.MakeSlice(f.Type, n, n)
			sValue := reflect.New(s.Type())
			p := reflect.ValueOf(sValue.Interface()).Elem()

			for j := 0; j < n; j++ {
				t := s.Type().Elem()
				v := reflect.New(t)
				ReadInput(t, v.Elem())
				p.Set(reflect.Append(p, reflect.ValueOf(v.Interface()).Elem()))
			}
			ft.Set(p)
		case reflect.Struct:
			fmt.Print(f.Name, ": \n")
			ReadInput(f.Type, ft)
		default:
			fmt.Print(f.Name, ": ")
			v, _ := reader.ReadString('\n')
			ft.SetString(strings.TrimSpace(v))
		}
	}
}
