package models

import (
	"bytes"
	"errors"
	"html"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/globalsign/mgo/bson"
	"github.com/jinzhu/gorm"
)

type Post struct {
	ID          uint64    `gorm:"primary_key;auto_increment" json:"id"`
	UserID      uint32    `gorm:"not null" json:"user_id"`
	Description string    `gorm:"size:255;not null;unique" json:"description"`
	Category    string    `gorm:"size:30;not null;" json:"category"`
	ImageURL    string    `gorm:"size:255;not null;" json:"image_url"`
	Budget      int       `gorm:"size:30;not null;" json:"budget"`
	Completed   bool      `gorm:"size:255;not null;" json:"completed"`
	Latitude    float32   `gorm:"not null" json:"latitude"`
	Longitude   float32   `gorm:"not null" json:"longitude"`
	User        User      `json:"user"`
	CreatedAt   time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt   time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (p *Post) Prepare() {
	p.ID = 0
	p.Description = html.EscapeString(strings.TrimSpace(p.Description))
	p.Category = html.EscapeString(strings.TrimSpace(p.Category))
	p.ImageURL = html.EscapeString(strings.TrimSpace(p.ImageURL))
	p.Completed = false
	p.User = User{}
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()
}

func (p *Post) Validate() error {
	if p.Description == "" {
		return errors.New("Required Description")
	}
	if p.Category == "" {
		return errors.New("Required Category")
	}
	if p.UserID < 1 {
		return errors.New("Required User ID")
	}
	// if p.ImageURL == "" {
	// 	return errors.New("Required Image")
	// }
	return nil
}

//Upload a new post
func (p *Post) UploadPost(db *gorm.DB) (*Post, error) {
	var err error
	err = db.Debug().Model(&Post{}).Create(&p).Error
	if err != nil {
		return &Post{}, err
	}
	if p.ID != 0 {
		err = db.Debug().Model(&User{}).Where("id = ?", p.UserID).Take(&p.User).Error
		if err != nil {
			return &Post{}, err
		}
	}
	return p, nil
}

//Return all posts
func (p *Post) FindAllPosts(db *gorm.DB) (*[]Post, error) {
	var err error

	posts := []Post{}

	err = db.Debug().Model(&Post{}).Limit(100).Find(&posts).Error
	if err != nil {
		return &[]Post{}, err
	}

	if len(posts) > 0 {
		for i, _ := range posts {
			err := db.Debug().Model(&User{}).Where("id = ?", posts[i].UserID).Take(&posts[i].User).Error
			if err != nil {
				return &[]Post{}, err
			}
		}
	}
	return &posts, nil
}

//Update an existing post
func (p *Post) UpdateAPost(db *gorm.DB) (*Post, error) {

	var err error

	err = db.Debug().Model(&Post{}).Where("id = ?", p.ID).Updates(Post{Description: p.Description, Category: p.Category, Budget: p.Budget, Completed: p.Completed, ImageURL: p.ImageURL, UpdatedAt: time.Now()}).Error
	if err != nil {
		return &Post{}, err
	}
	if p.ID != 0 {
		err = db.Debug().Model(&User{}).Where("id = ?", p.UserID).Take(&p.User).Error
		if err != nil {
			return &Post{}, err
		}
	}
	return p, nil
}

//Delete a post
func (p *Post) DeleteAPost(db *gorm.DB, pid uint64, uid uint32) (int64, error) {

	db = db.Debug().Model(&Post{}).Where("id = ? and author_id = ?", pid, uid).Take(&Post{}).Delete(&Post{})

	if db.Error != nil {
		if gorm.IsRecordNotFoundError(db.Error) {
			return 0, errors.New("Post not found")
		}
		return 0, db.Error
	}
	return db.RowsAffected, nil
}

func UploadPostPicToS3(path string, s *session.Session, file multipart.File, fileHeader *multipart.FileHeader) (string, error) {
	urlLink := "https://vickikbt-fixit.s3.us-east-2.amazonaws.com/"

	size := fileHeader.Size
	buffer := make([]byte, size)
	file.Read(buffer)

	// create a unique file name for the file
	tempFileName := path + "/" + bson.NewObjectId().Hex() + filepath.Ext(fileHeader.Filename)

	_, err := s3.New(s).PutObject(&s3.PutObjectInput{
		Bucket:               aws.String("vickikbt-fixit"), //Bucket name
		Key:                  aws.String(tempFileName),     //File name
		ACL:                  aws.String("public-read"),    // Access type- public
		Body:                 bytes.NewReader(buffer),
		ContentLength:        aws.Int64(int64(size)),
		ContentType:          aws.String(http.DetectContentType(buffer)),
		ContentDisposition:   aws.String("attachment"),
		ServerSideEncryption: aws.String("AES256"),
		StorageClass:         aws.String("INTELLIGENT_TIERING"),
	})
	if err != nil {
		return "", err
	}

	return urlLink + tempFileName, err
}
