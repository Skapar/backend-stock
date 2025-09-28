package repository

import (
	"context"
	"time"

	"github.com/onec-tech/bot/internal/models/entities"
	"github.com/onec-tech/bot/pkg/database"
	"github.com/onec-tech/bot/pkg/logger"
)

type pgRepository struct {
	DB  database.IDatabase
	log logger.Logger
}

func NewPGRepository(db database.IDatabase, log logger.Logger) PGRepository {
	return &pgRepository{
		DB:  db,
		log: log,
	}
}

func (r *pgRepository) CreateOrUpdateUser(ctx context.Context, user *entities.User) error {
	q := `
	INSERT INTO onec_user (tg_id, nickname, name, phone)
	VALUES ($1, $2, $3, $4)
	ON CONFLICT (tg_id) 
	DO UPDATE SET
		nickname = EXCLUDED.nickname,
		name = EXCLUDED.name,
		phone = EXCLUDED.phone
	RETURNING id;
	`

	return r.DB.Insert(ctx, user, q, user.TGID, user.Nickname, user.Name, user.Phone)
}

func (r *pgRepository) GetUserByTGID(ctx context.Context, tgID int64) (*entities.User, error) {
	q := `
	SELECT id, tg_id, nickname, name, phone
	FROM onec_user
	WHERE tg_id = $1;
	`
	var user entities.User
	err := r.DB.GetOne(ctx, &user, q, tgID)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *pgRepository) CreateReceipt(ctx context.Context, receipt *entities.Receipt) error {
	q := `
	INSERT INTO onec_receipt (user_id, file_path, status)
	VALUES ($1, $2, $3)
	RETURNING id;
	`
	return r.DB.Insert(ctx, receipt, q, receipt.UserID, receipt.FilePath, receipt.Status)
}

func (r *pgRepository) GetReceiptsByStatus(ctx context.Context, status entities.StatusType) ([]entities.ReceiptWithUser, error) {
	var receipts []entities.ReceiptWithUser
	q := `
		SELECT r.id, r.user_id, r.file_path, r.status, u.tg_id
		FROM onec_receipt r
		JOIN onec_user u ON u.id = r.user_id
		WHERE r.status = $1
	`
	err := r.DB.Get(ctx, &receipts, q, status)
	return receipts, err
}

func (r *pgRepository) UpdateReceiptStatus(ctx context.Context, status entities.StatusType, receiptID int64) error {
	q := `UPDATE onec_receipt SET status = $1 WHERE id = $2`
	return r.DB.Update(ctx, nil, q, status, receiptID)
}

func (r *pgRepository) GetDefaultSubscription(ctx context.Context) (*entities.Subscription, error) {
	var sub entities.Subscription
	q := `SELECT id, price, duration_days, total_cups FROM onec_subscription WHERE id = 1`
	err := r.DB.GetOne(ctx, &sub, q)
	return &sub, err
}

func (r *pgRepository) CreateUserSubscription(ctx context.Context, userID int64, sub *entities.Subscription) (int64, error) {
	start := time.Now()
	end := start.AddDate(0, 0, int(sub.DurationDays))

	var id int64
	q := `
		INSERT INTO onec_user_subscription 
		    (user_id, subscription_id, start_date, end_date, total_cups, remaining_cups) 
		VALUES ($1, $2, $3, $4, $5, $5)
		RETURNING id
	`
	err := r.DB.Insert(ctx, &id, q, userID, sub.ID, start, end, sub.TotalCups)
	return id, err
}

func (r *pgRepository) CreatePayment(ctx context.Context, userID, userSubID int64, amount int64) error {
	q := `
		INSERT INTO onec_payment (user_id, user_subscription_id, amount, status, payment_method)
		VALUES ($1, $2, $3, 'confirmed', 'manual')
	`
	return r.DB.Insert(ctx, nil, q, userID, userSubID, amount)
}
