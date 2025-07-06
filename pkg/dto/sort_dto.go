package dto

const (
	SortByID        = "id"
	SortByName      = "name"
	SortByCreatedAt = "created_at"
)

// todo
// ?sort_fields=id desc,name asc
type SortParam struct {
	SortField string `json:"sort_field"`
}
