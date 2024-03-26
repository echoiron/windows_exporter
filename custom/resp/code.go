package resp

// xxx 前三位采用谷歌http错误码  xx 模块 xx 错误类型
// 400 客户端错误
// 401 权限错误
// 500 服务器错误
// 501 数据库错误

var (
	OK           = Resp{Code: 200, Msg: "OK"}
	FAIL         = Resp{Code: 10001, Msg: "操作失败"}
	InvalidParam = Resp{Code: 10003, Msg: "参数错误"}
	Unauthorized = Resp{Code: 10004, Msg: "认证失败"}
	ForBidden    = Resp{Code: 10005, Msg: "没有权限"}
	NotFound     = Resp{Code: 10006, Msg: "未找到资源"}
	LoginFail    = Resp{Code: 10007, Msg: "登录失败"}

	// 数据库错误
	DbConnectError      = Resp{Code: 5010001, Msg: "数据库连接失败"}
	DbSelectError       = Resp{Code: 5010002, Msg: "数据查询失败"}
	DbUpdateError       = Resp{Code: 5010003, Msg: "数据更新失败"}
	DbInsertError       = Resp{Code: 5010004, Msg: "数据插入失败"}
	DbDeleteError       = Resp{Code: 5010005, Msg: "数据删除失败"}
	RedisConnectError   = Resp{Code: 5010006, Msg: "redis连接失败"}
	RedisExecError      = Resp{Code: 5010006, Msg: "redis执行失败"}
	DbSessionStartError = Resp{Code: 5010007, Msg: "数据库事务开启失败"}
	DbDataCheckError    = Resp{Code: 5010008, Msg: "数据去重验证失败"}
	EsConnectError      = Resp{Code: 5010009, Msg: "es连接失败"}

	AddSheetError      = Resp{Code: 5000001, Msg: "新建sheet失败"}
	UnableAccountLock  = Resp{Code: 5000002, Msg: "暂无账号锁定"}
	CmdExecError       = Resp{Code: 5000003, Msg: "执行shell命令失败"}
	FileExecError      = Resp{Code: 5000004, Msg: "文件执行失败"}
	YamlAnalysisError  = Resp{Code: 5000005, Msg: "yaml解析错误"}
	MarshalError       = Resp{Code: 5000006, Msg: "文件解析错误"}
	UploadFileError    = Resp{Code: 5000007, Msg: "文件上传失败"}
	GetFileStreamError = Resp{Code: 5000008, Msg: "获取文件流失败"}
	ReadFileError      = Resp{Code: 5000009, Msg: "读取文件失败"}
	FailedToParseFile  = Resp{Code: 5000010, Msg: "解析文件失败"}
)
