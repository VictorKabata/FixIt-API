package models

import (
	"errors"
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
)

type Booking struct {
	ID     uint32 `gorm:"primary_key;auto_increment" json:"id"`
	UserID uint32 `gorm:"not null" json:"user_id"`
	PostID uint32 `gorm:"not null" json:"post_id"`
	//Description string    `gorm:"not null" json:"description"`
	//Budget      string    `gorm:"not null" json:"budget"`
	Status    string    `gorm:"not null" json:"status"`
	User      User      `json:"user"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (b *Booking) Prepare() {
	b.ID = 0
	//b.Description = html.EscapeString(strings.TrimSpace(b.Description))
	//b.Budget = html.EscapeString(strings.TrimSpace(b.Budget))
	b.Status = "Pending"
	b.User = User{}
	b.CreatedAt = time.Now()
	b.UpdatedAt = time.Now()
}

func (b *Booking) Validate() error {

	// if b.Description == "" {
	// 	return errors.New("Required Description")
	// }
	// if b.Budget == "" {
	// 	return errors.New("Required Budget")
	// }
	if b.UserID < 1 {
		return errors.New("Required User ID")
	}
	if b.PostID < 1 {
		return errors.New("Required Post ID")
	}
	return nil
}

func (b *Booking) FindAllBookings(db *gorm.DB) (*[]Booking, error) {
	var err error

	booking := []Booking{}

	err = db.Debug().Model(&Booking{}).Limit(100).Find(&booking).Error
	if err != nil {
		return &[]Booking{}, err
	}

	if len(booking) > 0 {
		for i, _ := range booking {
			err := db.Debug().Model(&Post{}).Where("id = ?", booking[i].UserID).Take(&booking[i].User).Error
			if err != nil {
				return &[]Booking{}, err
			}
		}
	}

	// 	for i, _ := range booking {
	// 		err := db.Debug().Model(&User{}).Where("id = ?", booking[i].BookedPost.UserID).Take(&booking[i].BookedPost.User).Error
	// 		if err != nil {
	// 			return &[]Booking{}, err
	// 		}
	// 	}

	return &booking, nil
}

func (b *Booking) SaveBooking(db *gorm.DB) (*Booking, error) {
	var err error
	err = db.Debug().Model(&Booking{}).Create(&b).Error
	if err != nil {
		return &Booking{}, err
	}
	if b.ID != 0 {
		err = db.Debug().Model(&Post{}).Where("id = ?", b.UserID).Take(&b.User).Error
		if err != nil {
			fmt.Println(err)
			return &Booking{}, err
		}
	}
	return b, nil
}

func (b *Booking) UpdateBooking(db *gorm.DB) (*Booking, error) {

	var err error

	err = db.Debug().Model(&Booking{}).Where("id = ?", b.UserID).Updates(Booking{Status: b.Status, UpdatedAt: time.Now()}).Error
	if err != nil {
		return &Booking{}, err
	}
	if b.ID != 0 {
		err = db.Debug().Model(&User{}).Where("id = ?", b.UserID).Take(&b.User).Error
		if err != nil {
			return &Booking{}, err
		}
	}
	return b, nil
}
