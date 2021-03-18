package models

import (
	"errors"
	"html"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
)

type Review struct {
	ID        uint64    `gorm:"primary_key;auto_increment" json:"id"`
	UserID    uint32    `gorm:"not null;unique" json:"user_id"`
	WorkerID  uint32    `gorm:"not null;unique" json:"worker_id"`
	User      User      `json:"user"`
	Rating    uint32    `gorm:"not null" json:"rating"`
	Comment   string    `gorm:"not null" json:"comment"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (r *Review) Prepare() {
	r.ID = 0
	r.User = User{}
	//r.Rating = 0
	r.Comment = html.EscapeString(strings.TrimSpace(r.Comment))
	r.CreatedAt = time.Now()
	r.UpdatedAt = time.Now()
}

func (r *Review) Validate() error {
	if r.UserID < 1 {
		return errors.New("Required user id")
	}
	if r.WorkerID < 1 {
		return errors.New("Required worker id")
	}
	if r.Rating < 0 {
		return errors.New("Required rating")
	}
	if r.Comment == "" {
		return errors.New("Required comment")
	}
	return nil
}

//Upload a new rating
func (r *Review) UploadReview(db *gorm.DB) (*Review, error) {
	var err error
	err = db.Debug().Model(&Review{}).Create(&r).Error
	if err != nil {
		return &Review{}, err
	}

	if r.ID != 0 {
		err = db.Debug().Model(&User{}).Where("id=?", r.WorkerID).Take(&r.User).Error
		if err != nil {
			return &Review{}, err
		}
	}
	return r, nil
}

//Return all reviews
func (r *Review) FindAllReviews(db *gorm.DB) (*[]Review, error) {
	var err error

	reviews := []Review{}

	err = db.Debug().Model(&Post{}).Order("created_at desc").Limit(100).Find(&reviews).Error
	if err != nil {
		return &[]Review{}, err
	}

	if len(reviews) > 0 {
		for i, _ := range reviews {
			err := db.Debug().Model(&User{}).Where("id=?", reviews[i].UserID).Take(&reviews[i].User).Error
			if err != nil {
				return &[]Review{}, err
			}
		}
	}
	return &reviews, err
}

//Return all reviews for specific user
func (r *Review) FindUserReviews(db *gorm.DB, pid uint64) (*[]Review, error) {
	var err error
	//var ratings uint32

	reviews := []Review{}

	err = db.Debug().Model(&Review{}).Order("created_at desc").Limit(100).Where("user_id=?", pid).Find(&reviews).Error
	if err != nil {
		return &[]Review{}, err
	}

	if len(reviews) > 0 {
		for i, _ := range reviews {
			err := db.Debug().Model(&User{}).Where("id=?", reviews[i].UserID).Take(&reviews[i].User).Error
			if err != nil {
				//ratings += reviews[i].Rating
				//fmt.Println(ratings)
				return &[]Review{}, err
			}
		}
	}
	return &reviews, nil
}

//Update an existing review
func (r *Review) UpdateReview(db *gorm.DB) (*Review, error) {
	var err error

	err = db.Debug().Model(&Review{}).Where("id=?", r.ID).Updates(Review{Rating: r.Rating, Comment: r.Comment}).Error
	if err != nil {
		return &Review{}, err
	}
	if r.ID != 0 {
		err = db.Debug().Model(&User{}).Where("id=?", r.UserID).Take(&r.User).Error
		if err != nil {
			return &Review{}, err
		}
	}
	return r, nil
}

//Delete a review
func (r *Review) DeleteReview(db *gorm.DB, pid uint64, uid uint32) (int64, error) {

	db = db.Debug().Model(&Review{}).Where("id = ? and user_id = ?", pid, uid).Take(&Review{}).Delete(&Review{})

	if db.Error != nil {
		if gorm.IsRecordNotFoundError(db.Error) {
			return 0, errors.New("Review not found")
		}
		return 0, db.Error
	}
	return db.RowsAffected, nil
}
