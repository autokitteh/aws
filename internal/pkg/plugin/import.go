package plugin

import (
	"context"
	"fmt"
	"reflect"

	"go.autokitteh.dev/sdk/api/apivalues"
	"go.autokitteh.dev/sdk/pluginimpl"
)

func importServiceMethods(
	client interface{}, // only used to build the plugin, not actually used for connecting.
) map[string]pluginimpl.PluginMethodFunc {
	t := reflect.TypeOf(client)
	if t.Kind() != reflect.Ptr || t.Elem().Kind() != reflect.Struct {
		panic("client is not a pointer to a struct")
	}

	methods := make(map[string]pluginimpl.PluginMethodFunc, t.NumMethod())

	for mi := 0; mi < t.NumMethod(); mi++ {
		m := t.Method(mi)

		methods[m.Name] = func(
			ctx context.Context,
			name string,
			args []*apivalues.Value,
			kwargs map[string]*apivalues.Value,
			_ pluginimpl.FuncToValueFunc,
		) (*apivalues.Value, error) {
			var paramsArg *apivalues.Value

			if err := pluginimpl.UnpackArgs(args, kwargs, "params", &paramsArg); err != nil {
				return nil, err
			}

			if apivalues.GetConcretValue[apivalues.DictValue](paramsArg) == nil {
				return nil, fmt.Errorf("params must be a dict")
			}

			mt := m.Type

			// Expecting self, context, params, optFns.
			if mt.NumIn() != 4 {
				panic(fmt.Errorf("method %q numin %d != 4", m.Name, mt.NumIn()))
			}

			pt := mt.In(2)
			if pt.Kind() != reflect.Ptr || pt.Elem().Kind() != reflect.Struct {
				panic(fmt.Errorf("method %q param invalid type: %v", m.Name, pt))
			}

			paramsValue := reflect.New(pt.Elem())

			if err := ConvertToAWS(nil, paramsValue, paramsArg); err != nil {
				return nil, err
			}

			client := ec2Client()
			clientv := reflect.ValueOf(client)

			method := clientv.MethodByName(name)

			retvs := method.Call([]reflect.Value{
				reflect.ValueOf(ctx),
				paramsValue,
			})

			if len(retvs) != 2 {
				return nil, fmt.Errorf("call returned %d values != expected 2", len(retvs))
			}

			outv, errv := retvs[0], retvs[1]

			if !errv.IsNil() {
				if err, ok := errv.Interface().(error); ok {
					return nil, err
				}

				return nil, fmt.Errorf("invalid error return")
			}

			out, err := ConvertFromAWS(nil, outv)
			if err != nil {
				return nil, fmt.Errorf("return value conversion error: %w", err)
			}

			return out, nil
		}
	}

	return methods
}

func importService(
	name string,
	client interface{}, // only used to build the plugin, not actually used for connecting.
) *pluginimpl.PluginMember {
	methods := importServiceMethods(client)

	return pluginimpl.NewLazyValueMember(
		fmt.Sprintf("%s service API", name),
		func(ftov pluginimpl.FuncToValueFunc) *apivalues.Value {
			members := make(map[string]*apivalues.Value, len(methods))
			for k, f := range methods {
				members[k] = ftov(k, f)
			}

			return apivalues.Module(name, members)
		},
	)
}
