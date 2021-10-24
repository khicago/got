package kvstr

import "github.com/khicago/got/util/strs"

type (
	KVMap map[KeyName]string

	QueryOption struct {
		// accept more key that can be converted to the key of kvMap
		TryMapping map[string]KeyName
		TrySnake   bool
		TryCamel   bool
	}
)

func (m KVMap) TryGet(name string, opt *QueryOption) (string, bool) {
	str, ok := m[name]
	if ok {
		return str, true
	}

	if nil == opt {
		return str, ok
	}

	if nil != opt.TryMapping {
		mapKey, has := opt.TryMapping[name]
		if has {
			if str, ok = m[mapKey]; ok {
				return str, ok
			}
		}
	}

	if opt.TrySnake {
		if str, ok = m[strs.Conv2Snake(name)]; ok {
			return str, ok
		}
	}

	if opt.TryCamel {
		if str, ok = m[strs.Conv2Camel(name)]; ok {
			return str, ok
		}
	}

	return "", false
}
