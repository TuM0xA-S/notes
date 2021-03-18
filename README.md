# notes
useless, but production ready(?)

### features
* jwt authentification
* nice test organisation
* model/controller separation
* cool docker image(settings in .env file)

### api
action      | request
----------- | ---------------
create user | `POST /api/user/create`
login	    | `POST /api/user/login`
create note | `GET /api/me/notes`
create note | `POST /api/me/notes/create`
note detail | `GET /api/me/notes/{note_id}`
remove note | `POST /api/me/notes/{note_id}/remove`
user detail | `GET /api/me`

* to regiser/login provide username and password
* title and body to create note
