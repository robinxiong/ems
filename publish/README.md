# publish
1. 创建表db.DropTable， 如果db是draft模式，而table没有实现publishInterface, 则不会删除a_draft表，而是删除a表
2. 创建一条记录，不管是draft 或者production, 都是先将数据写入到_draft表中, 查看setTableAndPublishStatus
   如果是production, 在提交之前，还需要在production表中创建一个记录
   对于has-one和belong-to的关系，不会保存到draft表中，只保存到非_draft表，所以这些表可以不需要生成相应的drfat表 
   如果是production模式，则不会在draft表中的publish_status设置为true(DIRTY)
   如果是draft模式, 则设置publish_status为true
3. 删除记录 
   如果是draft模式，则直接安全删除draft表中的数据，同时设置publish_status为DIRTY
   如果是publish模式，则安全删除publish表中的数据和_draft表中的数据，但publish_status都依然为FALSE
    
4. 更新
   如果是draft模式，则_draft表中的数据得到更新,production表中的数据不会被更新，同时设置_draft表的publish_status为DIRTY
   如果是publish模式，则_draft表中的数据会更新，但状态依然为PUBLISHED(false), 同时production 表也得到相应的更新

5. 发布事件
    createResourcePublishInterface 结合具体的model struct, 比如product struct
    PublishEvent 是一个通用的类，所对应的数据库表为publish_event, 通过event.go中的注册createResourcePublishInterface,将
    publishevent的name与createResourcePublishInterface结合，在调用pb.Publish(&product{})
    
BUG: resolver.go rows, err = draftDB.Table(draftTable).Select(selectPrimaryKeys).Where("publish_status=?", DIRTY)
在db中设置了publish:publish_event时, publish_statue应该为false, 而不可能是DIRTY, 所以在PUBLISHEvent发布时，它不能找到DIRTY的记录    
其实对于publish:publish_event的使用，可以参考sorting.go move方法，它只是用于在publish_event表中创建一个事件记录，而不是用于创建一条draft记录，所以它不需要修改publish_status的状态
比如重新调整了行的排序，所以创建一条记录时，不能设置为publish:publish_event, 否则找不到关联系的数据。而对于修改sort的时候，我们同样不需要修改_draft的publish_status时，
我们可以设置publish:publish_event, 所以db.Set("publish:publish_event", true)， 必须跟PublishEvent相结合使用
