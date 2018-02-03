# media
model struct中包含了接口时

1. 先处理SerializableMetaInterface， 它的参数中有Media时，则处理
2. 接着检查model的字段是否实现了Media接口，这里为Image字段

```go
type QorJob struct {
  gorm.Model
  Name string
  Image oss.OSS 
  serializable_meta.SerializableMeta
}
```