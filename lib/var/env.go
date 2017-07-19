package var_

import "os"

func (_ Module) Env(name string) interface{} {
	return os.Getenv(name)
}
