package utils

type PaginationResponse[T any] struct {
	Data     []T                `json:"data"`
	Metadata PaginationMetadata `json:"metadata"`
}

type PaginationMetadata struct {
	Page      int64 `json:"page"`
	PerPage   int64 `json:"perPage"`
	TotalData int64 `json:"totalData"`
	TotalPage int64 `json:"totalPage"`
}
