package models

import (
	"errors"
	"time"

	"github.com/jinzhu/gorm"
)

type Work struct {
	ID        uint32    `gorm:"primary_key;auto_increment" json:"id"`
	UserID    uint32    `gorm:"not null" json:"user_id"`
	WorkerID  uint32    `gorm:"not null" json:"worker_id"`
	PostID    uint32    `gorm:"not null" json:"post_id"`
	User      User      `json:"user"`
	Worker    User      `json:"worker"`
	Post      Post      `json:"post"`
	Status    string    `gorm:"not null" json:"status"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (w *Work) Prepare() {
	w.ID = 0
	w.User = User{}
	w.Worker = User{}
	w.Post = Post{}
	w.Status = "In-Progress"
	w.CreatedAt = time.Now()
	w.UpdatedAt = time.Now()
}

func (w *Work) Validate() error {
	if w.UserID < 1 {
		return errors.New("Required User ID")
	}
	if w.WorkerID < 1 {
		return errors.New("Required Worker ID")
	}
	if w.PostID < 1 {
		return errors.New("Required Post ID")
	}
	if w.Status == "" {
		return errors.New("Required Status")
	}
	return nil
}

//Upload a new work
func (w *Work) UploadWork(db *gorm.DB) (*Work, error) {
	var err error
	err = db.Debug().Model(&Work{}).Create(&w).Error
	if err != nil {
		return &Work{}, err
	}

	if w.ID != 0 {
		err = db.Debug().Model(&User{}).Where("id = ?", w.UserID).Take(&w.User).Error
		if err != nil {
			return &Work{}, err
		}

		err = db.Debug().Model(&User{}).Where("id = ?", w.WorkerID).Take(&w.Worker).Error
		if err != nil {
			return &Work{}, err
		}

		err = db.Debug().Model(&Post{}).Where("id = ?", w.PostID).Take(&w.Post).Error
		if err != nil {
			return &Work{}, err
		}
	}

	return w, nil
}

//Get work for a particular post.
func (w *Work) FindWorkByID(db *gorm.DB, pid uint64) (*Work, error) {
	var err error

	err = db.Debug().Model(&Work{}).Where("post_id = ?", pid).Take(&w).Error
	if err != nil {
		return &Work{}, err
	}

	if w.ID != 0 {
		err = db.Debug().Model(&User{}).Where("id = ?", w.UserID).Take(&w.User).Error
		if err != nil {
			return &Work{}, err
		}

		err = db.Debug().Model(&User{}).Where("id = ?", w.WorkerID).Take(&w.Worker).Error
		if err != nil {
			return &Work{}, err
		}

		err = db.Debug().Model(&Post{}).Where("id = ?", w.PostID).Take(&w.Post).Error
		if err != nil {
			return &Work{}, err
		}
	}

	return w, nil
}

//Update an existing work
func (w *Work) UpdateWork(db *gorm.DB) (*Work, error) {

	var err error

	err = db.Debug().Model(&Work{}).Where("id = ?", w.ID).Updates(Work{Status: w.Status, UpdatedAt: time.Now()}).Error
	if err != nil {
		return &Work{}, err
	}

	if w.ID != 0 {
		err = db.Debug().Model(&User{}).Where("id = ?", w.UserID).Take(&w.User).Error
		if err != nil {
			return &Work{}, err
		}

		err = db.Debug().Model(&User{}).Where("id = ?", w.WorkerID).Take(&w.Worker).Error
		if err != nil {
			return &Work{}, err
		}

		err = db.Debug().Model(&Post{}).Where("id = ?", w.PostID).Take(&w.Post).Error
		if err != nil {
			return &Work{}, err
		}
	}

	return w, nil
}

//Find work based on user id
func (w *Work) FindUserWorks(db *gorm.DB, pid uint64) (*[]Work, error) {
	var err error

	work := []Work{}

	err = db.Debug().Model(&Work{}).Order("created_at desc").Limit(100).Find(&work).Error
	if err != nil {
		return &[]Work{}, err
	}

	if len(work) > 0 {

		for i := 0; i <= len(work); i++ {
			err = db.Debug().Model(&Work{}).Limit(100).Where("worker_id = ?", pid).Find(&work).Error
			if err != nil {
				return &[]Work{}, err
			}

			for i, _ := range work {
				err := db.Debug().Model(&User{}).Where("id = ?", work[i].UserID).Take(&work[i].User).Error
				if err != nil {
					return &[]Work{}, err
				}
			}

			for i, _ := range work {
				err := db.Debug().Model(&User{}).Where("id = ?", work[i].WorkerID).Take(&work[i].Worker).Error
				if err != nil {
					return &[]Work{}, err
				}
			}

			for i, _ := range work {
				err := db.Debug().Model(&Post{}).Where("id = ?", work[i].PostID).Take(&work[i].Post).Error
				if err != nil {
					return &[]Work{}, err
				}
			}

		}

	}

	return &work, err
}
