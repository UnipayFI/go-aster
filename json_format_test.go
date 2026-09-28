package aster

import (
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/UnipayFI/go-aster/v3/common"
)

// untaggedTimeFields are the time.Time fields the SDK fills in itself rather
// than decoding from JSON: kline rows arrive as positional arrays parsed by
// parseKline, and the user-data union copies EventTime from the decoded event.
// They carry no json tag, so marshaling them keeps the RFC 3339 default.
var untaggedTimeFields = map[string]bool{
	"futures.Kline.OpenTime":            true,
	"futures.Kline.CloseTime":           true,
	"futures.WsUserDataEvent.EventTime": true,
	"spot.Kline.OpenTime":               true,
	"spot.Kline.CloseTime":              true,
	"spot.WsUserDataEvent.EventTime":    true,
}

// TestTimeFieldsDeclareFormat requires every exported time.Time / *time.Time
// struct field with a json tag to declare its wire format with the `format`
// tag option (e.g. `json:"updateTime,format:unixmilli"`). Without one,
// encoding/json/v2 falls back to RFC 3339, which no Aster timestamp uses, so a
// missing option would only surface as a decode error against the live API.
// The option must also be one the standard library accepts, since a misspelt
// one fails every decode of the struct. A field without any json tag must be
// listed in untaggedTimeFields.
func TestTimeFieldsDeclareFormat(t *testing.T) {
	fset := token.NewFileSet()
	checked := 0
	seenUntagged := map[string]bool{}
	err := filepath.WalkDir(".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			if name := d.Name(); path != "." && (strings.HasPrefix(name, ".") || name == "testdata") {
				return filepath.SkipDir
			}
			return nil
		}
		if !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
			return nil
		}
		f, err := parser.ParseFile(fset, path, nil, parser.SkipObjectResolution)
		if err != nil {
			return err
		}
		// Name every struct type after the declared type that contains it;
		// struct literals inside function bodies stay unnamed.
		owner := map[*ast.StructType]string{}
		ast.Inspect(f, func(n ast.Node) bool {
			if ts, ok := n.(*ast.TypeSpec); ok {
				ast.Inspect(ts.Type, func(n ast.Node) bool {
					if st, ok := n.(*ast.StructType); ok {
						owner[st] = ts.Name.Name
					}
					return true
				})
			}
			return true
		})
		pkg := filepath.ToSlash(filepath.Dir(path))
		ast.Inspect(f, func(n ast.Node) bool {
			st, ok := n.(*ast.StructType)
			if !ok {
				return true
			}
			for _, field := range st.Fields.List {
				if !isTimeType(field.Type) || len(field.Names) == 0 || !field.Names[0].IsExported() {
					continue
				}
				id := pkg + "." + owner[st] + "." + field.Names[0].Name
				if field.Tag == nil {
					if !untaggedTimeFields[id] {
						t.Errorf("%s: field %s has no json tag; add one with a format option", fset.Position(field.Pos()), id)
					}
					seenUntagged[id] = true
					continue
				}
				raw, _ := strconv.Unquote(field.Tag.Value)
				tag, ok := reflect.StructTag(raw).Lookup("json")
				if tag == "-" {
					continue
				}
				checked++
				opts := strings.Split(tag, ",")[1:]
				if !ok || len(opts) == 0 || !strings.HasPrefix(opts[len(opts)-1], "format:") {
					t.Errorf("%s: field %s has json tag %q without a trailing format option", fset.Position(field.Pos()), id, tag)
				} else if err := roundTripTimeTag(field.Type, tag); err != nil {
					t.Errorf("%s: field %s has json tag %q that the codec rejects: %v", fset.Position(field.Pos()), id, tag, err)
				}
			}
			return true
		})
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if checked == 0 {
		t.Fatal("no time.Time fields found; is the test running from the module root?")
	}
	for id := range untaggedTimeFields {
		if !seenUntagged[id] {
			t.Errorf("untaggedTimeFields lists %s, which no longer exists or now has a json tag", id)
		}
	}
}

// roundTripTimeTag encodes and decodes a sample time through the SDK codec in
// a one-field struct carrying the given json tag, so encoding/json/v2 itself
// validates the format option.
func roundTripTimeTag(typ ast.Expr, tag string) error {
	sample := reflect.ValueOf(time.Date(2026, 9, 28, 12, 34, 56, 789e6, time.UTC))
	if _, ok := typ.(*ast.StarExpr); ok {
		p := reflect.New(sample.Type())
		p.Elem().Set(sample)
		sample = p
	}
	st := reflect.StructOf([]reflect.StructField{{Name: "T", Type: sample.Type(), Tag: reflect.StructTag(`json:"` + tag + `"`)}})
	v := reflect.New(st)
	v.Elem().Field(0).Set(sample)
	b, err := common.JSONMarshal(v.Interface())
	if err != nil {
		return err
	}
	return common.JSONUnmarshal(b, v.Interface())
}

func isTimeType(e ast.Expr) bool {
	if star, ok := e.(*ast.StarExpr); ok {
		e = star.X
	}
	sel, ok := e.(*ast.SelectorExpr)
	if !ok {
		return false
	}
	pkg, ok := sel.X.(*ast.Ident)
	return ok && pkg.Name == "time" && sel.Sel.Name == "Time"
}
