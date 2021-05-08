觉得需要 Wrap 这个 error，抛给上层，因为只有调用者才知道 sql.ErrNoRows 类型的错误需不需要特殊处理，
调用者可能也需要获取执行sql的一些信息
```
findSql :="select name from users where id = ?"
var name string
err = db.QueryRow(findSql, 1).Scan(&name)
if err != nil {
	if err == sql.ErrNoRows {
		return nil,fmt.Errorf("sql %s id:%d is empty: %w", name,id, err)
    }
    //do .... 
}
```