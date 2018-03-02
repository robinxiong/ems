package cache


type CacheStoreInterface interface{
	Get(key string) (string, error)
	Unmarshal(key string, object interface{}) error  //类似于Get, 但它将找到的值，保存进object
	Set(key string, value interface{}) error
	Fetch(key string, fc func()interface{}) (string, error) //找到key相关的值，如果没有找到，则保存相关的值到 CacheStore
	Delete(key string) error
}