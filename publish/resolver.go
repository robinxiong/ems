package publish

import (
	"database/sql"
	"ems/core/utils"
	"fmt"
	"reflect"
	"strings"

	"github.com/jinzhu/gorm"
)

type resolver struct {
	Records      []interface{}
	Events       []PublishEventInterface
	Dependencies map[string]*dependency //保存要发布的表，以及它关联表的主键(如果关联表也实现了publishInterface)
	DB           *gorm.DB
	publish      Publish //比如调用pulbish.logger用来输出数据库表
}

//dependency 包含了解析Records中的数据类型，以及跟它相并的many2many的数据, 以及主键信息
type dependency struct {
	Type                reflect.Type
	ManyToManyRelations []*gorm.Relationship
	PrimaryValues       [][][]interface{} //保存一个表多行，多例的值，一行的数据，以二维的[][]interface{}表示，每个列以[]interface{}组成
}

func includeValue(value [][]interface{}, values [][][]interface{}) bool {
	for _, v := range values {
		if fmt.Sprint(v) == fmt.Sprint(value) {
			return true
		}
	}
	return false
}

func (resolver *resolver) Publish() (err error) {

	resolver.GenerateDependencies() //找到所有要发布的数据，比如多个产品，每个产品相关的依赖(has-many, has-one, belong-to如果关联表实现了_draft表）， many2many等
	tx := resolver.DB.Begin()
	// Publish Events
	for _, event := range resolver.Events {
		resolver.publish.logger.Print("Publishing Event: ", stringify(event))
		event.Publish(tx)
	}

	// Publish dependencies
	// 将所有需要更新的值，从_draft表的状态为PUBLISHED(false)
	// 从production表中删除old记录
	// 复制_draft表的数据到production表中
	// 发布关联表数据many2many
	for _, dep := range resolver.Dependencies {
		value := reflect.New(dep.Type).Elem()
		productionScope := resolver.DB.Set(publishDraftMode, false).NewScope(value.Addr().Interface())
		productionTable := productionScope.TableName()
		draftTable := DraftTableName(productionTable)
		productionPrimaryKey := scopePrimaryKeys(productionScope, productionTable)
		draftPrimaryKey := scopePrimaryKeys(productionScope, draftTable)

		var columns []string
		for _, field := range productionScope.Fields() {
			if field.IsNormal {
				columns = append(columns, field.DBName)
			}
		}

		var productionColumns, draftColumns []string
		for _, column := range columns {
			productionColumns = append(productionColumns, productionScope.Quote(column))
			draftColumns = append(draftColumns, productionScope.Quote(column))
		}

		if len(dep.PrimaryValues) > 0 {
			queryValues := toQueryValues(dep.PrimaryValues)
			resolver.publish.logger.Print(fmt.Sprintf("Publishing %v: ", productionScope.GetModelStruct().ModelType.Name()), stringifyPrimaryValues(dep.PrimaryValues))

			// set status to published
			updateStateSQL := fmt.Sprintf("UPDATE %v SET publish_status = ? WHERE %v IN (%v)", draftTable, draftPrimaryKey, toQueryMarks(dep.PrimaryValues))

			var params = []interface{}{bool(PUBLISHED)}
			params = append(params, queryValues...)
			tx.Exec(updateStateSQL, params...)

			// delete old records
			deleteSQL := fmt.Sprintf("DELETE FROM %v WHERE %v IN (%v)", productionTable, productionPrimaryKey, toQueryMarks(dep.PrimaryValues))
			tx.Exec(deleteSQL, queryValues...)

			// insert new records
			publishSQL := fmt.Sprintf("INSERT INTO %v (%v) SELECT %v FROM %v WHERE %v IN (%v)",
				productionTable, strings.Join(productionColumns, " ,"), strings.Join(draftColumns, " ,"),
				draftTable, draftPrimaryKey, toQueryMarks(dep.PrimaryValues))
			tx.Exec(publishSQL, queryValues...)

			// publish join table data
			for _, relationship := range dep.ManyToManyRelations {
				productionTable := relationship.JoinTableHandler.Table(tx.Set(publishDraftMode, false))
				draftTable := relationship.JoinTableHandler.Table(tx.Set(publishDraftMode, true))
				var productionJoinKeys, draftJoinKeys []string
				var productionCondition, draftCondition string
				for _, foreignKey := range relationship.JoinTableHandler.SourceForeignKeys() {
					productionJoinKeys = append(productionJoinKeys, fmt.Sprintf("%v.%v", productionTable, productionScope.Quote(foreignKey.DBName)))
					draftJoinKeys = append(draftJoinKeys, fmt.Sprintf("%v.%v", draftTable, productionScope.Quote(foreignKey.DBName)))
				}

				if len(productionJoinKeys) > 1 {
					productionCondition = fmt.Sprintf("(%v)", strings.Join(productionJoinKeys, ","))
					draftCondition = fmt.Sprintf("(%v)", strings.Join(draftJoinKeys, ","))
				} else {
					productionCondition = strings.Join(productionJoinKeys, ",")
					draftCondition = strings.Join(draftJoinKeys, ",")
				}

				sql := fmt.Sprintf("DELETE FROM %v WHERE %v IN (%v)", productionTable, productionCondition, toQueryMarks(dep.PrimaryValues, relationship.ForeignFieldNames...))
				tx.Exec(sql, toQueryValues(dep.PrimaryValues, relationship.ForeignFieldNames...)...)

				rows, _ := tx.Raw(fmt.Sprintf("SELECT * FROM %s", draftTable)).Rows()
				joinColumns, _ := rows.Columns()
				rows.Close()
				if len(joinColumns) == 0 {
					continue
				}

				var productionJoinTableColumns, draftJoinTableColumns []string
				for _, column := range joinColumns {
					productionJoinTableColumns = append(productionJoinTableColumns, productionScope.Quote(column))
					draftJoinTableColumns = append(draftJoinTableColumns, productionScope.Quote(column))
				}

				publishSQL := fmt.Sprintf("INSERT INTO %v (%v) SELECT %v FROM %v WHERE %v IN (%v)",
					productionTable, strings.Join(productionJoinTableColumns, " ,"), strings.Join(draftJoinTableColumns, " ,"),
					draftTable, draftCondition, toQueryMarks(dep.PrimaryValues, relationship.ForeignFieldNames...))
				tx.Exec(publishSQL, toQueryValues(dep.PrimaryValues, relationship.ForeignFieldNames...)...)
			}
		}
	}

	if err = tx.Error; err == nil {
		return tx.Commit().Error
	}

	tx.Rollback()
	return err

}

func (resolver *resolver) Discard() (err error) {
	resolver.GenerateDependencies()
	tx := resolver.DB.Begin()

	// Discard Events
	for _, event := range resolver.Events {
		resolver.publish.logger.Print("Discarding Event: ", stringify(event))
		event.Discard(tx)
	}

	// Discard dependencies
	// 删除_draft表中的记录
	// 将production表中当前的记录复制到_draft

	for _, dep := range resolver.Dependencies {
		value := reflect.New(dep.Type).Elem()
		productionScope := resolver.DB.Set(publishDraftMode, false).NewScope(value.Addr().Interface())
		productionTable := productionScope.TableName()
		draftTable := DraftTableName(productionTable)

		productionPrimaryKey := scopePrimaryKeys(productionScope, productionTable)
		draftPrimaryKey := scopePrimaryKeys(productionScope, draftTable)

		var columns []string
		for _, field := range productionScope.Fields() {
			if field.IsNormal {
				columns = append(columns, field.DBName)
			}
		}

		var productionColumns, draftColumns []string
		for _, column := range columns {
			productionColumns = append(productionColumns, productionScope.Quote(column))
			draftColumns = append(draftColumns, productionScope.Quote(column))
		}

		if len(dep.PrimaryValues) > 0 {
			resolver.publish.logger.Print(fmt.Sprintf("Discarding %v: ", productionScope.GetModelStruct().ModelType.Name()), stringifyPrimaryValues(dep.PrimaryValues))

			// delete data from draft db
			deleteSQL := fmt.Sprintf("DELETE FROM %v WHERE %v IN (%v)", draftTable, draftPrimaryKey, toQueryMarks(dep.PrimaryValues))
			tx.Exec(deleteSQL, toQueryValues(dep.PrimaryValues)...)

			// delete join table
			for _, relationship := range dep.ManyToManyRelations {
				productionTable := relationship.JoinTableHandler.Table(tx.Set(publishDraftMode, false))
				draftTable := relationship.JoinTableHandler.Table(tx.Set(publishDraftMode, true))

				var productionJoinKeys, draftJoinKeys []string
				var productionCondition, draftCondition string
				for _, foreignKey := range relationship.JoinTableHandler.SourceForeignKeys() {
					productionJoinKeys = append(productionJoinKeys, fmt.Sprintf("%v.%v", productionTable, productionScope.Quote(foreignKey.DBName)))
					draftJoinKeys = append(draftJoinKeys, fmt.Sprintf("%v.%v", draftTable, productionScope.Quote(foreignKey.DBName)))
				}

				if len(productionJoinKeys) > 1 {
					productionCondition = fmt.Sprintf("(%v)", strings.Join(productionJoinKeys, ","))
					draftCondition = fmt.Sprintf("(%v)", strings.Join(draftJoinKeys, ","))
				} else {
					productionCondition = strings.Join(productionJoinKeys, ",")
					draftCondition = strings.Join(draftJoinKeys, ",")
				}

				sql := fmt.Sprintf("DELETE FROM %v WHERE %v IN (%v)", draftTable, draftCondition, toQueryMarks(dep.PrimaryValues, relationship.ForeignFieldNames...))
				tx.Exec(sql, toQueryValues(dep.PrimaryValues, relationship.ForeignFieldNames...)...)

				rows, _ := tx.Raw(fmt.Sprintf("select * from %v", draftTable)).Rows()
				joinColumns, _ := rows.Columns()
				rows.Close()
				if len(joinColumns) == 0 {
					continue
				}
				var productionJoinTableColumns, draftJoinTableColumns []string
				for _, column := range joinColumns {
					productionJoinTableColumns = append(productionJoinTableColumns, productionScope.Quote(column))
					draftJoinTableColumns = append(draftJoinTableColumns, productionScope.Quote(column))
				}

				publishSQL := fmt.Sprintf("INSERT INTO %v (%v) SELECT %v FROM %v WHERE %v IN (%v)",
					draftTable, strings.Join(draftJoinTableColumns, " ,"), strings.Join(productionJoinTableColumns, " ,"),
					productionTable, productionCondition, toQueryMarks(dep.PrimaryValues, relationship.ForeignFieldNames...))
				tx.Exec(publishSQL, toQueryValues(dep.PrimaryValues, relationship.ForeignFieldNames...)...)
			}

			// copy data from production to draft
			discardSQL := fmt.Sprintf("INSERT INTO %v (%v) SELECT %v FROM %v WHERE %v IN (%v)",
				draftTable, strings.Join(draftColumns, " ,"),
				strings.Join(productionColumns, " ,"), productionTable,
				productionPrimaryKey, toQueryMarks(dep.PrimaryValues))
			tx.Exec(discardSQL, toQueryValues(dep.PrimaryValues)...)
		}
	}

	if err = tx.Error; err == nil {
		return tx.Commit().Error
	}
	tx.Rollback()
	return err
}

func (resolver *resolver) GetDependencies(dep *dependency, primaryKeys [][][]interface{}) {
	value := reflect.New(dep.Type)
	fromScope := resolver.DB.NewScope(value.Interface())
	//返因一个新的db
	draftDB := resolver.DB.Set(publishDraftMode, true).Unscoped()

	//读书model字段的关系信息， 同时这个字段，也实现了publishInterface
	for _, field := range fromScope.Fields() {
		if relationship := field.Relationship; relationship != nil {
			if IsPublishableModel(field.Field.Interface()) {
				toType := utils.ModelType(field.Field.Interface())
				//查找_draft关联表, hasMany, hasOne, many2many
				toScope := draftDB.NewScope(reflect.New(toType).Interface())
				draftTable := DraftTableName(toScope.TableName())

				//保存找到的相关表的主键及值
				var dependencyKeys [][][]interface{}
				var rows *sql.Rows
				var err error
				var selectPrimaryKeys []string

				//找到关联表的主键
				for _, field := range toScope.PrimaryFields() {
					selectPrimaryKeys = append(selectPrimaryKeys, toScope.Quote(field.DBName))
				}

				//根据primaryKeys指定的主键和值，去检索依赖它的表的主键

				if relationship.Kind == "has_one" || relationship.Kind == "has_many" {
					//生成where条件, 返回 子表的外键是存在于当前表中
					//Product{Id Tag []Tag}
					//Tag {ProductId int, Name string}
					//toQueryCondition 返回外键的名称(ProductId, LanguageId) toQueryMarks 返回占位符(?, ?), toQueryValues返回primaryKeys中的值信息
					sql := fmt.Sprintf("%v IN (%v)", toQueryCondition(toScope, relationship.ForeignDBNames), toQueryMarks(primaryKeys, relationship.AssociationForeignDBNames...))

					rows, err = draftDB.Table(draftTable).Select(selectPrimaryKeys).Where("publish_status=?", DIRTY).Where(sql, toQueryValues(primaryKeys, relationship.AssociationForeignDBNames...)...).Rows()
				} else if relationship.Kind == "belongs_to" {
					//查找当前表信赖的表的主键值
					fromTable := DraftTableName(fromScope.TableName())
					// toTable := toScope.TableName()

					sql := fmt.Sprintf("%v IN (SELECT %v FROM %v WHERE %v IN (%v))",
						strings.Join(relationship.AssociationForeignDBNames, ","), strings.Join(relationship.ForeignDBNames, ","), fromTable, scopePrimaryKeys(fromScope, fromTable), toQueryMarks(primaryKeys))

					rows, err = draftDB.Table(draftTable).Select(selectPrimaryKeys).Where("publish_status = ?", DIRTY).Where(sql, toQueryValues(primaryKeys)...).Rows()

				}

				if rows != nil && err == nil {
					defer rows.Close()
					columns, _ := rows.Columns()
					for rows.Next() {
						//创建一个数据，用来保存找到的相关表的主键值
						var primaryValues = make([]interface{}, len(columns))
						for idx := range primaryValues {
							var value interface{}
							primaryValues[idx] = &value
						}
						rows.Scan(primaryValues...)

						var currentDependencyKeys [][]interface{}
						for idx, value := range primaryValues {
							currentDependencyKeys = append(currentDependencyKeys, []interface{}{columns[idx], reflect.Indirect(reflect.ValueOf(value)).Interface()})
						}

						dependencyKeys = append(dependencyKeys, currentDependencyKeys)
					}

					resolver.AddDependency(&dependency{Type: toType, PrimaryValues: dependencyKeys})
				}
			}

			if relationship.Kind == "many_to_many" {
				dep.ManyToManyRelations = append(dep.ManyToManyRelations, relationship)
			}
		}
	}
}

func (resolver *resolver) AddDependency(dep *dependency) {
	name := dep.Type.String()
	var newPrimaryKes [][][]interface{}
	//如果安前的数据类型缓存在的resolver.Dependencies中, 即之前解析过, 但是要保存新的值比如Product{id:1}, Product{id:2}
	if d, ok := resolver.Dependencies[name]; ok {
		for _, primaryKey := range dep.PrimaryValues {
			if !includeValue(primaryKey, d.PrimaryValues) {
				newPrimaryKes = append(newPrimaryKes, primaryKey)
				d.PrimaryValues = append(d.PrimaryValues, primaryKey)
			}
		}
	} else {
		resolver.Dependencies[name] = dep
		newPrimaryKes = dep.PrimaryValues
	}

	if len(newPrimaryKes) > 0 {
		resolver.GetDependencies(dep, newPrimaryKes)
	}
}

//GenerateDependencies 遍历要发布的记录Records, 解析每一个record的实际值，主键信息，many2many信息，event信息
// 如果这个 record has-one, has-many相关的关联表也实现了publishInterface,  则也保存到resolver.Dependencies中
// 数据的详细信息保存到Dependencies中，事件信息则保存到Events中
func (resolver *resolver) GenerateDependencies() {

	var addToDependencies = func(data interface{}) {
		//如果数据实现了publishInterface接口，即拥有publish_status数据列
		if IsPublishableModel(data) {
			//db中包含了publish draft或者publish等设置
			scope := resolver.DB.NewScope(data)
			//主键信息中第一个值为列名，第二个为实际值, 以2维数据保存
			var primaryValues [][]interface{}
			for _, field := range scope.PrimaryFields() {
				primaryValues = append(primaryValues, []interface{}{field.DBName, field.Field.Interface()})
			}
			resolver.AddDependency(&dependency{Type: utils.ModelType(data), PrimaryValues: [][][]interface{}{primaryValues}})
		}

		if event, ok := data.(PublishEventInterface); ok {
			resolver.Events = append(resolver.Events, event)
		}
	}

	for _, record := range resolver.Records {
		//value为指针，则获取引用的值
		reflectValue := reflect.Indirect(reflect.ValueOf(record))
		if reflectValue.Kind() == reflect.Slice {
			for i := 0; i < reflectValue.Len(); i++ {
				//以interface{}返回值
				addToDependencies(reflectValue.Index(i).Interface())
			}

		} else {
			addToDependencies(record)
		}
	}
}

//
func toQueryCondition(scope *gorm.Scope, columns []string) string {
	var newColumns []string
	for _, column := range columns {
		newColumns = append(newColumns, scope.Quote(column))
	}

	if len(columns) > 1 {
		return fmt.Sprintf("(%v)", strings.Join(newColumns, ","))
	}
	return strings.Join(columns, ",")
}

//传入主键值，它通常是[]interface{key, value}, 比如[][][]interface{[][]interface{[]interface{id:2}}}
//参照的列名
func toQueryMarks(primaryValues [][][]interface{}, columns ...string) string {
	var results []string
	for _, primaryValue := range primaryValues {
		//二维数组，一个数所据行所有主键以及值
		var marks []string
		for _, value := range primaryValue {
			//当个主键和值
			//如果参照的表，主键为0
			if len(columns) == 0 {
				marks = append(marks, "?")
			} else {
				for _, column := range columns {
					//当依赖中的主键主参考表的主键相同时, 添加 ？占位符
					if fmt.Sprintf("%v", value[0]) == fmt.Sprintf("%v", column) {
						marks = append(marks, "?")
					}
				}
			}
		}

		if len(marks) > 1 {
			results = append(results, fmt.Sprintf("(%v)", strings.Join(marks, ",")))
		} else {
			results = append(results, strings.Join(marks, ""))
		}
	}
	return strings.Join(results, ",")
}

//返回primaryValues中的值
func toQueryValues(primaryValues [][][]interface{}, columns ...string) (values []interface{}) {
	for _, primaryValue := range primaryValues {
		for _, value := range primaryValue {
			if len(columns) == 0 {
				values = append(values, value[1])
			} else {
				for _, column := range columns {
					if column == fmt.Sprint(value[0]) {
						values = append(values, value[1])
					}
				}
			}
		}
	}
	return values
}

func scopePrimaryKeys(scope *gorm.Scope, tableName string) string {
	var primaryKeys []string
	for _, field := range scope.PrimaryFields() {
		key := fmt.Sprintf("%v.%v", scope.Quote(tableName), scope.Quote(field.DBName))
		primaryKeys = append(primaryKeys, key)
	}
	if len(primaryKeys) > 1 {
		return fmt.Sprintf("(%v)", strings.Join(primaryKeys, ","))
	}
	return strings.Join(primaryKeys, "")
}
