package encryptor
//定义加密接口
type Interface interface{
	Digest(password string)(string, error)  //生成加密数字摘要
	Compare(hashedPassword ,password string) error //比对数字签名和密码名文
}