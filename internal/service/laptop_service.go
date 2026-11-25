package service

import (
	"biblioteca-api/internal/models"
	"biblioteca-api/internal/repository"
	"database/sql"
	"errors"
	"github.com/go-sql-driver/mysql"
	"strings"
	"time"
)

var (
	ErrLaptopUnavailable = errors.New("laptop no disponible")
	ErrLaptopNotFound    = errors.New("laptop o préstamo no encontrado")
	ErrLaptopForbidden   = errors.New("no puedes devolver un préstamo ajeno")
	ErrLaptopReturned    = errors.New("préstamo ya devuelto")
	ErrLaptopInvalid     = errors.New("datos de laptop inválidos")
	ErrLaptopSerial      = errors.New("serial ya registrado")
)

const laptopLoanDays = 7

type LaptopService struct {
	laptops *repository.LaptopRepository
	loans   *repository.LaptopLoanRepository
}

func NewLaptopService(l *repository.LaptopRepository, r *repository.LaptopLoanRepository) *LaptopService {
	return &LaptopService{l, r}
}
func (s *LaptopService) Create(req models.LaptopRequest) (*models.Laptop, error) {
	req.Serial = strings.TrimSpace(req.Serial)
	req.Brand = strings.TrimSpace(req.Brand)
	req.Model = strings.TrimSpace(req.Model)
	if req.Serial == "" || req.Brand == "" || req.Model == "" {
		return nil, ErrLaptopInvalid
	}
	l, err := s.laptops.Create(req)
	var me *mysql.MySQLError
	if errors.As(err, &me) && me.Number == 1062 {
		return nil, ErrLaptopSerial
	}
	return l, err
}
func (s *LaptopService) List() ([]models.Laptop, error) { return s.laptops.List() }
func (s *LaptopService) MyLoans(user int) ([]models.LaptopLoan, error) {
	return s.loans.ListByUser(user)
}
func (s *LaptopService) Borrow(user, id int) (*models.LaptopLoan, error) {
	tx, err := s.laptops.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	l, err := s.laptops.Lock(tx, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrLaptopNotFound
	}
	if err != nil {
		return nil, err
	}
	if !l.Available {
		return nil, ErrLaptopUnavailable
	}
	loan, err := s.loans.Create(tx, user, id, time.Now().AddDate(0, 0, laptopLoanDays))
	if err != nil {
		return nil, err
	}
	if err = s.laptops.SetAvailable(tx, id, false); err != nil {
		return nil, err
	}
	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return loan, nil
}
func (s *LaptopService) Return(user int, role string, id int) error {
	tx, err := s.laptops.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	l, err := s.loans.Lock(tx, id)
	if errors.Is(err, sql.ErrNoRows) {
		return ErrLaptopNotFound
	}
	if err != nil {
		return err
	}
	if l.UserID != user && role != "admin" {
		return ErrLaptopForbidden
	}
	if l.Status != "active" {
		return ErrLaptopReturned
	}
	if err = s.loans.MarkReturned(tx, id); err != nil {
		return err
	}
	if err = s.laptops.SetAvailable(tx, l.LaptopID, true); err != nil {
		return err
	}
	return tx.Commit()
}
