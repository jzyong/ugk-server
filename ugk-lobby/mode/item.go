package mode

// Item 道具
type Item interface {
	GetCount() uint64 //数量
	GetType() uint8   //道具类型
}

// CommonItem 通用道具属性
type CommonItem struct {
	Count uint64 `count` //数量
	Type  uint8  `type`  //类型
}

func (item *CommonItem) GetCount() uint64 {
	return item.Count
}
func (item *CommonItem) GetType() uint8 {
	return item.Type
}

// PropertyItem 属性类道具
type PropertyItem struct {
	CommonItem
}
