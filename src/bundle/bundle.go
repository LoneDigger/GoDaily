package bundle

import "time"

var (
	CodeOk           = "E-000"
	CodeFormat       = "E-001" // 格式錯誤
	CodeDb           = "E-002" // 資料庫
	CodeLogin        = "E-003" // 尚未登入
	CodeUsername     = "E-011" // 使用者名稱錯誤
	CodePassword     = "E-012" // 密碼錯誤
	CodeToken        = "E-013" // 授權碼錯誤
	CodeUserRepeat   = "E-014" // 帳號重複
	CodeHold         = "E-015" // 不是你的
	CodeDate         = "E-016" // 查詢日期錯誤
	CodeCache        = "E-017" // 快取沒資料
	CodeTypeRepeat   = "E-018" // 類別名稱重複
	CodeNoData       = "E-019" // 沒有資料
	CodeEmptyContent = "E-020" // 沒有輸入關鍵字
)

// 全部類型
type AllType struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	Subs []Sub  `json:"subs"`
}

// 主類別和
type MainSumMonthly struct {
	Name string `json:"name" db:"name"`
	Sum  int    `json:"sum" db:"sum"`
}

type Main struct {
	Id   int    `json:"id" swaggertype:"integer"`
	Name string `json:"name" swaggertype:"string"`
}

type Sub struct {
	Id       int    `json:"id" swaggertype:"integer"`
	Name     string `json:"name" swaggertype:"string"`
	Increase bool   `json:"increase" swaggertype:"boolean"`
}

type Item struct {
	Id     int       `json:"id" db:"id"`
	Name   string    `json:"name" db:"name"`
	MainId int       `json:"main_id" db:"main_id"`
	SubId  int       `json:"sub_id" db:"sub_id"`
	Price  int       `json:"price" db:"price"`
	Remark string    `json:"remark" db:"remark"`
	Date   time.Time `json:"date" db:"date"`
}

// 預覽項目
type PreviewItem struct {
	Id       int    `json:"id" db:"id"`
	MainId   int    `json:"main_id" db:"main_id"`
	MainName string `json:"main_name" db:"main_name"`
	SubId    int    `json:"sub_id" db:"sub_id"`
	SubName  string `json:"sub_name" db:"sub_name"`
	Name     string `json:"name" db:"name"`
	Increase int    `json:"increase" db:"increase"`
	Price    int    `json:"price" db:"price"`
	Date     string `json:"date" db:"date"`
}

// 月結花費
type Monthly struct {
	Sum  int       `json:"sum" db:"sum"`
	Date time.Time `json:"date" db:"date"`
}

// 建立使用者
// swagger:model CreateUserRequest
type CreateUserRequest struct {
	// 使用者帳號
	Username string `json:"username" binding:"required" validate:"required,min=3,max=32" swaggertype:"string" example:"username"`
	// 使用者密碼
	Password string `json:"password" binding:"required" validate:"required,min=3,max=64" swaggertype:"string" example:"password"`
	// 通行證
	Token string `json:"token" binding:"required" swaggertype:"string" example:"token"`
}

// 建立主類別請求
type CreateMainTypeRequest struct {
	Name string `json:"name" validate:"required,min=2,max=32" binding:"required" swaggertype:"string"`
}

// 建立主類別回應
type CreateMainTypeResponse struct {
	ErrorResponse
	MainId int    `json:"main_id"`
	Name   string `json:"name"`
}

// 建立子類別請求
// swagger:model CreateSubTypeRequest
type CreateSubTypeRequest struct {
	MainId   int    `json:"main_id" binding:"required" validate:"required,gt=0" swaggertype:"integer" example:"0"`
	Name     string `json:"name" binding:"required" validate:"required,min=2,max=32" swaggertype:"string" example:"name"`
	Increase bool   `json:"increase" swaggertype:"boolean" example:"false"`
}

// 建立子類別回應
type CreateSubTypeResponse struct {
	ErrorResponse
	SubId int    `json:"sub_id"`
	Name  string `json:"name"`
}

// 建立項目請求
// swagger:model CreateItemRequest
type CreateItemRequest struct {
	SubId  int    `json:"sub_id" binding:"required" validate:"required,gt=0" swaggertype:"integer" example:"0"`
	Name   string `json:"name" validate:"required,min=2,max=32" swaggertype:"string" example:"name"`
	Price  int    `json:"price" binding:"required" validate:"required,gt=0" swaggertype:"integer" example:"100"`
	Remark string `json:"remark" validate:"required,min=0,max=64" swaggertype:"string" example:""`
	Date   string `json:"date" binding:"required" time_format:"2006-01-02" example:"2006-01-02"`
}

// 建立項目回應
type CreateItemResponse struct {
	ErrorResponse
}

// 回應
type ErrorResponse struct {
	Code string `json:"code"` // 錯誤代號
}

// swagger:model LoginRequest
type LoginRequest struct {
	// 使用者帳號
	Username string `json:"username"  binding:"required" validate:"required,min=3,max=32" swaggertype:"string" example:"username"`
	// 使用者密碼
	Password string `json:"password" binding:"required" validate:"required,min=3,max=64" swaggertype:"string" example:"password"`
	// 通行證
	Token string `json:"token" binding:"required" swaggertype:"string" example:"token"`
}

// 取得全部類別
type GetAllTypeResponse struct {
	ErrorResponse
	List []AllType `json:"list"`
}

// 取得單個
type GetItemResponse struct {
	ErrorResponse
	Item Item `json:"item"`
}

// 取得多個
type GetItemsResponse struct {
	ErrorResponse
	List []PreviewItem `json:"list"`
}

// 取得主類別清單
type GetMainTypeResponse struct {
	ErrorResponse
	List []Main `json:"list"`
}

type GetSpendByMonthlyResponse struct {
	ErrorResponse
	List []Monthly `json:"list"`
}

// 取得子類別清單
type GetSubTypeResponse struct {
	ErrorResponse
	List []Sub `json:"list"`
}

type GetSumByMainTypeResponse struct {
	ErrorResponse
	List []MainSumMonthly `json:"list"`
}

// 更新項目請求
// swagger:model UpdateItemRequest
type UpdateItemRequest struct {
	ItemId int    `json:"item_id" binding:"required" validate:"required,gt=0" swaggertype:"integer"`
	SubId  int    `json:"sub_id" binding:"required" validate:"required,gt=0" swaggertype:"integer" example:"0"`
	Name   string `json:"name" validate:"required,min=2,max=32" swaggertype:"string" example:"name"`
	Price  int    `json:"price" binding:"required" validate:"required,gt=0" swaggertype:"integer" example:"100"`
	Remark string `json:"remark" validate:"required,min=0,max=64" swaggertype:"string" example:""`
	Date   string `json:"date" binding:"required" time_format:"2006-01-02" example:"2006-01-02"`
}

// 更新項目回應
type UpdateItemResponse struct {
	ErrorResponse
}

// 更新主類別請求
// swagger:model UpdateMainTypeRequest
type UpdateMainTypeRequest struct {
	MainId int    `json:"main_id" binding:"required" validate:"required,gt=0" swaggertype:"integer"`
	Name   string `json:"name" binding:"required" validate:"required,min=2,max=32" swaggertype:"string"`
}

// 更新主類別回應
type UpdateMainTypeResponse struct {
	ErrorResponse
}

// 更新子類別請求
// swagger:model UpdateSubTypeRequest
type UpdateSubTypeRequest struct {
	SubId    int    `json:"sub_id" binding:"required" validate:"required,gt=0" swaggertype:"integer"`
	Name     string `json:"name" binding:"required" validate:"required,min=2,max=32" swaggertype:"string"`
	Increase bool   `json:"increase" swaggertype:"boolean" example:"false"`
}

// 更新子類別回應
type UpdateSubTypeResponse struct {
	ErrorResponse
}

// 刪除主類別回應
type DeleteMainTypeResponse struct {
	ErrorResponse
}

// 刪除子類別回應
type DeleteSubTypeResponse struct {
	ErrorResponse
}
