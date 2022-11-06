package annoparse

type (
	SelfRefTag struct {
		Name string
	}
)

func (*SelfRefTag) TagName() string {
	return "anno"
}

var annoKeyMapping = map[string]map[FieldPath]string{
	(&SelfRefTag{}).TagName(): nil, // SelfRefTag 自身不支持
}

func getAnnoKeyMapping[TAnno IAnnotationObject](annoObj TAnno) (map[FieldPath]string, error) {
	if v, ok := annoKeyMapping[annoObj.TagName()]; ok {
		return v, nil
	}

	m, err := ExtractAnno(annoObj, &SelfRefTag{})
	if err != nil {
		return nil, err
	}
	ret := make(map[FieldPath]string)
	for fieldPath, tagSetting := range m {
		// support convert the field name to key name in the tag
		ret[fieldPath[1:]] = tagSetting.Name
	}
	annoKeyMapping[annoObj.TagName()] = ret
	// fmt.Println("====>", ret)
	return ret, nil
}
