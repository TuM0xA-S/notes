# notes
useless, but production ready(?)

### features
* jwt authentification
* nice test organisation
* model/controller separation
* cool docker image(settings in .env file)

### api
action      | url
----------- | ---------------
create user | `/api/user/create`
login	    | `/api/user/login`
create note | `/api/me/notes/create`
note detail | `/api/me/notes/{note_id}`
remove note | `/api/me/notes/{note_id}/remove`
user detail | `/api/me`

* to regiser/login provide username and password
* title and body to create note
