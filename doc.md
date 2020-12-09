# rlib
--
    import "roo.bo/rlib"

rlib: Make life super easy。

轮子很多，好胎不多。造好胎、用好胎、开好车，做卓越工程师。

相关功能索引

    1. config
    2. db

## Usage

```go
const (
	ERROR_SYSTEM          = -1
	ERROR_REMOTE_SERVER   = -2
	ERROR_RESPONSE_FORMAT = -3
)
```

```go
const RedisNil = redis.Nil
```

#### func  DB

```go
func DB(name string) *sql.DB
```
获得 name 的 *sql.DB

获得前必须保证 InitDB 过，否则会 panic

#### func  DBDefault

```go
func DBDefault() *sql.DB
```
获得 name 的 *sql.DB

获得前必须保证 `InitDB("default", "xxxx")` 过，否则会 panic

#### func  DBReg

```go
func DBReg(name, dbType string, dataSource string, maxOpenConns int) error
```
注册DB

获得前必须保证 InitDB 过，否则会 panic

#### func  DBX

```go
func DBX(name string) *sqlx.DB
```
获得*sqlx.DB

获得前必须保证 InitDB 过，否则会 panic

请不要使用migration特性

@see github.com/jmoiron/sqlx

#### func  DBXDefault

```go
func DBXDefault() *sqlx.DB
```
获得default *sqlx.DB

获得前必须保证 InitDB 过，否则会 panic

请不要使用migration特性

@see github.com/jmoiron/sqlx

#### func  GetHTTPClient

```go
func GetHTTPClient(service string) *http.Client
```

#### func  HttpServe

```go
func HttpServe(router *mux.Router, serveAt string) error
```

#### func  InitHTTPClient

```go
func InitHTTPClient(service string, config *HTTPClientConfig) error
```

#### func  JsonBodyTo

```go
func JsonBodyTo(r *http.Request, v interface{}) error
```

#### func  KafkaConsume

```go
func KafkaConsume(name string, consume func(context.Context, *kafka.Message))
```

#### func  KafkaProduce

```go
func KafkaProduce(name string, msg *KafkaMessage)
```

#### func  KafkaProducerGet

```go
func KafkaProducerGet(name string) *kafka.Producer
```

#### func  KafkaProducerStop

```go
func KafkaProducerStop(name string)
```

#### func  LocalCacheGet

```go
func LocalCacheGet(key string) (value interface{}, ok bool)
```
获取进程内存级缓存

#### func  LocalCacheSet

```go
func LocalCacheSet(key string, value interface{}, d time.Duration)
```
设置进程内存级缓存

#### func  LogDisableCodeinfo

```go
func LogDisableCodeinfo()
```
关闭显示代码与行信息

#### func  LogSetLevel

```go
func LogSetLevel(level LogLevel)
```
设置日志打印的级别

#### func  LogSetLocation

```go
func LogSetLocation(location *time.Location)
```
设置日志打印的时区信息

#### func  MicroTimestamp

```go
func MicroTimestamp(t time.Time) int64
```

#### func  MockForAPI

```go
func MockForAPI(service string, path string, responseBody interface{})
```
MockForApi mock single api

    	example:
    	 service string "baby"
      path  string  URL.Path  "/user-baby/add"
    	 responseBody string   `{"result":0,"msg":"mocked","data":{"items":"mocked data"}}`

#### func  MockForService

```go
func MockForService(service string) error
```
MockForService mock apis of service from config

#### func  NewMultiResponseWriter

```go
func NewMultiResponseWriter(w http.ResponseWriter, buf io.Writer) *multiResponseWriter
```

#### func  Redis

```go
func Redis(name string) *redis.Client
```

#### func  RedisDefault

```go
func RedisDefault() *redis.Client
```

#### func  RedisReg

```go
func RedisReg(name, addr, password string, db int) error
```

#### func  RequestId

```go
func RequestId(ctx context.Context) string
```

#### func  SetConfigFile

```go
func SetConfigFile(fn string)
```

#### func  SetDefaultLog

```go
func SetDefaultLog(log *RLog)
```
设置默认的Log

设置后，可以直接用 `rlib.Debug` 等打印日志

#### func  Timestamp

```go
func Timestamp(t time.Time) int64
```

#### func  ValidateField

```go
func ValidateField(v interface{}) error
```
验证实例字段 use for example:

    import rlib
    user = &struct {
    	Appid string `validate:"required,appid"`
    	SN    string `validate:"required,sn"`
    }{
    	SN:    "102102000100008F",
    	Appid: "TI1ZDE1NzJhZTEwO",
    }

    if err := rlib.FieldValid(user); err != nil {
    	fmt.Println(err)
    	rlib.NewErrResp(-422, err.Error(), "")
    }

#### func  ValidatorReg

```go
func ValidatorReg(tag string, f valid.Func) error
```
运行时动态加入验证规则

#### type ApiConfig

```go
type ApiConfig struct {
	BaseUrl      string                     `json:"baseUrl" yaml:"baseUrl"`
	Retry        int                        `json:"retry" yaml:"retry"`
	Timeout      int                        `json:"timeout" yaml:"timeout"`
	Supertoken   string                     `json:"supertoken" yaml:"supertoken"`
	SpecialRules map[string]*ApiSpecialRule `json:"specialRules" yaml:"SpecialRules"`
	Client       *HTTPClientConfig          `json:"client" yaml:"client"`
}
```


#### type ApiSpecialRule

```go
type ApiSpecialRule struct {
	Retry   int  `json:"retry" yaml:"retry"`
	Timeout int  `json:"timeout" yaml:"timeout"`
	Logable bool `json:"logable" yaml:"logable"`
}
```


#### type Config

```go
type Config struct {
	*Struct
}
```


#### func  DefaultConfig

```go
func DefaultConfig() *Config
```
获取默认配置

如果未通过`LoadDefaultConfig(filePath)` 或 `LoadConfig(name, filePath)` 加载过配置, 将会panic

    LoadDefaultConfig("./current/path/config.yml")
    GetDefaultConfig()

一般情况下，我们只需要使用 default 这套config就足够用

#### func  GetConfig

```go
func GetConfig(name string) *Config
```
根据配置名获取配置

    LoadConfig("myconfig", "./current/path/config.yml")
    GetConfig("myconfig")

#### func  RegConfig

```go
func RegConfig(name, filePath string) *Config
```
加载配置文件，并按name标识起来

#### func  RegDefaultConfig

```go
func RegDefaultConfig(filePath string) *Config
```
装在默认的配置文件，必须在使用前加载配置

#### type ContextLogFunc

```go
type ContextLogFunc func(ctx context.Context, v ...interface{})
```


```go
var (
	Debug   ContextLogFunc
	Info    ContextLogFunc
	Warning ContextLogFunc
	Error   ContextLogFunc
)
```
日志快捷方法: `Debug` `Info` `Warning` `Error`

要切换日志输出目录，使用 `SetDefaultLog()` 设置新的RLog实例

#### type DBConfig

```go
type DBConfig struct {
	DataSource   string `json:"dataSource" yaml:"dataSource"`
	MaxOpenConns int    `json:"maxOpenConns" yaml:"MaxOpenConns"`
	DBType       string `json:"dbType" yaml:"DBType"`
}
```


#### type Element

```go
type Element struct {
	Key, Value interface{}
}
```


#### func (*Element) Next

```go
func (e *Element) Next() *Element
```
Next returns the next element, or nil if it finished.

#### func (*Element) Prev

```go
func (e *Element) Prev() *Element
```
Prev returns the previous element, or nil if it finished.

#### type HTTPClient

```go
type HTTPClient struct {
}
```

HTTPClient struct

#### func  NewRooboHTTPClient

```go
func NewRooboHTTPClient(ctx context.Context, topic string, option *HTTPClientOption, client *http.Client) *HTTPClient
```
NewHRooboTTPClient method

#### func (HTTPClient) Get

```go
func (c HTTPClient) Get(urlstr string, body interface{}, header map[string]string) (*http.Response, error)
```
Get method body 内容可以为一下数据结构 url.Values,map[string]string, map[string][string]

#### func (HTTPClient) Post

```go
func (c HTTPClient) Post(urlstr string, body interface{}, header map[string]string) (*http.Response, error)
```
Post method body 如果传入内容 url.Values,map[string]string, map[string][string]
那么将按照Content-Type=application/x-www-form-urlencoded 进行请求
如果body实体中含有文件，那么Content-Type=multipart/form-data 传输 如果body 实体为
string,[]byte,io.Reader 那么直接进行post请求,额外的需要自定义content-type

#### func (HTTPClient) PostJSON

```go
func (c HTTPClient) PostJSON(urlstr string, data interface{}, header map[string]string) (*http.Response, error)
```
PostJSON method

#### func (HTTPClient) PostMultipart

```go
func (c HTTPClient) PostMultipart(urlstr string, body interface{}, header map[string]string) (*http.Response, error)
```
PostMultipart method

#### type HTTPClientConfig

```go
type HTTPClientConfig struct {
	Dialer                  *HTTPDialer `json:"dialer" yaml:"dialer"`
	MaxIdleConns            int         `json:"maxIdleConns" yaml:"maxIdleConns"`
	MaxIdleConnsPerHost     int         `json:"maxIdleConnsPerHost" yaml:"maxIdleConnsPerHost"`
	IdleConnTimeoutMs       int         `json:"idleConnTimeoutMs" yaml:"idleConnTimeoutMs"`
	ExpectContinueTimeoutMs int         `json:"expectContinueTimeoutMs" yaml:"expectContinueTimeoutMs"`
	ResponseHeaderTimeoutMs int         `json:"responseHeaderTimeoutMs" yaml:"responseHeaderTimeoutMs"`
}
```


#### type HTTPClientOption

```go
type HTTPClientOption struct {
	Retry           int
	RetryInterval   time.Duration
	CloseLog        bool
	UnableRequestID bool
	Timeout         time.Duration
	Mock            string
}
```

HTTPClientOption 配置选项

#### type HTTPDialer

```go
type HTTPDialer struct {
	TimeoutMs   int  `json:"timeoutMs" yaml:"timeoutMs"`
	KeepaliveMs int  `json:"keepaliveMs" yaml:"keepaliveMs"`
	DualStack   bool `json:"dualStack" yaml:"dualStack"`
}
```


#### type KafkaConsumerConfig

```go
type KafkaConsumerConfig struct {
	Servers string   `json:"servers" yaml:"servers"`
	Group   string   `json:"group" yaml:"group"`
	Offset  string   `json:"offset" yaml:"offset"`
	Topics  []string `json:"topics" yaml:"topics"`
}
```


#### type KafkaMessage

```go
type KafkaMessage struct {
	Action   string                 `json:"action"`
	ClientId string                 `json:"clientId"`
	Data     map[string]interface{} `json:"data"`
}
```


#### func (*KafkaMessage) DataAssignTo

```go
func (m *KafkaMessage) DataAssignTo(data interface{})
```

#### func (*KafkaMessage) DataValue

```go
func (m *KafkaMessage) DataValue(key string) interface{}
```

#### func (*KafkaMessage) SetData

```go
func (m *KafkaMessage) SetData(data map[string]interface{})
```

#### type KafkaProducerConfig

```go
type KafkaProducerConfig struct {
	Servers string `json:"servers" yaml:"servers"`
	Topic   string `json:"topic" yaml:"topic"`
}
```


#### type LogConfig

```go
type LogConfig struct {
	Path       string `json:"path" yaml:"path"`
	Level      string `json:"level" yaml:"level"`
	MaxMB      int    `json:"maxMB" yaml:"maxDB"`
	MaxDays    int    `json:"maxDays" yaml:"MaxDays"`
	MaxBackups int    `json:"maxBackups" yaml:"MaxBackups"`
}
```


#### type LogLevel

```go
type LogLevel int
```


```go
const (
	LEVEL_DEBUG LogLevel = iota - 1
	LEVEL_INFO
	LEVEL_WARNING
	LEVEL_ERROR
)
```

#### func  LogLevelFromString

```go
func LogLevelFromString(level string) LogLevel
```
string型日志级别快速转换为 `LogLevel`

#### type MockConfig

```go
type MockConfig struct {
	Root         string `json:"root" yaml:"root"`
	FilePaths    []string
	ServicePaths map[string][]string
}
```

MockConfig config mock path

#### type MockFunc

```go
type MockFunc func(r *http.Request) *http.Response
```

MockFunc implements to http.Transport

#### func (MockFunc) RoundTrip

```go
func (f MockFunc) RoundTrip(r *http.Request) (*http.Response, error)
```
RoundTrip http.RoundTrip

#### type OrderedMap

```go
type OrderedMap struct {
}
```


#### func  NewOrderedMap

```go
func NewOrderedMap() *OrderedMap
```

#### func (*OrderedMap) Back

```go
func (m *OrderedMap) Back() *Element
```
Back will return the element that is the last (most recent Set element). If
there are no elements this will return nil.

#### func (*OrderedMap) Delete

```go
func (m *OrderedMap) Delete(key interface{}) (didDelete bool)
```
Delete will remove a key from the map. It will return true if the key was
removed (the key did exist).

#### func (*OrderedMap) Front

```go
func (m *OrderedMap) Front() *Element
```
Front will return the element that is the first (oldest Set element). If there
are no elements this will return nil.

#### func (*OrderedMap) Get

```go
func (m *OrderedMap) Get(key interface{}) (interface{}, bool)
```
Get returns the value for a key. If the key does not exist, the second return
parameter will be false and the value will be nil.

#### func (*OrderedMap) GetElement

```go
func (m *OrderedMap) GetElement(key interface{}) *Element
```
GetElement returns the element for a key. If the key does not exist, the pointer
will be nil.

#### func (*OrderedMap) GetOrDefault

```go
func (m *OrderedMap) GetOrDefault(key, defaultValue interface{}) interface{}
```
GetOrDefault returns the value for a key. If the key does not exist, returns the
default value instead.

#### func (*OrderedMap) Keys

```go
func (m *OrderedMap) Keys() (keys []interface{})
```
Keys returns all of the keys in the order they were inserted. If a key was
replaced it will retain the same position. To ensure most recently set keys are
always at the end you must always Delete before Set.

#### func (*OrderedMap) Len

```go
func (m *OrderedMap) Len() int
```
Len returns the number of elements in the map.

#### func (*OrderedMap) Set

```go
func (m *OrderedMap) Set(key, value interface{}) bool
```
Set will set (or replace) a value for a key. If the key was new, then true will
be returned. The returned value will be false if the value was replaced (even if
the value was the same).

#### type Page

```go
type Page struct {
	Items      []interface{} `json:"items"`
	Pagination Pagination    `json:"pagination"`
}
```


#### func  PageIt

```go
func PageIt(items interface{}, page, pageSize int, total int) *Page
```

#### type Pagination

```go
type Pagination struct {
	Page     int `json:"page"`
	PageSize int `json:"pageSize"`
	Total    int `json:"total"`
}
```


#### func  PaginationFromRequest

```go
func PaginationFromRequest(r *http.Request) *Pagination
```

#### type RLog

```go
type RLog struct {
	Logfile string
	Debug   ContextLogFunc
	Info    ContextLogFunc
	Warning ContextLogFunc
	Error   ContextLogFunc
}
```


#### func  NewFileLog

```go
func NewFileLog(file string, loglv LogLevel, maxMB int, maxdays int, maxbackups int) *RLog
```

#### type RedisConfig

```go
type RedisConfig struct {
	Addr     string `json:"addr" yaml:"addr"`
	Password string `json:"password" yaml:"password"`
	DB       int    `json:"db" yaml:"db"`
}
```


#### type ResponseError

```go
type ResponseError struct {
	Code int
	Msg  string
	Desc string
}
```


#### func  NewResponseError

```go
func NewResponseError(code int, msg, desc string) ResponseError
```

#### func (ResponseError) Error

```go
func (e ResponseError) Error() string
```

#### func (ResponseError) Flush

```go
func (e ResponseError) Flush(w http.ResponseWriter)
```

#### type RooboConfig

```go
type RooboConfig struct {
	*Config
	DevMode       bool                            `json:"devMode" yaml:"devMode"`
	Log           *LogConfig                      `json:"log" yaml:"log"`
	DB            map[string]*DBConfig            `json:"db" yaml:"db"`
	Redis         map[string]*RedisConfig         `json:"redis" yaml:"redis"`
	Api           map[string]*ApiConfig           `json:"api" yaml:"api"`
	Mock          *MockConfig                     `json:"mock" yaml:"mock"`
	KafkaConsumer map[string]*KafkaConsumerConfig `json:"kafka_consumer" yaml:"kafka_consumer"`
	KafkaProducer map[string]*KafkaProducerConfig `json:"kafka_producer" yaml:"kafka_producer"`
}
```


#### func  DefaultRooboConfig

```go
func DefaultRooboConfig() *RooboConfig
```

#### func (*RooboConfig) Assign

```go
func (c *RooboConfig) Assign(key string, config interface{})
```

#### type RooboResponse

```go
type RooboResponse struct {
	Result int         `json:"result"`
	Msg    string      `json:"msg"`
	Desc   string      `json:"desc"`
	Data   interface{} `json:"data"`
}
```


#### func  NewErrResp

```go
func NewErrResp(result int, msg string, desc string) *RooboResponse
```

#### func  NewResp

```go
func NewResp(data interface{}) *RooboResponse
```

#### func  RooboPost

```go
func RooboPost(ctx context.Context, service, path string, data map[string]interface{}, header map[string]string) (err error, rbres *RooboResponse)
```

#### func (*RooboResponse) DataAssignTo

```go
func (rr *RooboResponse) DataAssignTo(data interface{})
```

#### func (*RooboResponse) Flush

```go
func (res *RooboResponse) Flush(w http.ResponseWriter) error
```

#### type Struct

```go
type Struct struct {
	Raw map[string]interface{}
}
```


#### func  JsonBody

```go
func JsonBody(r *http.Request) (st *Struct, err error)
```

#### func  NewStruct

```go
func NewStruct(data map[string]interface{}) *Struct
```

#### func (*Struct) DataAssignTo

```go
func (r *Struct) DataAssignTo(val interface{})
```

#### func (*Struct) Get

```go
func (r *Struct) Get(key string) (*Struct, bool)
```

#### func (*Struct) GetFloat

```go
func (r *Struct) GetFloat(key string) (float64, bool)
```

#### func (*Struct) GetInt

```go
func (r *Struct) GetInt(key string) (int, bool)
```

#### func (*Struct) GetSlice

```go
func (r *Struct) GetSlice(key string) ([]interface{}, bool)
```

#### func (*Struct) GetString

```go
func (r *Struct) GetString(key string) (string, bool)
```

#### func (*Struct) JsonMarshal

```go
func (r *Struct) JsonMarshal() []byte
```

#### func (*Struct) ValueAssignTo

```go
func (c *Struct) ValueAssignTo(keyPath string, valuePointer interface{}, default_val interface{})
```
将keyPath的设置复制到&value, 可以为空

    pageSize := 0
    c.ValueAssignTo("the.key.path.to.here", &pageSize, 10)

#### func (*Struct) ValueMustAssignTo

```go
func (c *Struct) ValueMustAssignTo(keyPath string, valuePointer interface{})
```
将keyPath的设置复制到&value, 必须存在

    pageSize := 0
    c.MustValueAssignTo("the.key.path.to.here", &pageSize)

#### type TLSConfig

```go
type TLSConfig struct {
	CrtPath string `json:"crt" yaml:"crt"`
	KeyPath string `json:"key" yaml:"key"`
}
```
