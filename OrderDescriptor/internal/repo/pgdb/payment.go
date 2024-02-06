package pgdb

import (
	"context"
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/nktinn/OrderDescriptor/OrderDescriptor/internal/model"
	"github.com/nktinn/OrderDescriptor/OrderDescriptor/pkg/postgres"
)

type PaymentRepo struct {
	pg *postgres.Postgres
}

func NewPaymentRepo(pg *postgres.Postgres) *PaymentRepo {
	return &PaymentRepo{pg}
}

func (r *PaymentRepo) CreatePayment(payment model.Payment) error {

	sql := `
	INSERT INTO payments (transaction_id, order_id, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`

	_, err := r.pg.Pool.Exec(context.Background(), sql,
		payment.TransactionID,
		payment.OrderID,
		payment.RequestID,
		payment.Currency,
		payment.Provider,
		payment.Amount,
		payment.PaymentDT,
		payment.Bank,
		payment.DeliveryCost,
		payment.GoodsTotal,
		payment.CustomFee,
	)
	if err != nil {
		return fmt.Errorf("failed to insert new payment: %v", err)
	}
	log.Infoln("New payment inserted")
	return nil
}

func (r *PaymentRepo) GetPayment(id string) (*model.Payment, error) {
	payment := &model.Payment{}

	query := `SELECT * FROM payments WHERE order_id = $1`
	err := r.pg.Pool.QueryRow(context.Background(), query, id).
		Scan(
			&payment.TransactionID,
			&payment.OrderID,
			&payment.RequestID,
			&payment.Currency,
			&payment.Provider,
			&payment.Amount,
			&payment.PaymentDT,
			&payment.Bank,
			&payment.DeliveryCost,
			&payment.GoodsTotal,
			&payment.CustomFee,
		)
	if err != nil {
		return nil, fmt.Errorf("failed to get payment: %v", err)
	}

	return payment, nil
}

func (r *PaymentRepo) GetAllPayments() []model.Payment {
	payments := make([]model.Payment, 0)

	query := `SELECT * FROM payments`
	rows, err := r.pg.Pool.Query(context.Background(), query)
	if err != nil {
		log.Errorf("Failed to get payments: %v", err)
		return nil
	}

	for rows.Next() {
		payment := model.Payment{}
		err := rows.Scan(
			&payment.TransactionID,
			&payment.OrderID,
			&payment.RequestID,
			&payment.Currency,
			&payment.Provider,
			&payment.Amount,
			&payment.PaymentDT,
			&payment.Bank,
			&payment.DeliveryCost,
			&payment.GoodsTotal,
			&payment.CustomFee,
		)
		if err != nil {
			log.Errorf("failed to scan payments: %v", err)
			return nil
		}
		payments = append(payments, payment)
	}

	log.Infoln("Payments wrote to memory")
	return payments
}

func (r *PaymentRepo) DeletePayment(uid string) error {
	sql := `DELETE FROM payments WHERE order_id = $1`

	_, err := r.pg.Pool.Exec(context.Background(), sql, uid)
	if err != nil {
		log.Errorf("Failed to delete payment from db: %v", err)
		return err
	}
	log.Infoln("Payment deleted from db")
	return nil
}

func (r *PaymentRepo) DeleteAllPayments() error {
	sql := `DELETE FROM payments`

	_, err := r.pg.Pool.Exec(context.Background(), sql)
	if err != nil {
		log.Errorf("Failed to delete all payments from db: %v", err)
		return err
	}
	log.Infoln("All payments deleted from db")
	return nil
}
