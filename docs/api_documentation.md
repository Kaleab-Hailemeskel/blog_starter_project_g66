End points
1. Registration

Endpoint: /registration
Method: POST
Body: {"email": "string", "password", "username": "string", "personal_bio": "string", "Role": "string"}
Response 201 Created  "message": "Please enter your otp to successfully register"  
400 Bad Request    "error": "user already exists"

Endpoint: /registration/verification
Body: {"email": "string", "otp": "string"}
Response 200 OK  "message": "User verified successfully"  
400 Bad Request  "message": "Invalid or expired OTP"


2.Login

Endpoint: /login
Method: POST
Body: {"email": "string", "password": "string"}
Response 200 OK  
"message": "User logged in successfully",
    "token": {
        "AccessToken": "string"
        "RefreshToken": "string"
    }


3.Forgot password


Endpoint: /forgot_password
Method: POST
Body: {"email":"string"}
Response 200 OK
    "message": "Password reset link sent"

Response 400 Bad Request
    "error": "user not found"


4.Reset password


Endpoint: /reset_password
Body: {"Token": "string", "NewPassword": "string"}
Method: PUT
Response 200 Ok
    "message": "Password reset link sent"


5.Edit profile


Endpoint: /user/edit_profile
Body: {"UserName": "string", "PersonalBio" :   "string", "ProfilePic" : "string", "Email": "string", "PhoneNum": "string", "TelegramHandle": "string", "Password":  "string"}
Authorization: Bearer <your_token>
Method: PUT
Response 200 OK
 "message": "Profile updated", "user":{the updated information}


6.Promote user 


Endpoint: /promote_user
Body: {
   "target_email": string
}
Authorization: Bearer <your_token>
Method:POST
Response 200 OK
  "message": "User promoted to ADMIN successfully"


7. Demote user


Endpoint: /demote_user
Body:  {
   "target_email": string
}
Authorization: Bearer <your_token>
Method:POST
Response 200 OK
 "message": "Admin demoted to user successfully"
8.


9.logout


Endpoint: /logout
Body:
Method: POST
Response 200 OK
 "message": "logged out successfully"


10.Create blog


Endpoint: /blog
Body: {"content": string, "description":string, "tag":string, "title":string, "author": string}
Authorization: Bearer <your_token>
Method: POST
Response 201 CREATED
"blog"{blog is displayed}, 
"message": "blog created"


11.Get blog


Endpoint: /blog, /blog?tag=something, /blog?author= ....
Body:
Authorization: Bearer <your_token>
Method: GET

Response 200 OK
"result" : [
        {
            "Title": "MongoDB Integration in Go",
            "Tags": [
                "mongodb",
                "go",
                "database"
            ],
            "Author": "",
            "Description": "Getting started with MongoDB integration in Go projects.",
            "LastUpdate": "2025-08-07T10:48:06.348Z"
        },
        .....
]

all blogs are returned


12.Edit blog


Endpoint:/blog/:id
Body:{
    "content": string,
    "description": string,
    "tags": string or array of strings,
    "title": string,
    "author":string
}
Authorization: Bearer <your_token>
Method: PUT
Response 202 Accepted
"message": "blog Updated"


13.Delete blog


Endpoint:/blog/:id
Body:
Authorization: Bearer <your_token>
Method: DELETE
Response 202 Accepted
 "message": "Blog Deleted"

 Response 401 unauthorized
  "error": "Invalid token"


14.Like blog


Endpoint:
Body:
Method:
Response


15. Dislike blog


Endpoint:
Body:
Method:
Response


16.Comment blog


Endpoint:
Body:
Method:
Response

17. AI routes

Endpoint: /ai/comment
Body:
Method: GET
Response 

Endpoint: /ai/:id
Body:
Method: GET
Response

Endpoint: /ai/filter
Body:
Method: GET
Response





