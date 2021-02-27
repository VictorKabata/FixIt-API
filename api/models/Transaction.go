package models

import (
	"errors"
	"html"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
)

type Transaction struct {
	ID        uint64    `gorm:"primary_key;auto_increment" json:"id"`
	UserID    uint32    `gorm:"not null" json:"user_id"`
	WorkerID  uint32    `gorm:"not null" json:"worker_id"`
	PostID    uint32    `gorm:"not null" json:"post_id"`
	WorkID    uint32    `gorm:"not null" json:"work_id"`
	Amount    string    `gorm:"not null" json:"amount"`
	Type      string    `gorm:"not null" json:"type"` //Mpesa or Cash
	User      User      `json:"user"`
	Worker    User      `json:"worker"`
	Post      Post      `json:"post"`
	Work      Work      `json:"work"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
}

func (t *Transaction) Prepare() {
	t.ID = 0
	t.Amount = html.EscapeString(strings.TrimSpace(t.Amount))
	t.Type = html.EscapeString(strings.TrimSpace(t.Type))
	t.User = User{}
	t.Worker = User{}
	t.Post = Post{}
	t.Work = Work{}
	t.CreatedAt = time.Now()
}

func (t *Transaction) Validate() error {
	if t.Type == "" {
		return errors.New("Required Type")
	}
	if t.Amount == "" {
		return errors.New("Required Amount")
	}
	if t.UserID < 1 {
		return errors.New("Required User ID")
	}
	if t.WorkerID < 1 {
		return errors.New("Required Worker ID")
	}
	if t.PostID < 1 {
		return errors.New("Required Post ID")
	}
	if t.WorkID < 1 {
		return errors.New("Required Work ID")
	}
	return nil
}

//Create a new transaction
func (t *Transaction) UploadTransaction(db *gorm.DB) (*Transaction, error) {
	var err error

	err = db.Debug().Model(&Transaction{}).Create(&t).Error
	if err != nil {
		return &Transaction{}, err
	}

	if t.ID != 0 {
		err = db.Debug().Model(&User{}).Where("id = ?", t.UserID).Take(&t.User).Error
		if err != nil {
			return &Transaction{}, err
		}

		err = db.Debug().Model(&User{}).Where("id = ?", t.WorkerID).Take(&t.Worker).Error
		if err != nil {
			return &Transaction{}, err
		}

		err = db.Debug().Model(&Post{}).Where("id = ?", t.PostID).Take(&t.Post).Error
		if err != nil {
			return &Transaction{}, err
		}

		err = db.Debug().Model(&Work{}).Where("id = ?", t.WorkID).Take(&t.Work).Error
		if err != nil {
			return &Transaction{}, err
		}
	}
	return t, nil
}

//Return all transactions
func (t *Transaction) FindAllTransactions(db *gorm.DB) (*[]Transaction, error) {
	var err error

	transactions := []Transaction{}

	err = db.Debug().Model(&Transaction{}).Limit(100).Find(&transactions).Error
	if err != nil {
		return &[]Transaction{}, err
	}

	if len(transactions) > 0 {
		for i, _ := range transactions {
			err := db.Debug().Model(&User{}).Where("id = ?", transactions[i].UserID).Take(&transactions[i].User).Error
			if err != nil {
				return &[]Transaction{}, err
			}
		}
		for i, _ := range transactions {
			err := db.Debug().Model(&User{}).Where("id = ?", transactions[i].WorkerID).Take(&transactions[i].Worker).Error
			if err != nil {
				return &[]Transaction{}, err
			}
		}
		for i, _ := range transactions {
			err := db.Debug().Model(&Post{}).Where("id = ?", transactions[i].PostID).Take(&transactions[i].Post).Error
			if err != nil {
				return &[]Transaction{}, err
			}
		}
		for i, _ := range transactions {
			err := db.Debug().Model(&Work{}).Where("id = ?", transactions[i].WorkID).Take(&transactions[i].Work).Error
			if err != nil {
				return &[]Transaction{}, err
			}
		}
	}
	return &transactions, nil
}

//Return transaction based on id
func (t *Transaction) FindTransactionByID(db *gorm.DB, pid uint64) (*Transaction, error) {
	var err error

	err = db.Debug().Model(&Transaction{}).Where("id = ?", pid).Take(&t).Error
	if err != nil {
		return &Transaction{}, err
	}

	if t.ID != 0 {
		err = db.Debug().Model(&User{}).Where("id = ?", t.UserID).Take(&t.User).Error
		if err != nil {
			return &Transaction{}, err
		}

		err = db.Debug().Model(&User{}).Where("id = ?", t.WorkerID).Take(&t.Worker).Error
		if err != nil {
			return &Transaction{}, err
		}

		err = db.Debug().Model(&Post{}).Where("id = ?", t.PostID).Take(&t.Post).Error
		if err != nil {
			return &Transaction{}, err
		}

		err = db.Debug().Model(&Work{}).Where("id = ?", t.WorkID).Take(&t.Work).Error
		if err != nil {
			return &Transaction{}, err
		}
	}

	return t, nil
}

//Return transaction based on user id
func (t *Transaction) FindTransactionByUserID(db *gorm.DB, pid uint64) (*[]Transaction, error) {
	var err error

	transactions := []Transaction{}

	err = db.Debug().Model(&Transaction{}).Where("user_id = ?", pid).Limit(100).Find(&transactions).Error
	if err != nil {
		return &[]Transaction{}, err
	}

	if len(transactions) > 0 {
		for i, _ := range transactions {
			err := db.Debug().Model(&User{}).Where("id = ?", transactions[i].UserID).Take(&transactions[i].User).Error
			if err != nil {
				return &[]Transaction{}, err
			}
		}
		for i, _ := range transactions {
			err := db.Debug().Model(&User{}).Where("id = ?", transactions[i].WorkerID).Take(&transactions[i].Worker).Error
			if err != nil {
				return &[]Transaction{}, err
			}
		}
		for i, _ := range transactions {
			err := db.Debug().Model(&Post{}).Where("id = ?", transactions[i].PostID).Take(&transactions[i].Post).Error
			if err != nil {
				return &[]Transaction{}, err
			}
		}
		for i, _ := range transactions {
			err := db.Debug().Model(&Work{}).Where("id = ?", transactions[i].WorkID).Take(&transactions[i].Work).Error
			if err != nil {
				return &[]Transaction{}, err
			}
		}
	}
	return &transactions, nil
}
