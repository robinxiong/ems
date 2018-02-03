package filesystem

import (
	"testing"
	"ems/oss/tests"
)

func TestNew(t *testing.T){
	fileSystem := New("/tmp")  //创建一个保存资源的平台，可以是本地的文件系统,也可以是s3， 它用于获取对像，文件，保存对像
	tests.TestAll(fileSystem, t)
}