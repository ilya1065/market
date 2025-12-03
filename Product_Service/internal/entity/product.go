package entity

import "time"

type Product struct {
	Id          uint      `db:"id" json:"id"`
	Price       float64   `db:"price" json:"price"`
	Name        string    `db:"name" json:"name"`
	Description string    `db:"description" json:"description"`
	Stock       int       `db:"stock" json:"stock"`
	ImageURL    string    `db:"image_url" json:"image_url"`
	CategoryID  int       `db:"category_id" json:"categoryID"`
	CreatedAt   time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt   time.Time `db:"updated_at" json:"updatedAt"`
}
