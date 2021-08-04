package handlers

import (
	"io/ioutil"
	"log"
	"net/http"
    "bytes"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

const object = `
{
	"name": "functions_demo",
	"description": "this is a test object",
	"ext": [
		"platform V",
		"functions"
	]
}
`

func MakeInfoHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Body != nil {
			defer r.Body.Close()
		}

	region := "test"
	endpoint := "obs.ru-moscow-1.hc.sbercloud.ru"
	objectKey := "functions_test_json"
	bucket := "cpm-functions-demo"
	
	sess, err := session.NewSession(&aws.Config{
		Credentials: credentials.NewStaticCredentials("RONNJFJCQVISHLP6RK6W", "K08tc6dNegACi09KgVsqCyOhZPAKEq6umtcDhziu", ""),
		Region:      &region,
		Endpoint:    &endpoint,
	})
	if err != nil {
		log.Printf("error creating storage client: %s\n", err)
	}

	c := s3.New(sess)

	_, err = c.PutObject(&s3.PutObjectInput{
		Key:    &objectKey,
		Bucket: &bucket,
		Body:   bytes.NewReader([]byte(object)),
	})
	if err != nil {
		log.Printf("put object error: %s\n", err)
	}

	object, err := c.GetObject(&s3.GetObjectInput{
		Bucket: &bucket,
		Key:    &objectKey,
	})
	if err != nil {
		log.Printf("get object error: %s", err)
	}

	b, err := ioutil.ReadAll(object.Body)
	if err != nil {
		log.Printf("read object body error: %s", err)
	}

	log.Printf("object: %s\n", string(b))

		w.WriteHeader(http.StatusOK)

		_, err = w.Write(b)
		if err != nil {
			log.Printf("error write answer - %v", err)
		}
	}
}
