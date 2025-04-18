package utils

type BindFunc func(interface{}) error

func BindFlow(obj interface{}, binds ...BindFunc) error {
	for _, bind := range binds {
		if err := bind(obj); err != nil {
			return err
		}
	}
	return nil
}
