package plugin

import (
	"fmt"
	"reflect"
	"strings"

	"go.autokitteh.dev/sdk/api/apivalues"
)

func ConvertFromAWS(pathElems []string, src reflect.Value) (*apivalues.Value, error) {
	path := strings.Join(pathElems, ".")

	switch t := src.Type(); t.Kind() {
	case reflect.Ptr:
		if src.IsNil() {
			return apivalues.None, nil
		}

		return ConvertFromAWS(pathElems, src.Elem())
	case reflect.Bool:
		return apivalues.Boolean(src.Interface().(bool)), nil
	case reflect.Int:
		return apivalues.Integer(int64(src.Interface().(int))), nil
	case reflect.Int64:
		return apivalues.Integer(src.Interface().(int64)), nil
	case reflect.Int32:
		return apivalues.Integer(int64(src.Interface().(int32))), nil
	case reflect.Int16:
		return apivalues.Integer(int64(src.Interface().(int16))), nil
	case reflect.Int8:
		return apivalues.Integer(int64(src.Interface().(int8))), nil
	case reflect.String:
		return apivalues.String(src.Convert(reflect.TypeOf("")).Interface().(string)), nil
	case reflect.Slice:
		if src.IsNil() {
			return apivalues.None, nil
		}

		vs := make([]*apivalues.Value, src.Len())
		for i := 0; i < src.Len(); i++ {
			var err error
			if vs[i], err = ConvertFromAWS(
				append(pathElems, fmt.Sprintf("%d", i)),
				src.Index(i),
			); err != nil {
				return nil, err
			}
		}
		return apivalues.List(vs...), nil
	case reflect.Struct:
		m := make(map[string]*apivalues.Value, src.NumField())
		for fi := 0; fi < src.NumField(); fi++ {
			ft, fv := t.Field(fi), src.Field(fi)

			// TODO: Include metadata.
			if ft.Name == "ResultMetadata" {
				continue
			}

			if !ft.IsExported() {
				continue
			}

			var err error
			if m[ft.Name], err = ConvertFromAWS(
				append(pathElems, ft.Name),
				fv,
			); err != nil {
				return nil, err
			}
		}
		return apivalues.DictFromMap(m), nil
	default:
		return nil, fmt.Errorf("%q: unhandled data type", path)
	}
}

func ConvertToAWS(pathElems []string, dst reflect.Value, in *apivalues.Value) error {
	path := strings.Join(pathElems, ".")

	switch t := dst.Type(); t.Kind() {
	case reflect.Ptr:
		switch tt := t.Elem(); tt.Kind() {
		case reflect.Int64, reflect.Int32, reflect.Int16, reflect.Int8, reflect.Int:
			iptr := apivalues.GetConcretValue[apivalues.IntegerValue](in)
			if iptr == nil {
				return fmt.Errorf("%q: must have integer value", path)
			}

			if err := apivalues.UnwrapInto(dst.Interface(), *iptr); err != nil {
				return fmt.Errorf("%q: %w", path, err)
			}

			return nil
		case reflect.Bool:
			bptr := apivalues.GetConcretValue[apivalues.BooleanValue](in)
			if bptr == nil {
				return fmt.Errorf("%q: must have boolean value", path)
			}

			if err := apivalues.UnwrapInto(dst.Interface(), *bptr); err != nil {
				return fmt.Errorf("%q: %w", path, err)
			}

			return nil
		case reflect.String:
			sptr := apivalues.GetConcretValue[apivalues.StringValue](in)
			if sptr == nil {
				return fmt.Errorf("%q: must have string value", path)
			}

			if err := apivalues.UnwrapInto(dst.Interface(), *sptr); err != nil {
				return fmt.Errorf("%q: %w", path, err)
			}

			return nil
		case reflect.Struct:
			indict := apivalues.GetConcretValue[apivalues.DictValue](in)
			if indict == nil {
				return fmt.Errorf("%q: must have a dict value", path)
			}

			indictm := make(map[string]*apivalues.Value, len(*indict))
			indict.ToStringValuesMap(indictm)

			v := dst.Elem()

			for fi := 0; fi < v.NumField(); fi++ {
				fv, ft := v.Field(fi), tt.Field(fi)

				if !ft.IsExported() {
					continue
				}

				inv, ok := indictm[ft.Name]
				if !ok {
					// no such field in input, leave at zero value
					continue
				}

				switch ft.Type.Kind() {
				case reflect.Ptr:
					fv.Set(reflect.New(ft.Type.Elem()))
				case reflect.Slice:
					fv.Set(reflect.MakeSlice(ft.Type, 0, 0))
				default:
					return fmt.Errorf("%q: not a ptr, array, or slice", path)
				}

				if err := ConvertToAWS(append(pathElems, ft.Name), fv, inv); err != nil {
					return err
				}
			}

			return nil
		default:
			return fmt.Errorf("%q is not of a handled type", path)
		}
	case reflect.Slice:
		inlist := apivalues.GetConcretValue[apivalues.ListValue](in)
		if inlist == nil {
			return fmt.Errorf("%q: must have a list value", path)
		}

		for i, inv := range *inlist {
			v := reflect.New(t.Elem())

			if err := ConvertToAWS(
				append(pathElems, fmt.Sprintf("%d", i)),
				v,
				inv,
			); err != nil {
				return err
			}

			dst.Set(reflect.Append(dst, v.Elem()))
		}

		return nil
	default:
		return fmt.Errorf("%q is not a pointer, slice, or an array", path)
	}
}
