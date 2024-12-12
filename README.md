# go-mysql-crud
Sample crud operation using Golang and MySql

## API ENDPOINTS

### All Characters
- Path : `/characters`
- Method: `GET`
- Response: `200`

### Pagination & limit Characters
- Path : `/characters?page=1&limit=10`
- Method: `GET`
- Response: `200`

### Search Characters
- Path : `/characters/search?name=sasuke`
- Method: `GET`
- Response: `200`

### Create Post
- Path : `/characters`
- Method: `POST`
- Response: `201`

### Details a characters
- Path : `/characters/{slug}`
- Method: `GET`
- Response: `200`

### Update characters
- Path : `/characters/{slug}`
- Method: `PUT`
- Response: `200`

### Delete characters
- Path : `/characters/{slug}`
- Method: `DELETE`
- Response: `204`
