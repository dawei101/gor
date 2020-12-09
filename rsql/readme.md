# rsql
基于 [sqlx](https://github.com/jmoiron/sqlx) 和 [gosql](https://github.com/ilibs/rsql)

## 使用

Quick start!

```go
import (
    "roo.bo/rlib/rsql"
)

func main(){
    rsql.QueryRowx("select * from users where id = 1")
}

```

数据库的连接基于rlib下的封装，默认使用default connection, 如果需要使用其他connection，可以指定

```go
// default db connection
rsql.Queryx("select * from users")

// appcourse db connection
rsql.Use("appcourse").Queryx("select * from users")
```
rsql包裹一层sqlx，所以可以直接使用sqlx的相关 function

使用sqlx的原始方法, 请查看 https://github.com/jmoiron/sqlx

```go
//Exec
rsql.Exec("insert into users(name,email,created_at,updated_at) value(?,?,?,?)","test","test@gmail.com",time.Now(),time.Now())

//Queryx
rows,err := rsql.Queryx("select * from users")
for rows.Next() {
    user := &User{}
    err = rows.StructScan(user)
}
rows.Close()

//QueryRowx
user := &User{}
err := rsql.QueryRowx("select * from users where id = ?",1).StructScan(user)

//Get
user := &User{}
err := rsql.Get(user,"select * from users where id = ?",1)

//Select
users := make([]User)
err := rsql.Select(&users,"select * from users")

//Change database
db := rsql.Use("test")
db.Queryx("select * from tests")
```

你可以使用其他database的connection

```go
rsql.Use("appcourse").Queryx("select * from users")
```

> 因此你所有使用的方法都是基于 connection `apocourse`

## 使用结构体

### Model interface 的定义
```go
type IModel interface {
	TableName() string
	PK() string
}
```

### 结构体demo

```go
type User struct {
	Id        int       `db:"id"`
	Name      string    `db:"name"`
	Email     string    `db:"email"`
	Status    int       `db:"status"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func (u *User) TableName() string {
	return "users"
}

func (u *User) PK() string {
	return "id" // 返回table中的field名
}

//Get
user := &User{}
rsql.Model(user).Where("id=?",1).Get()

//All
user := make([]User,0)
rsql.Model(&user).All()

//Create and auto set CreatedAt
rsql.Model(&User{Name:"test",Email:"test@gmail.com"}).Create()

//Update
rsql.Model(&User{Name:"test2",Email:"test@gmail.com"}).Where("id=?",1).Update()

//Delete
rsql.Model(&User{}).Where("id=?",1).Delete()

```

如果使用 struct 直接当做where的条件查询:

```go
//Get where id = 1 and name = "test1"
user := &User{Id:1,Name:"test1"}
rsql.Model(&user).Get()

//Update 默认使用主键值作为where条件
// update user set id = 1, name = "test2" where id = 1
rsql.Model(&User{Id:1,Name:"test2"}).Update()

//Use custom conditions
//Builder => UPDATE users SET `id`=?,`name`=?,`updated_at`=? WHERE (status = ?)
rsql.Model(&User{Id:1,Name:"test2"}).Where("status = ?",1).Update()

//Delete
rsql.Model(&User{Id:1}).Delete()
```

⚠️  特别注意：结构体中的int default=0的问题处理

```go
user := &User{Id:1,Status:0}
// 如果你需要status加入where筛选中，需要在Get() 指定zero int字段
// where id=1 and status=0
rsql.Model(&user).Get("status")

// where id=1;
rsql.Model(&user).Get()


//如果你想update status=0 你需要在Update() 指定zero int字段
// update user set status=0 where id = 1;
rsql.Model(&User{Status:0}).Where("id=?",1).Update("status")

// 直接使用struct Update时，where 只使用主键筛选
// update user set id = 1, name = "test2" where id = 1
rsql.Model(&User{Id:1,Name:"test2"}).Update()
```

> 从db中生成结构体小工具 [genstruct](https://github.com/fifsky/genstruct)

## 事务
`Tx` 有一个实例方法，如果方法返回error，事务会回滚

```go
rsql.Tx(func(tx *rsql.DB) error {
    for id := 1; id < 10; id++ {
        user := &User{
            Id:    id,
            Name:  "test" + strconv.Itoa(id),
            Email: "test" + strconv.Itoa(id) + "@test.com",
        }
		
		//do some database operations in the transaction (use 'tx' from this point, not 'rsql')
        tx.Model(user).Create()

        if id == 8 {
            return errors.New("interrupt the transaction")
        }
    }

    //query with transaction
    var num int
    err := tx.QueryRowx("select count(*) from user_id = 1").Scan(&num)

    if err != nil {
        return err
    }

    return nil
})
```

带context的事务
```go
rsql.Txx(r.Context(), func(ctx context.Context, tx *rsql.DB) error {
    // do something
})

```

事务也可以使用流程化的 rsql.Begin() / rsql.Use("other").Begin() 完成:
```go
tx, err := rsql.Begin()
if err != nil {
    return err
}

for id := 1; id < 10; id++ {
    _, err := tx.Exec("INSERT INTO users(id,name,status,created_at,updated_at) VALUES(?,?,?,?,?)", id, "test"+strconv.Itoa(id), 1, time.Now(), time.Now())
    if err != nil {
        return tx.Rollback()
    }
}

return tx.Commit()
```

## 创建/更新时间自动添加
如果表结构里有以下字段， 将会自动添加上当前时间

```
AUTO_CREATE_TIME_FIELDS = []string{
    "create_time",
    "create_at",
    "created_at",
    "update_time",
    "update_at",
    "updated_at",
}
AUTO_UPDATE_TIME_FIELDS = []string{
    "update_time",
    "update_at",
    "updated_at",
}
```


## 使用 Map
`Create` `Update` `Delete` `Count` 支持 `map[string]interface`, 举例:

```go
//Create
rsql.Table("users").Create(map[string]interface{}{
    "id":         1,
    "name":       "test",
    "email":      "test@test.com",
    "created_at": "2018-07-11 11:58:21",
    "updated_at": "2018-07-11 11:58:21",
})

//Update
rsql.Table("users").Where("id = ?", 1).Update(map[string]interface{}{
    "name":  "fifsky",
    "email": "fifsky@test.com",
})

//Delete
rsql.Table("users").Where("id = ?", 1).Delete()

//Count
rsql.Table("users").Where("id = ?", 1).Count()

//Change database
rsql.Use("db2").Table("users").Where("id = ?", 1).Count()

//Transaction `tx`
tx.Table("users").Where("id = ?", 1}).Count()
```



## rsql.Expr
怎么使用raw sql，原始语句?
```go
rsql.Table("users").Update(map[string]interface{}{
    "id":2,
    "count":rsql.Expr("count+?",1)
})
//Builder SQL
//UPDATE `users` SET `count`=count + ?,`id`=?; [1 2]
```


## "In" 语句

database/sql 使用in语法，需要自己组合 sql，总是比较复杂, 现在你可以使用sqlx的简单写法：

```go


//SELECT * FROM users WHERE level IN (?);

var levels = []int{4, 6, 7}
rows, err := rsql.Queryx("SELECT * FROM users WHERE level IN (?);", levels)

//or

user := make([]User, 0)
err := rsql.Select(&user, "select * from users where id in(?)",[]int{1,2,3})

```

## Relation
rsql 使用 golang structure 来描述 tables 间的relationships, 你需要使用 `relation` Tag 来定义相关关系, 实例：

⚠️ 跨db的查询，可以直接使用`connection` tag， 来定义对应的connection


```go
type RichMoment struct {
	models.Moment
	User   *models.User    `json:"user" db:"-" relation:"user_id,id"`         //one-to-one
	Photos []models.Photos `json:"photos" db:"-" relation:"id,moment_id" connection:"db2"`     //one-to-many
}
```

单条查询

```go
moment := &RichMoment{}
err := rsql.Model(moment).Where("status = 1 and id = ?",14).Get()
//output User and Photos and you get the result
```

SQL:

```sql
2018/12/06 13:27:54
	Query: SELECT * FROM `moments` WHERE (status = 1 and id = ?);
	Args:  []interface {}{14}
	Time:  0.00300s

2018/12/06 13:27:54
	Query: SELECT * FROM `moment_users` WHERE (id=?);
	Args:  []interface {}{5}
	Time:  0.00081s

2018/12/06 13:27:54
	Query: SELECT * FROM `photos` WHERE (moment_id=?);
	Args:  []interface {}{14}
	Time:  0.00093s
```

多条查询结果, many-to-many

```go
var moments = make([]RichMoment, 0)
err := rsql.Model(&moments).Where("status = 1").Limit(10).All()
//You get the total result  for *UserMoment slice
```

SQL:

```sql
2018/12/06 13:50:59
	Query: SELECT * FROM `moments` WHERE (status = 1) LIMIT 10;
	Time:  0.00319s

2018/12/06 13:50:59
	Query: SELECT * FROM `moment_users` WHERE (id in(?));
	Args:  []interface {}{[]interface {}{5}}
	Time:  0.00094s

2018/12/06 13:50:59
	Query: SELECT * FROM `photos` WHERE (moment_id in(?, ?, ?, ?, ?, ?, ?, ?, ?, ?));
	Args:  []interface {}{[]interface {}{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}}
	Time:  0.00087s
```


Relation中的Where:

```go
moment := &RichMoment{}
err := rsql.Relation("User" , func(b *rsql.Builder) {
    //this is builder instance,
    b.Where("gender = 0")
}).Get(moment , "select * from moments")
```

## 使用钩子-Hooks
在一个Model的执行创建、查询、修改、删除语句时，你可以添加执行前或执行后的方法。

您仅需要为 `Model` 定义相关约定的方法即可，`rsql` 会在相应的操作时，自动调用钩子。如果任意一个钩子方法返回了error，`rsql`会停止接下来的操作，并将当前操作的事务回滚。

```
// begin transaction
BeforeChange
BeforeCreate
// update timestamp `CreatedAt`, `UpdatedAt`
// save
AfterCreate
AfterChange
// commit or rollback transaction
```
Example:

```go
func (u *User) BeforeCreate() (err error) {
  if u.IsValid() {
    err = errors.New("can't save invalid data")
  }
  return
}

func (u *User) AfterCreate(tx *rsql.DB) (err error) {
  if u.Id == 1 {
    u.Email = "after@test.com"
    tx.Model(u).Update()
  }
  return
}
```

> BeforeChange / AfterChange 仅在创建、更新、删除时执行

所有钩子方法Hooks:

```
BeforeChange
AfterChange
BeforeCreate
AfterCreate
BeforeUpdate
AfterUpdate
BeforeDelete
AfterDelete
BeforeFind
AfterFind
```

钩子方法支持多种传参和返回值的组合:

```
func (u *User) BeforeCreate()
func (u *User) BeforeCreate() (err error)
func (u *User) BeforeCreate(tx *rsql.DB)
func (u *User) BeforeCreate(tx *rsql.DB) (err error)
```


