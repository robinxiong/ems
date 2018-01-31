回调
=========
1. 对于创建一个新的行，则它会在原先存在行的position+1
2. 查询sortingInterface行时，修改sql为position desc或者position asc
3. 删除行时，需要对原先的记录重新进行排序，特别是对于locale和publish的处理