# go-naruto-api
Sample crud operation using Golang and NoSql(MongoDB)

## API ENDPOINTS

### All Characters/Tailedbeast
- Path : `/characters` /  `/tailedbeast`
- Method: `GET`
- Response: `200`

### Pagination & limit Characters/Tailedbeast
- Path : `/characters?page=1&limit=10` / `/tailedbeast?page=1&limit=10`
- Method: `GET`
- Response: `200`

### Search Characters/Tailedbeast
- Path : `/characters/search?name=sasuke` / `/tailedbeast/search?name=kurama`
- Method: `GET`
- Response: `200`

### Create Post
- Path : `/characters` /  `/tailedbeast`
- Method: `POST`
- Response: `201`
- `https://www.postman.com/muhammadhafizhzikry/narutoapi/request/tsrtd6x/storecharacter?origin=request` / `https://www.postman.com/muhammadhafizhzikry/narutoapi/request/pbafrnl/storetailedbeast`

### Details a Characters/Tailedbeast
- Path : `/characters/{slug}` / `/tailedbeast/{slug}`
- Method: `GET`
- Response: `200`

### Update characters
- Path : `/characters/{slug}` / `/tailedbeast/{slug}`
- Method: `PUT`
- Response: `200`
- `https://www.postman.com/muhammadhafizhzikry/narutoapi/request/itakuvr/updatetailedbeast` / `https://www.postman.com/muhammadhafizhzikry/narutoapi/request/5zesfy8/updatecharacter`

### Delete characters
- Path :  `/characters/{slug}` / `/tailedbeast/{slug}`
- Method: `DELETE`
- Response: `204`
