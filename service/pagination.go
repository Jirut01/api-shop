package service

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
)

const (
	DefaultQueryPage = 1
	DefaultQuerySize = 20
)

type Pagination struct {
	Page int
	Size int
}

func GetfiberPagination(ctx *fiber.Ctx) Pagination {
	p := Pagination{
		Page: atoi(ctx.Query("page"), DefaultQueryPage),
		Size: atoi(ctx.Query("size"), DefaultQuerySize),
	}
	return p
}

func (p *Pagination) Offset() int {
	return (p.Page - 1) * p.Size
}

func atoi(s string, v int) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		return v
	}
	return i
}
