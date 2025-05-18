package email

import (
	"context"
	_ "embed"
	"errors"
	"github.com/jackc/pgx/v5/pgxpool"
	"net"
)

//go:embed sql/insert_sent_email.sql
var insertIntoEmailSql string

//go:embed sql/get_sent_emails.sql
var getFromEmailSql string

type emailRepository struct {
	pool *pgxpool.Pool
}

type Repository interface {
	AddSentEmail(ctx context.Context, email string, ip net.IP, isSignUp bool) error
	GetSentEmails(ctx context.Context) ([]*SentEmail, error)
}

func (e *emailRepository) AddSentEmail(ctx context.Context, email string, ip net.IP, isSignUp bool) error {
	conn, err := e.pool.Exec(ctx, insertIntoEmailSql, email, ip, isSignUp)
	if err != nil {
		return err
	}

	if conn.RowsAffected() == 0 {
		return errors.New("no rows affected")
	}

	return nil
}
func (e *emailRepository) GetSentEmails(ctx context.Context) ([]*SentEmail, error) {
	rows, err := e.pool.Query(ctx, getFromEmailSql)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	var sentEmails []*SentEmail
	for rows.Next() {
		var sentEmail SentEmail
		err := rows.Scan(
			&sentEmail.Id,
			&sentEmail.CreatedAt,
			&sentEmail.Email,
			&sentEmail.Ip,
		)

		if err != nil {
			return nil, err
		}

		sentEmails = append(sentEmails, &sentEmail)
	}

	return sentEmails, nil
}

func NewRepo(pool *pgxpool.Pool) Repository {
	return &emailRepository{
		pool: pool,
	}
}
