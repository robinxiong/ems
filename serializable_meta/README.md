#serializable_meta
将一个对像序列化后保存到数据库

```go
type QorJob struct {
  gorm.Model
  Name string
  /*
    Kind  string
  	Value serializableArgument `sql:"size:65532"`
  */
  serializable_meta.SerializableMeta
}

func (qorJob QorJob) GetSerializableArgumentResource() *admin.Resource {
  return jobsArgumentsMap[qorJob.Kind]
}


//注册不同的model struct, 用于序列化
var jobsArgumentsMap = map[string]*admin.Resource{
  "newsletter": admin.NewResource(&sendNewsletterArgument{}),
  "import_products": admin.NewResource(&importProductArgument{}),
}

type sendNewsletterArgument struct {
  Subject string
  Content string
}

type importProductArgument struct {}

var qorJob QorJob
qorJob.Name = "sending newsletter"
qorJob.Kind = "newsletter"
qorJob.SetSerializableArgumentValue(&sendNewsletterArgument{
  Subject: "subject",
  Content: "content",
})

db.Create(&qorJob)
//INSERT INTO "qor_jobs" (kind, value) VALUES (`newsletter`, `{"Subject":"subject","Content":"content"}`);
var result QorJob
db.First(&result, "name = ?", "sending newsletter")

var argument = result.GetSerializableArgument(result)
argument.(*sendNewsletterArgument).Subject // "subject"
argument.(*sendNewsletterArgument).Content // "content"
```
每一个QorJob都有一个类型，根据Kind的不同，返回不同的Value(serializableArgument), serializableArgument包含了SerializeValue和OriginalValue. 它跟JSON不同在于，反序列化后，可以调用这些定义和注册的struct的方法，就像一个正常的go struct, 而json数据不是

