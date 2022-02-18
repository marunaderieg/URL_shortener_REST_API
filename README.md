# Url Shortener REST API written in GoLang and SQLite

## Description

This is REST API provides two functionalities:
1. Assigning a unique short url to the passed long url. The short url will consist of 8 random characters.
2. Redirecting to the original url when loading the short url.

## Interface
### - Assign a unique short url to a long url
At http://127.0.0.1:8080/shorten use the **POST** method to pass a JSON key-value pair. Use the keyword "url", 
and set the value to the url you want to shorten. <br><br>
**Example:** You could POST the following JSON message: `{ "url" : "https://drewderieg.com/thisIsAVeryLongUrl" }`.
<br><br>**Returns:** This will return a JSON response in the following format: `{ "short_url" : "F8eiKwl9" }` 
### - Redirect to the original url
In order to redirect from a short url to the original url use the **GET** method at http://127.0.0.1:8080/{short_url}.
<br><br>**Example:** Call http://127.0.0.1:8080/F8eiKwl9 and you will be redirected to https://drewderieg.com/thisIsAVeryLongUrl.

## Installation
Since this API is deploying a SQLite database you must ensure that:
- you set the environment variable `CGO_ENABLED=1`
- you have `gcc` compile present within your path

Download the project into your $GOPATH, move to the new restAPI directory and run the command `$ go run main.go`. 
<br><br>An alternative is to build the binary by running `$ go build` in the project directory. This will create a binary named `restAPI`.

## Note
- The current database is deleted every time you rerun the server.
- You must pass the full url including the protocol when you use the http://127.0.0.1:8080/shorten POST entry point.
- This API is made to handle a decent number of "long url" - "short url" pairs. You need to modify the code if you want to create and store billions of short urls.
