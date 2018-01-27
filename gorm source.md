package main

import (
	"log"
	"github.com/jinzhu/gorm"
)

## Associations
> https://gist.github.com/ShionRyuu/00385f4959884386ac72
> https://ruby-china.github.io/rails-guides/association_basics.html
将实体与实体的关系，反应到最终数据库的设计上来，将关系分为：一对一，一对多，多对多
所有的关系都是指的表与表之间的关系

一对一
一张表的一条记录一定只能与另外一张表的一条记录进行对应，反之亦然。

学生表：姓名，性别，年龄，身高，体重，籍贯，家庭住址，紧急联系人
其中姓名、性别、年龄、身高，体重属于常用数据，但是籍贯、住址和联系人为不常用数据
如果每次查询都是查询所有数据，不常用的数据就会影响效率，实际又不用
常用信息表：ID(P)，姓名，性别，年龄，身高，体重
不常用信息表：ID(P)，籍贯，家庭住址，紧急联系人

解决方案：将常用的和不常用的信息分享存储，分成两张表
不常用信息表和常用信息表，保证不常用信息表与常用信息表能够对应上：找一个具有唯一性的
字段来共同连接两张表。
一个常用表中的一条记录永远只能在一张不常用表中匹配一条记录，反之亦然。


一对多
一张表中有一条记录可以对应另外一张表中的多条记录；但是反过来，另外一张表的一条记录
只能对应第一张表的一条记录，这种关系就是一对多或多对一
母亲与孩子的关系：母亲，孩子两个实体
母亲表：ID(P),名字，年龄，性别
孩子表：ID(P),名字，年龄，性别
以上关系：一个妈妈可以在孩子表中找到多条记录（也可能是一条），但是一个孩子只能找到一个妈妈
是一种典型的一对多的关系。
但是以上设计：解决了实体的设计表问题，但是没有解决关系问题，孩子找不到母亲，母亲也找不到孩子

解决方案：在某一张表中增加一个字段，能够找到另外一张表中的记录:在孩子表中增加一个字段
指向母亲表，因为孩子表的记录只能匹配到一条母亲表的记录。
母亲表：ID(P),名字，年龄，性别
孩子表：ID(P),名字，年龄，性别，母亲表ID（母亲表主键）


多对多
一对表中（A）的一条记录能够对应另外一张表（B）中的多条记录；同时B表中的一条记录
也能对应A表中的多条记录

老师和学生
老师表 T_ID(P),姓名，性别
学生表 S_ID(P),姓名，性别
以上设计方案：实现了实体的设计，但是没有维护实体的关系
一个老师教过多个学生，一个学生也被多个老师教过

解决方案：增加一张中间关系表
老师与学生的关系表：ID(P),T_ID,S_ID
老师表与中间表形成一对多的关系，而中间表是多表；维护了能够唯一找到一表的关系；
同样的学生表与中间表也是一个一对多的关系;
学生找老师：找出学生ID--->中间表寻找匹配记录（多条）--->老师表匹配（一条）
老师找学生：找出老师ID--->中间表寻找匹配记录（多条）--->学生表匹配（一条）

### Belongs To
belongs_to 关联创建两个模型之间一对一的关系，声明所在的模型实例属于另一个模型的实例。例如，如果应用中有作者和图书两个模型，而且每本书只能指定给一位作者，就要这么声明图书模型, 属于一对一

```
// `User` belongs to `Profile`, `ProfileID` is the foreign key 默认情况下不会在User表中保存信息， 而且以字段值+ID为外键
// 可以通过`gorm:"ForeignKey:ProfileRefer"`指定哪个字段作为外键
type User struct {
	gorm.Model
	Profile   Profile
	ProfileID int
}

type Profile struct {
	gorm.Model
	Name string
}
user := User{
		Profile: Profile{
			Name: "test",
		},
	}

db.Create(&user)

```
也可以指定外键为某一列
```
type Profile struct {
    gorm.Model
    Refer int
    Name  string
}

type User struct {
    gorm.Model
    Profile   Profile `gorm:"ForeignKey:ProfileID;AssociationForeignKey:Refer"`
    ProfileID int
}
```
### 在 belongs_to 和 has_one 之间选择

如果想建立两个模型之间的一对一关系，要在一个模型中添加 belongs_to，在另一模型中添加 has_one。但是怎么知道在哪个模型中添加哪个呢？

二者之间的区别是在哪里放置外键（外键在 belongs_to 关联所在模型对应的表中），不过也要考虑数据的语义。has_one 的意思是某样东西属于我，即哪个东西指向你。例如，说供应商有一个账户，比账户拥有供应商更合理

###Has Many
```
// User has many emails, UserID is the foreign key
type User struct {
    gorm.Model
    Emails   []Email
}

type Email struct {
    gorm.Model
    Email   string
    UserID  uint
}

db.Model(&user).Related(&emails)
```
指定外键
```
type Profile struct {
  gorm.Model
  Name      string
  UserRefer uint
}

type User struct {
  gorm.Model
  Profiles []Profile `gorm:"ForeignKey:UserRefer"`
}
```

指定外键和关联键
```
type Profile struct {
  gorm.Model
  Name   string
  UserID uint
}

type User struct {
  gorm.Model
  Refer   uint
  Profiles []Profile `gorm:"ForeignKey:UserID;AssociationForeignKey:Refer"`
}
```

###Polymorphism
关联还有一种高级形式——多态关联（polymorphic association）。在多态关联中，在同一个关联中，一个模型可以属于多个模型。例如，图片模型可以属于雇员模型或者产品模型，模型的定义如下
当前gorm只支持has-many and has-one 关系，不支持belongs-to and many-to-many
在gorm中默认的是通过 OwnerType和OwnerId来连接
```
type Cat struct {
    Id    int
    Name  string
    Toy   Toy `gorm:"polymorphic:Owner;"`
  }

  type Dog struct {
    Id   int
    Name string
    Toy  Toy `gorm:"polymorphic:Owner;"`
  }

  type Toy struct {
    Id        int
    Name      string
    OwnerId   int
    OwnerType string
  }
```
###Association Mode
Association Mode 包含许多Helper方法来处理关系
```
// Start Association Mode
var user User
db.Model(&user).Association("Languages")
// `user` is the source, it need to be a valid record (contains primary key)
// `Languages` is source's field name for a relationship.
// If those conditions not matched, will return an error, check it with:
// db.Model(&user).Association("Languages").Error


// Query - Find out all related associations 找到当前user相关联的数据
db.Model(&user).Association("Languages").Find(&languages)


// Append - Append new associations for many2many, has_many, will replace current association for has_one, belongs_to
db.Model(&user).Association("Languages").Append([]Language{languageZH, languageEN})
db.Model(&user).Association("Languages").Append(Language{Name: "DE"})


// Delete - Remove relationship between source & passed arguments, won't delete those arguments
db.Model(&user).Association("Languages").Delete([]Language{languageZH, languageEN})
db.Model(&user).Association("Languages").Delete(languageZH, languageEN)


// Replace - Replace current associations with new one
db.Model(&user).Association("Languages").Replace([]Language{languageZH, languageEN})
db.Model(&user).Association("Languages").Replace(Language{Name: "DE"}, languageEN)


// Count - Return the count of current associations
db.Model(&user).Association("Languages").Count()


// Clear - Remove relationship between source & current associations, won't delete those associations
db.Model(&user).Association("Languages").Clear()
```

## 预备知识

```
//gorm join_table_handler.go
type JoinTableHandler struct {
	TableName   string          `sql:"-"`
	Source      JoinTableSource `sql:"-"`
	Destination JoinTableSource `sql:"-"`
}
joinTableHandler := JoinTableHandler{}
joinTableHandler.Setup(relationship, "product_categories", Source, Destination)
```
生成关系表时处理器, 通常使用Setup方法， relationship包含类型， 外键字段名称，关联字段名称, 第二个参数为关联表的类型，第三个为source model类型
第四人为 distination model struct


//gorm model_struct.go
type Relationship struct {
	Kind                         string
	PolymorphicType              string
	PolymorphicDBName            string
	PolymorphicValue             string
	ForeignFieldNames            []string
	ForeignDBNames               []string
	AssociationForeignFieldNames []string
	AssociationForeignDBNames    []string
	JoinTableHandler             JoinTableHandlerInterface
}
```

Kind可以为many_to_many, has_many等
PolymorphicType为设置了Polymorphic时，才会设置的值, 表示的是has-one, has-many, belong-to的关系
ForeignFieldNames 保存source的主键名称
ForeignDBNames 生成关系表时的名称比如product_id
AssociationForeignFieldNames 相关表的主键名称
AssociationForeignDBNames 生成关系表时，Destination的主键dbname,  比如category_id



生成_draft表

```
type Product struct {
	ID         int        `gorm:"primary_key"`
	Categories []Category `gorm:"many2many:product_categories;ForeignKey:id;AssociationForeignKey:id"`  //2-j
	Brand      Brand  //参考 2->g
	locale //参考2->i
}
type Category struct {
	ID   int `gorm:"primary_key"`
	Name string
}
type Brand struct {
	ID   int `gorm:"primary_key"`
	Name string
	Locale
}
type Locale struct {
	LanguageCode string `sql:"size:20" gorm:"primary_key"`
}
func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	db := utils.TestDB()
	db.DropTableIfExists(&Product{})
	db.DropTableIfExists(&Category{})
	db.Exec("drop table product_categories;")
	//创建producs, categories表
	db.AutoMigrate(&Product{}, &Category{}, &Brand)
}
```
以下是普通的AutoMigrate过程
1. 创建&Product{}, 生成表名products
2. 遍历每一个字段_, field := range scope.GetModelStruct().StructFields, 解析过程如下
    首先解析tag
    ---------------------
    a. 解析字段的sql， gorm tag, 并且以;分割每个gorm设置, 返回一个map
    b. 如果在gorm设置为"-", 则field.IsIgnored = true
    c. "PRIMARY_KEY", 则field.IsPrimaryKey= true, 同时将这个字段保存到model.PrimaryFields数组中(哪些是主键StrctField字段)
    d. "DEFAULT" tag field.HasDefaultValue = true
    e. "AUTO_INCREMENT" tag field.HasDefaultValue = true
    接下来解析model struct中字段的类型
    f. 如果字段的类型为指针，则返回它的实际类型fieldStruct.Type.Elem(), 通过reflect.New(indirectType).Interface()返回一个接口值
    g. 接口值实现了sql.Scanner接口, 比如sql.NULLString, 其它的原生类型(int, string, bool, float64)或者自定义类似都没有实现此接口，设置field.IsScanner, field.IsNormal = true, true (sql.Scanner可以参考https://golang.org/src/database/sql/sql.go),
       如果这个字段的类型同时是一个struct类型，则解析这个子struct的每一个字段的tag, 如果找到而且field.TagSettings中没有设置，则添加
    h. 如果没有实现Scanner, 则判断是否为*time.Time， 如果是，则field.IsNormal = true
    i. 如果也不是*time.Time, 判断是否设置了tag为"EMBEDDED", 同时fieldStruct.Anonymous, 则调用调用递归GetModelStruct()，然后continue跳过后面的代码， 执行下一个field解析, 解析出它所有的字段, 过程如下
        1. 将当前field的名字，添加到 subField.Names中
        2. 如果field中设置了 EMBEDDED_PREFIX tag, 则subFeild在数据库表的名字为subField.DBName = prefix + subField.DBName
        3. 如果subfield.IsPrimaryKey, 则添加到当前modelStruct.PrimaryFields
        4. 如果subField.Relationship和subField.Relationship.JoinTableHandler 非nil, 即tag中设置了many2many tag。则设置一个新的JoinTableHandler = newJoinTableHandler.Setup(subField.Relationship, joinTableHandler.TableName, reflectType, joinTableHandler.Destination.ModelType)

    j. 如果字段类型是一个slice, 即这个字段包含多个值, 它使用了defer特性Last In First Out, 即从最后一个字段开始解析它的跟其它表的关系, 表示many-to-many has-many
        1. 字段是否设置了FOREIGNKEY, 如果设置了，则添加到foreignKeys数据中
        2. ASSOCIATIONFOREIGNKEY tag, 添加到associationForeignKeys
        3. 解析当前slice的值的类型
        4. 哪果字段设置了MANY2MANY, 同时slice保存的是struct类型
            1). 创建一个model_struct.go Relationship > relationship
            2). relationship.Kind = "many_to_many"
            3). 如果之前没有解析到FOREIGNKEY  tag， 则将整个modelstruct中包含的主键添加到 foreignKeys数据中
            4). foreignKeys中保存的是字符串，所以我们需要通过getForeignField方法来查找它所对应的field, 如果找到field, 比如id所对应的field(这也是为什么使用LIFO的原因, 能够查找所有之前设置的field)
                将这个field的DBName添加到relationship.ForeignFieldNames数组中(source foreight keys), 同时为source设置连接表的外键为reflectType.Name() +"_" +foreignField.DBName, 即product_id
                将它添加到ForeignDBNames中
            5). 如果没有找到ASSOCIATIONFOREIGNKEY，则创建一个新的toScope(scope.New(&Category)), 它的值为字段的类型(这里为Category), 它也会调用GetModelStruct()进解析， 递归，然后将它的主键添加到associationForeignKeys
            6). 类似于step 4, category中的主键添加到relationship.AssociationForeignFieldNames中, 同时每一个字段的名称为category_主键名称, 并添加到relationship.AssociationForeignDBNames中
            7). 设置新建一个JoinTableHandler, 调用它的Setup方法，第一个参数为刚刚设置的relationship, 第二个为many2many的值product_categories, 第三个参数为 Product struct， 第四个为 Category struct
                将这个JoinTableHandler设置为relationship.JoinTableHandler, 然后field.relationship = relationship
        5. 如果没有设置任何tag, 则默认为has_many, 比如User has many comments, associationType 为 User, comment 使用 UserID 作为 foreign key
            1). relationship.Kind = "has_many"， associationType=表名(默认)
            2). 如果field设置Polymorphic, 比polymorphic:Owner, 则在toFileds中找到OwnerType字段, 设置assoctionType为PLYMORPHIC标签的值
                找到字段 relationship.PolymorphicType = polymorphicType.Name (OwnerType), relationship.PolymorphicDBName = polymorphicType.DBName
                如果field中设置了POLYMORPHIC_VALUE, 则为relationship.PolymorphicValue = value, 否则为表的名字
                设置OwenType字段的IsForeignKey = true
            3). 如果没有指定foreignkey, 也没有指定associationForeignKey, 遍历当前表的主键, 所以foreignKeys为associationType(表名)+field.Name(主键), 而associationForeignKeys为field.Name
            4). 如果没有批定foreignkey, 而指定了associationForeignKey, 则遍历associationForeignKey, 并按照上一步的方式，生成foreignKeys
            5). 指定foreignKey, 则按照上面两步的方法，生成foreignKeys,associationForeignKeys
            6). 遍历foreignKeys, 从toFields中找到foreignKey的Field作为foreignField, 在从本表中找到对应的associationField
                设置foreignField.IsForeignKey=true
                relationship.AssociationForeignFieldNames为当前字段的名字
                relationship.AssociationForeignDBNames为当前字段的数据库名字
                relationship.ForeignFieldNames为toScope中相关字段的名字
                relationship.ForeignDBNames为toScope数据库表的名字
            7)  只有在relationship.ForeignFieldNames != 0时，才设置 field.Relationship = relationship


    k. 如果字段是一个struct类型，同时没有实现sql.Scanner接口， 表示has-one或者belongs to. 如果 user has one profile, associationType 为 User, profile 使用 UserID 作为 foreign key
        如果user belongs to profile, associationType 为 Profile, user 使用 ProfileID 作为 foreign key, 此步聚跟j->5类似, 以之前belong to中讲过的例子为例
        ```
type Profile struct {
	gorm.Model
	Refer int
	Name  string
}
//belongs-to
type User struct {
	gorm.Model
	Profile   Profile `gorm:"ForeignKey:ProfileID;AssociationForeignKey:Refer"`
	ProfileID int
}

//has-one
type Cat struct {
	Id    int
	Name  string
	Toy   Toy `gorm:"polymorphic:Owner;"`
}

type Dog struct {
	Id   int
	Name string
	Toy  Toy `gorm:"polymorphic:Owner;"`
}

type Toy struct {
	Id        int
	Name      string
	OwnerId   int
	OwnerType string
}
```
        1. 检查聚合的struct是否包含FOREIGNKEY 标签, 如果包含，则添加到tagForeignKeys中, 然后检查是否包含ASSOCIATIONFOREIGNKEY, 添加到tagAssociationForeignKeys
        2. 检查是否设置了POLYMORPHIC标签, 如果设置了，比如Owner,则在聚合的类中(Toy)查找Owner+"Type"字段，则设置区分不同表的列为OwnerType, 即通过OwnerType来查询toy所对应的拥有者，cat or dog
           associationType变量由表名Cat改为Owner
           relationship.PolymorphicType = OwnerType
           relationship.PolymorphicDBName = owner_type,
           接着检查field中是否设置了POLYMORPHIC_VALUE
           是：relationship.PolymorphicValue=设置的值
           否：relationship.PolymorphicValue=表名(通过PolymorphicValue在PolymorphicDBName查询它相关的值, cat or dog)
        3. has-one的检查
            1).没有定义foreignKey 标签, 也没有定义associationForeignKey标签，则外键为 表名+表的主键, 比如这里为UserId 或者 CatId, DogId, OwnerId(多态), 关联键为user表的Id键 (has one)
            2).没有定义foreignKey标签，但设置了associationForeignKey, 则从associationForeignKey中生成外键, 从当前表中查找associationForeignKey字段，找到后外键为表名+Found Value, 关联键为 found field.name
            3).设置了foreignKey标签, 没有设置associationForeignKey, 外键是否以表名开头, 如果是，关联键去除掉表名, 如果在当前表中存在这个关联键，则添加到associationForeignKeys, 如果在当前表中不存在此键，则关联键为当前表的主键, associationForeignKeys > foreignKeys
            4). 即存在foreignKey也存在associationForeignKeys, 检查是否相等，如果不相等，则报错

        4. 在聚合表中查找外键(Toy)中查询外键, 这里为OwnerId, 如果找到, 则在当前表对应的关联键associationForeignKey, 如果都找到
            foreignField.IsForeignKey = true (Toy.OwnerId field)
            relationship.AssociationForeignFieldNames 设置为当前表的id
            relationship.ForeignFieldNames 则设置为Toy表的owner_id
            所以Cat的Toy字段可以描述为
            OwnerId或者CatId是Toy 模型的外键, 并且设置Toy模型(自动解析Toy模型，并缓存为gorm.ModelStruct)中OwnerId或者CatId字段为外键foreignField.IsForeignKey
            它参照的键为Cat表的Id
        5. 如果在聚合表中找到了foreignKey字段(UserId, 或者说OwnerId), 则relationship.Kind="has_one", 否则继续检查是否为belong_to，即第6步, 在我们的例子中的User, 它就不在Profile中包含UserId, 所以它不是has-one的关系
        6. belongs-to
            1). 没有找到foreignKey, 同时也没有设置associationForeignKey, 则外键为当前字段名字+toScope的主键, ProfileId， associationForeignKey为Profile的Id
            2). 没有找到foreignKey, 但设置了associationforeignKey, 则在Profile表中查找 "Refer"字段， 如果找到，则外键为当前字段名+Refer, 即ProfileRefer, associationForekgnKeys为Refer, len(associationForeignKeys) > len(foreignKeys)
            3). 设置了foreignKey, 没有设置associationForeignKey, (ProfileId), 先删除Profile, -> Id, 接着在Profile struct中查找id, 找到了，则赋值给associationForeignKeys
            4). 同时设置了foreignKey, associationforeignKey， 检查是否一样
        7. 在当前表中查找foreignKey, 这里为ProfileId, 然后在聚合表Profile中查找associationForeignKey, 如果找到，则
            当前字段为外键, User.Profile.IsForeignKey = true
            它的关系描术为，User表的(ForeignFieldName ProfileId)参考Profile(AassociationForeignFieldName 的 id或refer字段, distination)
        8. 不管是has-one or belong-to,  它们包含ForeignFieldNames为外键的字段名，相关联的AssociationForeignDBNames（参考)字段的名称，以及当前字段IsForeignKey是否为true, 这三个值
           如果包含polymorphic，则还包含polymorphicType字段，用于区分不同的表, 而它并不会向many2many 那样，加入JoinTableHandler, 所以在生成表时，不会生成关系表

    l. 常规值(int, string)，设置field.IsNormal=true
    m. 缓存当前解析的struct
3. scope.CreateTable 会遍历每一个当前表的field, 如果这个field主键，则添加到primaryKeys, 如果这个字段有relationship，则生成关系表, 全部字段处理完成之后，才生成当前数据库表, 以下是createJoinTable的处理方式
    a. 字段有Relationship, 并且relationship.JoinTableHandler不为nil
    b. 通过joinTableHandler可以获取到关系表的表名， 赋值给，并判断是否存在于数据库中
    c. 根据当前字段类型，创建一个toScope
    d. 遍历relationship中的ForeignFieldNames字段, 获取这些字段的名称，然后在当前scope中查找这些字段，然后添加到sqlTypes和primaryKeys
    e. 同理，在toScope中找到relationship.AssociationForeignFieldNames指定的键，也添加到sqlTypes, primaryKeys
    f. 调用原生sql, 将sqlTypes, primaryKeys传入给它， 创建joinTable表
    g. 调用scope.NewDB().Table(joinTable).AutoMigrate(joinTableHandler)，此句可以不需要, GetModelStruct并不能解析joinTableHandler成数据库表


##创建
Create and save都会调用, 默认的DefaultCallback
```
//对于作何创建，db.Create(&Product{})都开启一个事务, 回调函数位于callback_save.go文件, 并且将事务保存进scope.db.Db
DefaultCallback.Create().Register("gorm:begin_transaction", beginTransactionCallback)
DefaultCallback.Create().Register("gorm:before_create", beforeCreateCallback)   //位于callback_create.go, 它执行Product中定义的BeforeSave或者BeforeCreate方法
/*
更新或者创建时调用, 确认是否先保存关联表
先判断是否在scope中设置了gorm:save_associations
比如db.UpdateColumns会将Set("gorm:save_associations", false), 就是更新某些字段时不需要保存关联表数据, 只更新当前表, 通常这个值没有设置，所以默认为true
/*
举例来说 association_test.go

func TestSkipSaveAssociation(t *testing.T) {
	type Company struct {
		gorm.Model
		Name string
	}

	type User struct {
		gorm.Model
		Name      string
		CompanyID uint
		Company   Company `gorm:"save_associations:false"`
	}
	DB.AutoMigrate(&Company{}, &User{})

	DB.Save(&User{Name: "jinzhu", Company: Company{Name: "skip_save_association"}})
    //正常来说应该会在company中保存一条记录，但这里使用save_associations:false， 则直接跳过
	if !DB.Where("name = ?", "skip_save_association").First(&Company{}).RecordNotFound() {
		t.Errorf("Company skip_save_association should not been saved")
	}
}
注：tag的设置并不影响scope.Get("gorm:save_associations")
*/

以下是详细的方法
func saveBeforeAssociationsCallback(scope *Scope) {
	if !scope.shouldSaveAssociations() { //检查是否设置了gorm:save_associations
		return
	}
	//scope保存了value 为&User{Profile: Profile{name:"hello"}}
	//scope.Fields会调用scope.GetModelStruct()解析model的字段
	//在调用reflect.Indirect(scope.IndirectValue).FieldByName来查询 User中对应字段的值
	for _, field := range scope.Fields() {
		//saveFieldAsAssociation查看是否设置了SAVE_ASSOCIATIONS标签为skip or false, 如果是，则返回false
		//saveFieldAsAssociation同时类型为belongs_to,则先创建Profile的值
		if ok, relationship := saveFieldAsAssociation(scope, field); ok && relationship.Kind == "belongs_to" {

			//聚合model的值, profile
			fieldValue := field.Field.Addr().Interface()
			//保存Profile的值
			scope.Err(scope.NewDB().Save(fieldValue).Error)

			//设置当前表的外键值 ProfileId
			if len(relationship.ForeignFieldNames) != 0 {
				// set value's foreign key
				for idx, fieldName := range relationship.ForeignFieldNames {
					//查找relationship中跟foreignFieldNames外键相对应的associationForeignKey的列名, 这里为Profile表中的id
					associationForeignName := relationship.AssociationForeignDBNames[idx]
					//User中Profile查找id的值, 然后设置User外键的值ProfileId
					if foreignField, ok := scope.New(fieldValue).FieldByName(associationForeignName); ok {
						scope.Err(scope.SetColumn(fieldName, foreignField.Field.Interface()))
					}
				}
			}
		}
	}
}
*/
DefaultCallback.Create().Register("gorm:save_before_associations", saveBeforeAssociationsCallback)//位于callback_save.go文件,


//如果models中有CreateAt或者UpdateAt字段，则设置它们的值
DefaultCallback.Create().Register("gorm:update_time_stamp", updateTimeStampForCreateCallback)


/*
首先遍历所有的字段, 如果字段是需要更新的
	for _, field := range scope.Fields() {
			if scope.changeableField(field) {
                //需要保存到当前表的字段，可以查看model_struct.GetModelStruct, 通常为int, string, bool, float等正常的字段
				if field.IsNormal {
                    //scopes.Fields()会为每一个Field设置值，如果这个Field的值是零值"", false, 0，则IsBlank为true
                    //GetModelStruct解析model时，如果字段包含了default或者auto_increment, 则HasDefaultvalue = true
                    //所以当这两个条件满足时，需要在数据保存到数据库后，重新设置这些列的值为默认值或者auto_increment,
                    //详细了解 参考forceReloadAfterCreateCallback
					if field.IsBlank && field.HasDefaultValue {
						blankColumnsWithDefaultValue = append(blankColumnsWithDefaultValue, scope.Quote(field.DBName))
						scope.InstanceSet("gorm:blank_columns_with_default_value", blankColumnsWithDefaultValue)
					} else if !field.IsPrimaryKey || !field.IsBlank {

						columns = append(columns, scope.Quote(field.DBName))
						placeholders = append(placeholders, scope.AddToVars(field.Field.Interface()))
					}
				} else if field.Relationship != nil && field.Relationship.Kind == "belongs_to" {

					for _, foreignKey := range field.Relationship.ForeignDBNames {
						if foreignField, ok := scope.FieldByName(foreignKey); ok && !scope.changeableField(foreignField) {
							columns = append(columns, scope.Quote(foreignField.DBName))
							placeholders = append(placeholders, scope.AddToVars(foreignField.Field.Interface()))
						}
					}
				}
			}
		}
执行sql语句
*/
DefaultCallback.Create().Register("gorm:create", createCallback)
//通过db.Select选择createCallback中设置的从数据库中有默认值，而当前scope中字段值为空的例, 然后赋值给当前scope
/*
if blankColumnsWithDefaultValue, ok := scope.InstanceGet("gorm:blank_columns_with_default_value"); ok {
        //blank_columns_with_default_value中保存的列
		db := scope.DB().New().Table(scope.TableName()).Select(blankColumnsWithDefaultValue.([]string))
		for _, field := range scope.Fields() {
		    //field主键，常见的为id, 它会在createCallback中，设置为lastInsertId
			if field.IsPrimaryKey && !field.IsBlank {
				db = db.Where(fmt.Sprintf("%v = ?", field.DBName), field.Field.Interface())
			}
		}
		db.Scan(scope.Value)
	}
*/
DefaultCallback.Create().Register("gorm:force_reload_after_create", forceReloadAfterCreateCallback)
//它跟saveBeforeAssociationsCallback （主要用于belong_to关系) 相反，saveAfterAssociationsCallback用于has_one, has_many, many_to_many的更新
DefaultCallback.Create().Register("gorm:save_after_associations", saveAfterAssociationsCallback)
//调用model中设置的AfterCreate和AfterSave方法
DefaultCallback.Create().Register("gorm:after_create", afterCreateCallback)
//提交回调，如果有错误发生，则回滚事务
DefaultCallback.Create().Register("gorm:commit_or_rollback_transaction", commitOrRollbackTransactionCallback)
```

