# OSS
oss提供了公共接口，用来操作云存储，ftp, 文件系统中的文件

# Usage

当前, OSS仅支持S3, Qiniu. 你需要实现自己的通用接口，实现自己的storage

```go
//StorageInterface定义公用API用来操作存储
type StorageInterface interface {
  Get(path string) (*os.File, error)
  Put(path string, reader io.Reader) (*Object, error)
  Delete(path string) error
  List(path string) ([]*Object, error)
  GetEndpoint() string
}
```
以下是一个使用s3的OSS.
```go
import (
	"github.com/oss/filesystem"
	"github.com/oss/s3"
	awss3 "github.com/aws/aws-sdk-go/s3"
)

func main() {
	storage := s3.New(s3.Config{AccessID: "access_id", AccessKey: "access_key", Region: "region", Bucket: "bucket", Endpoint: "cdn.getqor.com", ACL: awss3.BucketCannedACLPublicRead})
	// storage := filesystem.New("/tmp")

	// Save a reader interface into storage
	storage.Put("/sample.txt", reader)

	// Get file with path
	storage.Get("/sample.txt")

	// Delete file with path
	storage.Delete("/sample.txt")

	// List all objects under path
	storage.List("/")
}
```