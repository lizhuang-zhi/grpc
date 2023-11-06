package common

func DataList() []string {
	data := make([]string, 0)
	for i := 0; i < 10; i++ {
		data = append(data, "测试数据")
	}
	return data
}
