package main

import (
	"log"
	"net/http"
	"os"
	"text/template"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/loozhengyuan/ical"
	"github.com/rs/xid"
	// "github.com/aws/aws-sdk-go/service/dynamodb"
	// "github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

func uploadToS3(bucket string, filename string) {
	// Open file based on filename
	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("Unable to open file %q, %v", filename, err)
	}

	defer file.Close()

	// Print the error and exit.
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("ap-southeast-1")},
	)
	if err != nil {
		log.Fatalf("Unable to create session object, %v", err)
	}

	// Setup the S3 Upload Manager. Also see the SDK doc for the Upload Manager
	// for more information on configuring part size, and concurrency.
	//
	// http://docs.aws.amazon.com/sdk-for-go/api/service/s3/s3manager/#NewUploader
	uploader := s3manager.NewUploader(sess)

	// Upload the file's body to S3 bucket as an object with the key being the
	// same as the filename.
	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucket),

		// Can also use the `filepath` standard library package to modify the
		// filename as need for an S3 object key. Such as turning absolute path
		// to a relative path.
		Key: aws.String(filename),

		// The file to be uploaded. io.ReadSeeker is preferred as the Uploader
		// will be able to optimize memory when uploading large content. io.Reader
		// is supported, but will require buffering of the reader's bytes for
		// each part.
		Body: file,
	})
	if err != nil {
		// Print the error and exit.
		log.Fatalf("Unable to upload %q to %q, %v", filename, bucket, err)
	}

	log.Printf("Successfully uploaded %q to %q\n", filename, bucket)
}

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	t, err := template.ParseFiles("static/" + tmpl)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = t.Execute(w, p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func createEvent(eventName string) error {
	// Create new calendar object
	c := ical.NewCalendar()

	// Create new event object
	e := ical.NewEvent()
	e.SUMMARY = eventName
	timeNow := time.Now()
	e.DTSTART = &timeNow

	// Add event to calendar object
	c.EVENT = append(c.EVENT, *e)

	// Export to file
	o := c.GenerateCalendarProp()
	guid := xid.New()
	filename := guid.String() + ".ics"
	ical.OutputToFile(filename, []byte(o), 0644)
	uploadToS3("bookyourtime-development", filename)

	return nil
}
