CREATE TABLE IF NOT EXISTS Authors (
    id SERIAL PRIMARY KEY,
    name VARCHAR(64),
    password VARCHAR(255),
    biography TEXT,
    birth_date DATE,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS isbn (
    id SERIAL PRIMARY KEY,
    name VARCHAR(20) UNIQUE NOT NULL
);

CREATE TABLE IF NOT EXISTS Books (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255),
    author_id INTEGER REFERENCES Authors(id),
    publication_date DATE,
    isbn INTEGER REFERENCES ISBN(id),
    description TEXT,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);


-- 1. **book_id (Primary Key):**
--    - Har bir kitob uchun noyob identifikator.
--    - Type: Integer or UUID.

-- 2. **title:**
--    - Kitobning nomi.
--    - Type: String.

-- 3. **author_id (Foreign Key):**
--    - Kitobni yozgan muallifga havola.
--    - Type: Integer or UUID.

-- 4. **publication_date:**
--    - Kitobning nashr etilgan sanasi.
--    - Type: Date.

-- 5. **isbn:**
--    - Kitobning xalqaro standart kitob raqami (ISBN).
--    - Type: String.

-- 6. **description:**
--    - Kitobning qisqacha tavsifi yoki xulosasi.
--    - Type: Text or String.

-- 7. **created_at:**
--    - Kitob yozuvi yaratilgan vaqt tamg'asi.
--    - Type: Timestamp or Datetime.

-- 8. **updated_at:**
--    - Kitob yozuvi oxirgi marta yangilangan vaqt belgisi.
--    - Type: Timestamp or Datetime.

