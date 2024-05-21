package models

import "time"

/// https://en.wikipedia.org/wiki/ISBN

type Book struct {
	id              int       //
	Title           string    `json:"title"`
	authorId        int       //
	PublicationDate time.Time `json:"publication_date"`
	isbn            string    //
	Description     string    `json:"description"`
	createdAt       time.Time //
	updatedAt       time.Time //
	isDeleted       bool
}

func (b *Book) GetId() int {
	return b.id
}

func (b *Book) SetDeleted() {
	b.isDeleted = true
}

func (b *Book) SetAuthor(authorId int) {
	b.authorId = authorId
}

func (b *Book) GetAuthorId() int {
	return b.authorId
}

func (b *Book) GetIsbn() string {
	return b.isbn
}

func (b *Book) SetIsbn(isbn string) {
	b.isbn = isbn
}

func (b *Book) GetCreatedTime() time.Time {
	return b.createdAt
}

func (b *Book) SetCreatedTime(t time.Time) {
	b.createdAt = t
}

func (b *Book) GetUpdatedTime() time.Time {
	return b.updatedAt
}

func (b *Book) SetUpdatedTime(t time.Time) {
	b.updatedAt = t
}

/*
## Books Table Fields (Kitoblar jadvali ustunlari)

1. **book_id (Primary Key):**
   - Har bir kitob uchun noyob identifikator.
   - Type: Integer or UUID.

2. **title:**
   - Kitobning nomi.
   - Type: String.

3. **author_id (Foreign Key):**
   - Kitobni yozgan muallifga havola.
   - Type: Integer or UUID.

4. **publication_date:**
   - Kitobning nashr etilgan sanasi.
   - Type: Date.

5. **isbn:**
   - Kitobning xalqaro standart kitob raqami (ISBN).
   - Type: String.

6. **description:**
   - Kitobning qisqacha tavsifi yoki xulosasi.
   - Type: Text or String.

7. **created_at:**
   - Kitob yozuvi yaratilgan vaqt tamg'asi.
   - Type: Timestamp or Datetime.

8. **updated_at:**
   - Kitob yozuvi oxirgi marta yangilangan vaqt belgisi.
   - Type: Timestamp or Datetime.

*/
