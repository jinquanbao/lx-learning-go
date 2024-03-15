package excelutil

/***
excel tag支持的标签:
name:标题名称
index: 标题索引，从0开始
cell:cell标题索引, 从"A"开始
converter: 方法转换器，由于是反射调用，方法名称必须是大写

type User struct {
	Id      int       `json:"id" excel:"index:1 ; name:用户id; "`
	Name    string    `json:"name" excel:"name:用户名; "`
	Time    time.Time `json:"time" excel:"name:创建时间; "`
	Roles   []string  `json:"roles" excel:"name:角色; converter:ConvertRoles; "`
	Message string    `json:"message" excel:"index:5; name:消息; "`
}

func (u *User) ConvertRoles(cellVal string) error {
	if cellVal == "" {
		return nil
	}
	u.Roles = strings.Split(cellVal, "|")
	return nil
}

*/
