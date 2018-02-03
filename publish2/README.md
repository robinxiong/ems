#publish2
提供了带版本的草稿，可见性，以及任务(预设发布时间)
 
```go
type Product struct {
  ...

  publish2.Version
  publish2.Schedule
  publish2.Visible
}
//e.g.
admin.New(&qor.Config{DB: db.DB.Set(publish2.VisibleMode, publish2.ModeOff).Set(publish2.ScheduleMode, publish2.ModeOff)})

```