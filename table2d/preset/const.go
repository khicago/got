// Package preset
// e.g. a preset table of `hero`
//
// |@    |ID     |INT           |Float    |[         |       |      ]|[        |{|          |#        |          |#    |          |}|]|BOOL  |STRING          |
// |     |link(@)|test($>1,$<50)|test($%2)|link(item)|       |       |         | |link(wear)|         |link(wear)|     |link(wear)| | |      |enum(A|S|SR|SSR)|
// |     |LvUp   |Power         |Magic    |InitItems |       |       |InitStuff| |Weapon    |         |Shoes     |     |Hat       | | |IsHero|Gene            |
// |10001|10002  |12            |         |          |1010001|1010002|         | |1020001   |The Sword|1020101   |Speed|1030001   | | |Y     |S               |
package preset

import "errors"

var (
	ErrColOutOfRange = errors.New("preset err: col out of range")
	ErrPropertyNil   = errors.New("preset err: the property is nil")
	ErrPropertyType  = errors.New("preset err: property type error")
)

const (
	InvalidCol Col = -1
)
