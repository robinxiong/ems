package serializable_meta

import "ems/admin"

type SerializableMetaInterface interface {
	GetSerializableArgumentResource() *admin.Resource
	GetSerializableArgument(SerializableMetaInterface) interface{}
	GetSerializableArgumentKind() string
	SetSerializableArgumentKind(name string)
	SetSerializableArgumentValue(interface{})
}

type SerializableMeta struct {
	Kind  string
	Value serializableArgument `sql:"size:65532"`
}
type serializableArgument struct {
	SerializedValue string
	OriginalValue   interface{}
}
