# notes
useless, but production ready(?)

### features
* jwt authentification
* ~~nice test organisation~~(just tests)
* model/controller separation
* cool docker image(settings in .env file)

### api
action              | request
------------------- | ---------------
create user		    | `POST /api/users`
login	    	    | `POST /api/me`
notes list  	    | `GET /api/me/notes`
create note 	    | `POST /api/me/notes`
note detail(yours)  | `GET /api/me/notes/{note_id}`
update note 	    | `PUT /api/me/notes/{note_id}`
remove note 	    | `DELETE /api/me/notes/{note_id}`
user detail 	    | `GET /api/me`
published notes     | `GET /api/notes`
note detail(public) | `GET /api/notes/{note_id}`
user detail			| `GET /api/users/{user_id}`

* to regiser/login provide `username` and `password`
* `title` and `body` to create note
* note can be published by setting `published` field to `true`
* `page` query parameter for specifying page
* `no_body=true` for omit note body

### tasks
* [x] api base
* [x] pagination
* [ ] filters
* [x] rework error handling
* [ ] reorganize tests
* [x] create desktop gui ([notesclient](https://github.com/TuM0xA-S/notesclient))

