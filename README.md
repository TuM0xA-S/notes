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
create user		| `POST /api/user/create`
login	    	| `POST /api/user/login`
notes list  	| `GET /api/me/notes`
create note 	| `POST /api/me/notes/create`
note detail 	| `GET /api/me/notes/{note_id}`
update note 	| `PUT /api/me/notes/{note_id}`
remove note 	| `DELETE /api/me/notes/{note_id}`
user detail 	| `GET /api/me`
published notes | `GET /api/notes`

* to regiser/login provide `username` and `password`
* `title` and `body` to create note
* note can be published by setting `published` field to `true`

tasks
* [x] api base
* [ ] pagination
* [ ] filters

