package structcopy

import "encoding/json"

func JsonDeepCopy(src, dest interface{}) error {
	jsrc, err := json.Marshal(src)
	if err != nil {
		return err
	}

	err = json.Unmarshal(jsrc, &dest)
	if err != nil {
		return err
	}

	return nil
}

func ReflectDeepCopy(src, dest interface{}) error {

	return nil
}
