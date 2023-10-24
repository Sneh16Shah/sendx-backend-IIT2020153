<h1 align='center'>Sendx-backend-IIT2020153</h1>

## About the Website

Completed the backend heavy assignment and completed the requirement given by sendx.
- Requirements:
- Pages at times may not be available so you may have to retry The application may have paying customers and non paying customers. So your server needs to always give priority for crawling to paying customers. For simplicity you can differentiate between paying and non paying customers via query parameter being passed to the backend from the frontend API call.Using go colly in golang the project will crawl the url and fetch the required details.
- Ensuring that sometimes Pages may not be available so it will retry 3 times or for paying customers 5 times.
Priority is given to the paying customers is done with the help of priority_queue which will put the paying customers to the top in the queue.

# Video Link
https://www.dropbox.com/scl/fi/8bucs7hh7tb1n2uecjlfs/Desktop-2023.10.24-21.17.01.09.mp4?rlkey=q82m63b8wg020bkal5ap12347&dl=0

## Steps to run the website

**To run the web app on your local computer, clone this repository-**

**1.Open the terminal in the client directory and run the following command :**
<br>```go mod tidy```</br>
```go run main.go```
