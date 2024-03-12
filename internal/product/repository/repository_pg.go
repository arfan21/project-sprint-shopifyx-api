package productrepo

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/arfan21/project-sprint-shopifyx-api/internal/entity"
	"github.com/arfan21/project-sprint-shopifyx-api/internal/model"
	dbpostgres "github.com/arfan21/project-sprint-shopifyx-api/pkg/db/postgres"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type Repository struct {
	db dbpostgres.Queryer
}

func New(db dbpostgres.Queryer) *Repository {
	return &Repository{db: db}
}

func (r Repository) Begin(ctx context.Context) (tx pgx.Tx, err error) {
	return r.db.Begin(ctx)
}

func (r Repository) WithTx(tx pgx.Tx) *Repository {
	r.db = tx
	return &r
}

func (r Repository) Create(ctx context.Context, data entity.Product) (err error) {
	query := `
		INSERT INTO products (id, name, price, imageUrl, stock, condition, tags, isPurchaseable, user_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8,  $9)
	`

	_, err = r.db.Exec(ctx, query,
		data.ID,
		data.Name,
		data.Price,
		data.ImageUrl,
		data.Stock,
		data.Condition,
		data.Tags,
		data.IsPurchaseable,
		data.UserID,
	)
	if err != nil {
		err = fmt.Errorf("product.repository.Create: failed to create product: %w", err)
		return
	}

	return
}

func (r Repository) Update(ctx context.Context, data entity.Product) (err error) {
	query := `
		UPDATE products
		SET name = $1, price = $2, imageUrl = $3, condition = $4, tags = $5, isPurchaseable = $6
		WHERE id = $7
	`

	_, err = r.db.Exec(ctx, query,
		data.Name,
		data.Price,
		data.ImageUrl,
		data.Condition,
		data.Tags,
		data.IsPurchaseable,
		data.ID,
	)
	if err != nil {
		err = fmt.Errorf("product.repository.Update: failed to update product: %w", err)
		return
	}

	return
}

func (r Repository) GetByID(ctx context.Context, id uuid.UUID) (product entity.Product, err error) {
	query := `
		SELECT id, name, price, imageUrl, stock, condition, tags, isPurchaseable, user_id
		FROM products
		WHERE id = $1
	`

	err = r.db.QueryRow(ctx, query, id).Scan(
		&product.ID,
		&product.Name,
		&product.Price,
		&product.ImageUrl,
		&product.Stock,
		&product.Condition,
		&product.Tags,
		&product.IsPurchaseable,
		&product.UserID,
	)
	if err != nil {
		err = fmt.Errorf("product.repository.GetByID: failed to get product by id: %w", err)
		return
	}

	return
}

func (r Repository) Delete(ctx context.Context, id uuid.UUID) (err error) {
	query := `
		DELETE FROM products
		WHERE id = $1
	`

	_, err = r.db.Exec(ctx, query, id)
	if err != nil {
		err = fmt.Errorf("product.repository.Delete: failed to delete product: %w", err)
		return
	}

	return
}

func (r Repository) queryGetListWithFilter(ctx context.Context, query string, filter model.ProductGetListRequest) (rows pgx.Rows, err error) {
	arrArgs := []interface{}{}
	andStatement := " AND "
	whereQuery := ""

	if filter.UserOnly {
		arrArgs = append(arrArgs, filter.UserID)
		lenArgsStr := strconv.Itoa(len(arrArgs))

		whereQuery += fmt.Sprintf("user_id = $%s %s", lenArgsStr, andStatement)
	}

	if filter.Condition != "" {
		arrArgs = append(arrArgs, filter.Condition)
		lenArgsStr := strconv.Itoa(len(arrArgs))

		whereQuery += fmt.Sprintf("condition = $%s %s", lenArgsStr, andStatement)
	}

	if len(filter.Tags) > 0 {
		arrArgs = append(arrArgs, filter.Tags)
		lenArgsStr := strconv.Itoa(len(arrArgs))

		whereQuery += fmt.Sprintf("tags = ANY($%s) %s", lenArgsStr, andStatement)
	}

	if !filter.ShowEmptyStock {
		arrArgs = append(arrArgs, 0)
		whereQuery += fmt.Sprintf("stock > $%d %s", len(arrArgs), andStatement)
	}

	if filter.MaxPrice > 0 {
		arrArgs = append(arrArgs, filter.MinPrice)
		lenArgsStr := strconv.Itoa(len(arrArgs))

		whereQuery += fmt.Sprintf("price >= $%s %s", lenArgsStr, andStatement)

		arrArgs = append(arrArgs, filter.MaxPrice)
		lenArgsStr = strconv.Itoa(len(arrArgs))

		whereQuery += fmt.Sprintf("price <= $%s %s", lenArgsStr, andStatement)
	}

	if filter.Search != "" {
		arrArgs = append(arrArgs, strings.ToLower(filter.Search))
		lenArgsStr := strconv.Itoa(len(arrArgs))

		whereQuery += fmt.Sprintf("LOWER(name) LIKE '%%$%s%%' %s", lenArgsStr, andStatement)

	}

	// if lenArgs  > 0, add WHERE statement and remove last AND
	if lenArgs := len(arrArgs); lenArgs > 0 {
		whereQuery = "WHERE " + whereQuery[:len(whereQuery)-len(andStatement)] + " "
	}

	query += whereQuery

	if !filter.DisableOrder {
		sortBy := "id"
		if filter.SortBy != "" && filter.SortBy != "date" {
			sortBy = "price"
		}

		query += fmt.Sprintf("ORDER BY %s ", sortBy)

		orderBy := "DESC"
		if filter.OrderBy != "" && filter.OrderBy != "desc" {
			orderBy = "ASC"
		}
		query += fmt.Sprintf("%s ", orderBy)
	}

	if !filter.DisablePaging && filter.Limit > 0 {
		arrArgs = append(arrArgs, filter.Limit)
		query += fmt.Sprintf("LIMIT $%d ", len(arrArgs))

		arrArgs = append(arrArgs, filter.Offset)
		query += fmt.Sprintf("OFFSET $%d ", len(arrArgs))
	}
	return r.db.Query(ctx, query, arrArgs...)
}

func (r Repository) GetList(ctx context.Context, filter model.ProductGetListRequest) (res []entity.Product, err error) {
	query := `
		SELECT id, name, price, imageUrl, stock, condition, tags, isPurchaseable
		FROM products
	`

	rows, err := r.queryGetListWithFilter(ctx, query, filter)
	if err != nil {
		err = fmt.Errorf("product.repository.GetList: failed to get list product: %w", err)
		return
	}

	defer rows.Close()

	for rows.Next() {
		var product entity.Product
		if err = rows.Scan(
			&product.ID,
			&product.Name,
			&product.Price,
			&product.ImageUrl,
			&product.Stock,
			&product.Condition,
			&product.Tags,
			&product.IsPurchaseable,
		); err != nil {
			err = fmt.Errorf("product.repository.GetList: failed to scan product: %w", err)
			return
		}
		res = append(res, product)
	}

	return
}

func (r Repository) GetTotal(ctx context.Context, filter model.ProductGetListRequest) (total int, err error) {
	query := `
		SELECT COUNT(id)
		FROM products
	`

	filter.DisableOrder = true
	filter.DisablePaging = true

	rows, err := r.queryGetListWithFilter(ctx, query, filter)
	if err != nil {
		err = fmt.Errorf("product.repository.GetTotal: failed to get total product: %w", err)
		return
	}

	defer rows.Close()

	rows.Next()
	err = rows.Scan(&total)
	if err != nil {
		err = fmt.Errorf("product.repository.GetTotal: failed to scan total product: %w", err)
		return
	}

	return
}

func (r Repository) GetDetailByID(ctx context.Context, id uuid.UUID) (product entity.Product, err error) {
	query := `
		SELECT id, name, price, imageUrl, stock, condition, tags, isPurchaseable, user_id
		FROM products
		WHERE id = $1
	`

	err = r.db.QueryRow(ctx, query, id).Scan(
		&product.ID,
		&product.Name,
		&product.Price,
		&product.ImageUrl,
		&product.Stock,
		&product.Condition,
		&product.Tags,
		&product.IsPurchaseable,
		&product.UserID,
	)
	if err != nil {
		err = fmt.Errorf("product.repository.GetDetailByID: failed to get product by id: %w", err)
		return
	}

	return
}

func (r Repository) GetStockByIDForUpdate(ctx context.Context, id uuid.UUID) (stock int, err error) {
	query := `
		SELECT stock
		FROM products
		WHERE id = $1
		FOR UPDATE
	`

	err = r.db.QueryRow(ctx, query, id).Scan(&stock)
	if err != nil {
		err = fmt.Errorf("product.repository.GetStockByIDForUpdate: failed to get stock by id for update: %w", err)
		return
	}

	return
}

func (r Repository) UpdateStock(ctx context.Context, id uuid.UUID, stock int) (err error) {
	query := `
		UPDATE products
		SET stock = $1
		WHERE id = $2
	`

	_, err = r.db.Exec(ctx, query, stock, id)
	if err != nil {
		err = fmt.Errorf("product.repository.UpdateStock: failed to update stock: %w", err)
		return
	}

	return
}
