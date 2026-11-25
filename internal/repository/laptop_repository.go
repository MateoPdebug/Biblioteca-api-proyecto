package repository

import (
	"biblioteca-api/internal/models"
	"github.com/jmoiron/sqlx"
	"time"
)

type LaptopRepository struct{ db *sqlx.DB }

func NewLaptopRepository(db *sqlx.DB) *LaptopRepository { return &LaptopRepository{db} }
func (r *LaptopRepository) Create(req models.LaptopRequest) (*models.Laptop, error) {
	res, err := r.db.Exec(`INSERT INTO laptops(serial,brand,model) VALUES (?,?,?)`, req.Serial, req.Brand, req.Model)
	if err != nil {
		return nil, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	var l models.Laptop
	err = r.db.Get(&l, `SELECT * FROM laptops WHERE id=?`, id)
	return &l, err
}
func (r *LaptopRepository) List() ([]models.Laptop, error) {
	l := []models.Laptop{}
	err := r.db.Select(&l, `SELECT * FROM laptops ORDER BY id`)
	return l, err
}
func (r *LaptopRepository) Begin() (*sqlx.Tx, error) { return r.db.Beginx() }
func (r *LaptopRepository) Lock(tx *sqlx.Tx, id int) (models.Laptop, error) {
	var l models.Laptop
	err := tx.Get(&l, `SELECT * FROM laptops WHERE id=? FOR UPDATE`, id)
	return l, err
}
func (r *LaptopRepository) SetAvailable(tx *sqlx.Tx, id int, available bool) error {
	_, err := tx.Exec(`UPDATE laptops SET available=? WHERE id=?`, available, id)
	return err
}

type LaptopLoanRepository struct{ db *sqlx.DB }

func NewLaptopLoanRepository(db *sqlx.DB) *LaptopLoanRepository { return &LaptopLoanRepository{db} }

const loanColumns = `id,user_id,laptop_id,loaned_at,due_at,returned_at,status`

func (r *LaptopLoanRepository) Create(tx *sqlx.Tx, userID, laptopID int, due time.Time) (*models.LaptopLoan, error) {
	res, err := tx.Exec(`INSERT INTO laptop_loans(user_id,laptop_id,due_at) VALUES (?,?,?)`, userID, laptopID, due)
	if err != nil {
		return nil, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	var l models.LaptopLoan
	err = tx.Get(&l, `SELECT `+loanColumns+` FROM laptop_loans WHERE id=?`, id)
	return &l, err
}
func (r *LaptopLoanRepository) Lock(tx *sqlx.Tx, id int) (models.LaptopLoan, error) {
	var l models.LaptopLoan
	err := tx.Get(&l, `SELECT `+loanColumns+` FROM laptop_loans WHERE id=? FOR UPDATE`, id)
	return l, err
}
func (r *LaptopLoanRepository) MarkReturned(tx *sqlx.Tx, id int) error {
	_, err := tx.Exec(`UPDATE laptop_loans SET status='returned',returned_at=NOW() WHERE id=?`, id)
	return err
}
func (r *LaptopLoanRepository) ListByUser(id int) ([]models.LaptopLoan, error) {
	l := []models.LaptopLoan{}
	err := r.db.Select(&l, `SELECT `+loanColumns+` FROM laptop_loans WHERE user_id=? ORDER BY id DESC`, id)
	return l, err
}
