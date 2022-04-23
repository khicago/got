package annoparse

import (
	"github.com/khicago/got/util/reflecto"
	"github.com/khicago/got/util/strs/kvstr"
)

type (
	FieldPath = string

	IAnnotationObject interface {
		TagName() string
	}

	AnnoTable[TAnno IAnnotationObject] map[FieldPath]TAnno
)

func (table AnnoTable[TAnno]) Get(fieldPath FieldPath) (TAnno, bool) {
	t, ok := table[fieldPath]
	return t, ok
}

// ExtractAnno
// create AnnoTable of TAnno : IAnnotationObject for each field of given targetObject's structure
// cache by tagName is recommended
func ExtractAnno[TAnno IAnnotationObject](targetObject any, defaultAnno TAnno) (AnnoTable[TAnno], error) {
	ret := make(map[FieldPath]TAnno)
	spawn := reflecto.NewTSpawner(defaultAnno)
	if err := reflecto.ForEachField(targetObject, func(fieldContext reflecto.FieldContext) error {
		anno := spawn()
		tagContent, ok := fieldContext.Tag.Lookup(anno.TagName())
		if !ok {
			return nil
		}
		if err := ExtractAnnoValuesFromTag(anno, tagContent); err != nil {
			return err
		}
		ret[fieldContext.Path] = anno

		return nil
	},
		reflecto.ForEachFieldOptions.OnlyExported(),
		reflecto.ForEachFieldOptions.Drill(-1),
	); err != nil {
		return nil, err
	}

	return ret, nil
}

func ExtractAnnoValuesFromTag[TAnno IAnnotationObject](annoObj TAnno, tag string) error { 
	_, err = kvstr.KVStr(tag).ReflectTo(annoObj, &kvstr.QueryOption{ 
		TrySnake:   true})
	return err
}
