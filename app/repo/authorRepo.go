package Repo

import (
	"fmt"
	. "github.com/gcamlicali/RESTHW/app/model"
	"gorm.io/gorm"
	"time"
)

type AuthorRepository struct {
	db *gorm.DB
}

//NewAuthorRepository Crate and Return new Repository
func NewAuthorRepository(db *gorm.DB) *AuthorRepository {
	return &AuthorRepository{db: db}
}

//Migrations  Author Table
func (a *AuthorRepository) Migrations() error {
	err := a.db.AutoMigrate(&Author{})
	if err != nil {
		fmt.Println("Migration Error")
		return err
	}

	return nil
}

//FillAuthorData Fills JSON data into DB
func (a *AuthorRepository) FillAuthorData() {
	fmt.Println("fill author a geldi")
	authors := GetAllAuthorsFromJson()

	for _, author := range authors {

		fmt.Println(author.AuthorID)
		fmt.Println(author.AuthorName)
		a.db.Where(Author{AuthorID: author.AuthorID}).
			Attrs(Author{AuthorID: author.AuthorID, AuthorName: author.AuthorName}).
			FirstOrCreate(&author)
	}
}

//GetAllAuthors Get All Authors
func (a *AuthorRepository) GetAllAuthors() ([]Author, error) {
	var author []Author

	result := a.db.Where("deleted_at IS NULL").Find(&author)
	if result.Error != nil {
		//log.Fatal("Db error: ", result.Error)
		return nil, result.Error
	}

	return author, nil
}

//GetAuthorByIdWithBooks Get Author by id and with books
func (a *AuthorRepository) GetAuthorByIdWithBooks(id int) ([]AuthorWithBook, error) {
	var author []AuthorWithBook

	result := a.db.Joins("left join books on authors.author_id = books.authors_id").
		Where("authors.author_id = ?", id).
		Table("authors").
		Select("books.book_id ,books.name, authors.author_name").
		Scan(&author)

	if result.Error != nil {
		return nil, result.Error
	}
	return author, nil
}

//SearchAuthorByName Get Author by name
func (a *AuthorRepository) SearchAuthorByName(name string) ([]Author, error) {
	var authors []Author

	result := a.db.Where("author_name ILIKE ? ", "%"+name+"%").Find(&authors)
	if result.Error != nil {
		return nil, result.Error
	}

	return authors, nil
}

//GetAuthorById Get By Id <SELECT * FROM Authors WHERE ID = id>
func (a *AuthorRepository) GetAuthorById(id int) (Author, error) {
	var author Author

	result := a.db.Where("id = ?", id).Find(&author)
	if result.Error != nil {
		return author, result.Error
	}

	return author, nil
}

func (a *AuthorRepository) DeleteAuthorById(id int) error {
	var author Author

	result := a.db.Where(Book{BookID: id}).Find(&author).Model(&author).Update("deleted_at", time.Now())
	if result.Error != nil {
		//log.Fatal("Delete book by id is not completed: ", result.Error)
		return result.Error
	}

	return nil
}
