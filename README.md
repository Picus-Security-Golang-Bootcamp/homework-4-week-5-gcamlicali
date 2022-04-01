## Homework | Week 5

- Creates Server on 3000 Port
- Reads all book and author data from JSON and writes DB

Query on 	
GET		"/book" 					-> Get all Books with authors
GET		"/book/name/{name}" 		-> Search Books by book name
GET 	"/book/id/{id:[0-9]+}"		-> Get books by book id
PATCH	"/book/id/{id:[0-9]+}"		-> Buy books by book id and requires quantity JSON in html body
DELETE	"/book/id/{id:[0-9]+}"		-> Soft deletes book by id

GET		"/author"					-> Get all Author with books
GET		"/author/name/{name}"		-> Search Authors by author name
GET		"/author/id/{id:[0-9]+}"	-> Get author by author id
DELETE	"/author/id/{id:[0-9]+}"	-> Soft deletes author by id
