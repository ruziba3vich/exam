package models

import "time"

type Author struct {
	id        int
	Name      string     `json:"name"`
	Password  string     `json:"password"`
	BirthDate *time.Time `json:"birth_date"`
	Biography string     `json:"biography"`
	createdAt time.Time
	updatedAt time.Time
}

func (a *Author) GetId() int {
	return a.id
}

func (a *Author) SetId(id int) {
	a.id = id
}

func (a *Author) GetCreatedTime() time.Time {
	return a.createdAt
}

func (a *Author) SetCreatedTime(t time.Time) {
	a.createdAt = t
}

func (a *Author) GetUpdatedTime() time.Time {
	return a.updatedAt
}

func (a *Author) SetUpdatedTime(t time.Time) {
	a.updatedAt = t
}

/*
## Authors Table Fields (Mualliflar jadvali ustunlari)

1. **author_id (Primary Key):**
   - Har bir muallif uchun noyob identifikator.
   - Type: Integer or UUID (Universally Unique Identifier).

2. **name:**
   - Muallifning ismi.
   - Type: String.

3. **birth_date:**
   - Muallifning tug'ilgan sanasi.
   - Type: Date.

4. **biography:**
   - Muallifning qisqacha tarjimai holi yoki tavsifi.
   - Type: Text.

5. **created_at:**
   - Muallif yozuvi yaratilgan vaqt tamg'asi.
   - Type: Timestamp or Datetime.

6. **updated_at:**
   - Muallif yozuvi oxirgi marta yangilangan vaqt belgisi.
   - Type: Timestamp or Datetime.

*/
