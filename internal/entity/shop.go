package entity

import (
	"time"
)

type Shop struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	OpenTime    time.Time `json:"open_time"`
	CloseTime   time.Time `json:"close_time"`
	IsClosed    bool      `json:"is_closed"`
	CreatedAt   time.Time `json:"created_at"`
}

type CreateShop struct {
	Name        string     `json:"name"`
	Description string     `json:"description"`
	OpenTime    time.Time  `json:"open_time"`
	CloseTime   time.Time  `json:"close_time"`
	IsClosed    bool       `json:"is_closed"`
	Menuitems   []MenuItem `json:"menuitems"`
	UserId      int64      `json:"user_id"`
}

type MenuItem struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Photo       string    `json:"photo"`
	Price       int32     `json:"price"`
	Shop        int64     `json:"shop_id"`
	CreatedAt   time.Time `json:"created_at"`
}

type UpdateShopInfo struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	OpenTime    time.Time `json:"open_time"`
	CloseTime   time.Time `json:"close_time"`
	IsClosed    bool      `json:"is_closed"`
	UserId      int64     `json:"user_id"`
}

type CreateMenuItem struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Photo       string `json:"photo"`
	Price       int32  `json:"price"`
	ShopID      int64  `json:"shop_id"`
}

type UpdateMenuItem struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       int32  `json:"price"`
}
