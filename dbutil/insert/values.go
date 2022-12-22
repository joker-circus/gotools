package insert

import (
	"database/sql/driver"
	"fmt"
	"reflect"
	"strings"
	"time"
)

func FmtField(field interface{}) string {
	ts, ok := field.(time.Time)
	if ok {
		field = ts.Format(DatetimeFormat)
	}
	return fmt.Sprintf("'%v'", field)
}

type Writer interface {
	WriteByte(byte) error
	WriteString(string) (int, error)
}

type Builder interface {
	Writer
	WriteColumn(field interface{})
	WriteValues(Writer, ...interface{})
}

type write struct {
	strings.Builder
	err error
}

func (w *write) WriteColumn(field interface{}) {
	w.WriteString(fmt.Sprintf("`%v`", field))
}

func (w *write) WriteValues(writer Writer, vars ...interface{}) {

	for idx, v := range vars {
		if idx > 0 {
			writer.WriteByte(',')
		}

		rv := reflect.ValueOf(v)
		switch vv := v.(type) {
		// 正确可参考：https://github.com/denisenkom/go-mssqldb/blob/v0.12.2/mssql.go#L901
		case driver.Valuer:
			if rv.Kind() == reflect.Ptr && reflect.ValueOf(v).IsNil() {
				writer.WriteString("NULL")
				continue
			}

			value, err := vv.Value()
			w.Error(err)

			writer.WriteString(FmtField(value))

		case []interface{}:
			if len(vv) == 0 {
				writer.WriteString("(NULL)")
				continue
			}

			writer.WriteByte('(')
			w.WriteValues(writer, vv...)
			writer.WriteByte(')')

		default:
			switch rv.Kind() {
			case reflect.Slice, reflect.Array:
				if rv.Len() == 0 {
					writer.WriteString("(NULL)")
					continue
				}

				//if rv.Type().Elem() == reflect.TypeOf(uint8(0)) {
				//  writer.WriteString(Field(v))
				//  continue
				//}

				writer.WriteByte('(')
				for i := 0; i < rv.Len(); i++ {
					if i > 0 {
						writer.WriteByte(',')
					}
					w.WriteValues(writer, rv.Index(i).Interface())
				}
				writer.WriteByte(')')
			case reflect.Interface, reflect.Ptr:
				if rv.IsNil() {
					writer.WriteString("NULL")
					continue
				}

				writer.WriteString(FmtField(v))
			default:
				writer.WriteString(FmtField(v))
			}
		}
	}
}

func (w *write) Error(err error) error {
	if err == nil {
		return nil
	}

	if w.err == nil {
		w.err = err
	}

	return err
}

// Fork Values
type SQL struct {
	Table   string
	Columns []string
	Values  [][]interface{}
}

func (s SQL) BuildInsertSQL() (string, error) {
	var w write
	w.WriteString(fmt.Sprintf("INSERT INTO `%s` ", s.Table))
	s.buildInsertSQL(&w)
	if w.err != nil {
		return "", w.err
	}

	return w.String(), nil
}

// Build is build from clause
func (s SQL) buildInsertSQL(builder Builder) {
	if len(s.Columns) > 0 {
		builder.WriteByte('(')
		for idx, column := range s.Columns {
			if idx > 0 {
				builder.WriteByte(',')
			}
			builder.WriteColumn(column)
		}
		builder.WriteByte(')')

		builder.WriteString(" VALUES ")

		for idx, value := range s.Values {
			if idx > 0 {
				builder.WriteByte(',')
			}

			builder.WriteByte('(')
			builder.WriteValues(builder, value...)
			builder.WriteByte(')')
		}
	} else {
		builder.WriteString("DEFAULT VALUES")
	}
}

func (s SQL) ExplainSQL() (string, error) {
	preSQL, err := s.PreSQL()
	if err != nil {
		return "", err
	}

	var args []interface{}
	for _, v := range s.Values {
		args = append(args, v...)
	}

	return ExplainSQL(preSQL, nil, `'`, args...), nil
}
func (s SQL) PreSQL() (string, error) {
	var w write
	w.WriteString(fmt.Sprintf("INSERT INTO `%s` ", s.Table))
	s.buildPreSQL(&w)
	if w.err != nil {
		return "", w.err
	}

	return w.String(), nil
}

func (s SQL) buildPreSQL(builder Builder) {
	if len(s.Columns) > 0 {
		builder.WriteByte('(')
		for idx, column := range s.Columns {
			if idx > 0 {
				builder.WriteByte(',')
			}
			builder.WriteColumn(column)
		}
		builder.WriteByte(')')

		builder.WriteString(" VALUES ")

		for idx, value := range s.Values {
			if idx > 0 {
				builder.WriteByte(',')
			}

			builder.WriteByte('(')
			for j := range value {
				if j > 0 {
					builder.WriteByte(',')
				}
				builder.WriteByte('?')
			}
			builder.WriteByte(')')
		}
	} else {
		builder.WriteString("DEFAULT VALUES")
	}
}
