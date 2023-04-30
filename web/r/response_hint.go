package r

const (
	FailParamParse     = "参数解析失败"
	FailParamValid     = "参数校验不通过"
	FailIDNotNil       = "ID值不能为空"
	FailRecordNotFound = "数据不存在"

	// CRUD

	OKCreate = "创建成功"
	OkRead   = "获取成功"
	OKUpdate = "更新成功"
	OKDelete = "删除成功"
	OkSave   = "保存成功"

	FailCreate = "创建失败"
	FailRead   = "获取失败"
	FailUpdate = "更新失败"
	FailDelete = "删除失败"
	FailSave   = "保存失败"

	// File

	OKUpload   = "上传成功"
	OkDownload = "下载成功"

	FailUpload   = "上传失败"
	FailDownload = "下载失败"

	FailFileRead   = "文件读取失败"
	FailFileParse  = "文件解析失败"
	FailFileFormat = "文件格式不支持"
)
