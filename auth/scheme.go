package auth
//Schema 认证模式，可以是邮箱，手机，github, facebook帐号等
type Schema struct {
	Provider string
	UID      string

	Name      string
	Email     string
	FirstName string
	LastName  string
	Location  string
	Image     string
	Phone     string
	URL       string

	RawInfo interface{} //原始的req, 它包含了注册form信息
}
