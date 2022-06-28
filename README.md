## go, echo, sqlite3 rest-api examples

### API List
=========

Name              Url                   Method
----------------------------------------------
Get Member List   /api/v1/members       GET	                                        
  - Request Sample
  - Response Sample
    200 [{"id":"hong", "name":"Gildong Hong", "email":"gildonhong@mail.comn"},{},...]
    404

Create Member     /api/v1/members/:id   POST
  - Request Sample
    {"id":"hong", "name":"Gildong Hong", "email":"gildonhong@mail.comn","favorites":["watching","food","sleeping"]}	
  - Response Sample
    201 {"result":"ok"}
    400
    500

Modify Member     /api/v1/members/:id	  PUT
  - Request Sample
    {"name":"gdong Hong", "email":"gildonhong@mail.comn","favorites":["watching","food","sleeping"]}	
  - Response Sample
    200 {"result":"ok"}

204

Delete Member     /api/v1/members/:id	  DELETE
  - Request Sample
  - Response Sample
    200 {"result":"ok"}
    204
    500

Get Member        /api/v1/members/:id	  GET
  - Request Sample
  - Response Sample
    200 {"id":"hong", "name":"홍길동", "email":"gildonhong@mail.comn","favorites":["watching","food","sleeping"]}
    404



### Test
====
VSCode Plugin : Thunder Client
import /test/api_go_homework.json
