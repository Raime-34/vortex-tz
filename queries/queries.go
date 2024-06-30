package queries

import (
	"database/sql"
	"fmt"
	"vortex-tz/models"
	"vortex-tz/utils"
)

func GetOrderBook(exchangeName string, pair string) (*models.OrderBook, error) {
	config, err := utils.LoadConfig("config.json")
	if err != nil {
		return nil, err
	}

	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.Host, config.Port, config.User, config.Password, config.DBName)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	id, err := getOrderBookID(db, exchangeName, pair)
	if err != nil {
		return nil, err
	}

	orderBook := models.OrderBook{
		Exchange: exchangeName,
		Pair:     pair,
		Asks:     make([]*models.DepthOrder, 0),
		Bids:     make([]*models.DepthOrder, 0),
	}

	asks, bids, err := getDepthOrders(db, id)
	if err != nil {
		return nil, err
	}

	orderBook.Asks = asks
	orderBook.Bids = bids

	return &orderBook, nil
}

func getOrderBookID(db *sql.DB, exchangeName string, pair string) (int64, error) {
	rows, err := db.Query(
		"SELECT id FROM public.\"OrderBook\" WHERE exchange = $1 and pair = $2 LIMIT 1",
		exchangeName, pair,
	)
	if err != nil {
		return -1, err
	}

	var id int64
	for rows.Next() {
		err = rows.Scan(&id)
		if err != nil {
			return -1, err
		}
	}
	return id, nil
}

func getDepthOrders(db *sql.DB, orderBookID int64) (asks, bids []*models.DepthOrder, err error) {
	rows, err := db.Query(
		"SELECT * FROM public.\"DepthOrder\" WHERE order_book_id = $1",
		orderBookID,
	)
	if err != nil {
		return
	}

	asks = make([]*models.DepthOrder, 0)
	bids = make([]*models.DepthOrder, 0)

	for rows.Next() {
		var (
			id            int64
			price         float64
			qty           float64
			t             string
			order_book_id int64
		)

		err := rows.Scan(&id, &price, &qty, &t, &order_book_id)
		if err != nil {
			continue
		}

		switch t {
		case "asks":
			asks = append(asks, &models.DepthOrder{
				Price:   price,
				BaseQty: qty,
			})
		case "bids":
			bids = append(bids, &models.DepthOrder{
				Price:   price,
				BaseQty: qty,
			})
		}
	}

	return
}

func SaveOrderBook(exchangeName string, pair string, asks []*models.DepthOrder, bids []*models.DepthOrder) error {
	config, err := utils.LoadConfig("config.json")
	if err != nil {
		return err
	}

	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.Host, config.Port, config.User, config.Password, config.DBName)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return err
	}

	orderBookId, err := insertOrderBook(db, exchangeName, pair)
	if err != nil {
		return err
	}

	insertArrayOfDepthOfOrder(db, "asks", orderBookId, asks)
	insertArrayOfDepthOfOrder(db, "bids", orderBookId, bids)

	return nil
}

func insertArrayOfDepthOfOrder(db *sql.DB, t string, orderBookId int64, depthOrders []*models.DepthOrder) {
	for i := 0; i < len(depthOrders); i++ {
		_, err := insertDepthOrder(db, depthOrders[i], t, orderBookId)
		if err != nil {
			println(err.Error())
		}
	}
}

func insertDepthOrder(db *sql.DB, depthOrders *models.DepthOrder, t string, orderBookId int64) (int64, error) {
	rows, err := db.Query(
		"insert into public.\"DepthOrder\" (price, qty, type, order_book_id) values ($1, $2, $3, $4) returning id",
		depthOrders.Price, depthOrders.BaseQty, t, orderBookId,
	)
	if err != nil {
		return -1, err
	}
	var id int64
	for rows.Next() {
		err = rows.Scan(&id)
		if err != nil {
			return -1, err
		}
	}
	return id, nil
}

func insertOrderBook(db *sql.DB, exchangeName string, pair string) (int64, error) {
	rows, err := db.Query(
		"insert into public.\"OrderBook\" (exchange, pair) values ($1, $2) returning id",
		exchangeName, pair,
	)

	if err != nil {
		return -1, err
	}

	var id int64
	for rows.Next() {
		err = rows.Scan(&id)
		if err != nil {
			return -1, err
		}
	}
	return id, nil
}
