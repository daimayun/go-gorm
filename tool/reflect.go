package tool

import "reflect"

// GetStructFullName get reflect.Type name with package path.
func GetStructFullName(typ reflect.Type) string {
	return typ.PkgPath() + "." + typ.Name()
}

func GetFuncReturnValue(val reflect.Value, key string) (ok bool, value reflect.Value) {
	if fun := val.MethodByName(key); fun.IsValid() {
		values := fun.Call([]reflect.Value{})
		if len(values) > 0 {
			ok = true
			value = values[0]
		}
	}
	return
}
