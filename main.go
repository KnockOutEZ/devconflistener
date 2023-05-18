package main

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/joho/godotenv"
	"github.com/sfreiberg/gotwilio"
)

//go:embed templates/*
var resources embed.FS

var t = template.Must(template.ParseFS(resources, "templates/*"))

type message struct {
	Content string
	To      string
	Medium  string
}

type Config struct {
	TwilioAccountSid string
	TwilioAuthToken  string
	SenderPhone      string
	RecieverPhone    string
}

var cfg = Config{
	TwilioAccountSid: goDotEnvVariable("TWILIO_ACCOUNT_SID"),
	TwilioAuthToken:  goDotEnvVariable("TWILIO_AUTH_TOKEN"),
	SenderPhone:      goDotEnvVariable("SENDER_PHONE"),
	RecieverPhone:    goDotEnvVariable("RECEIVER_PHONE"),
}

const (
	targetURL = "https://devconfbd.com/"
	// targetURL   = "http://localhost:4000/"
	// targetClass = ".test"
	targetClass = ".btn.btn-primary.btn-md.bg-gradient-to-r"
)

var infoLog *log.Logger
var errorLog *log.Logger

func main() {
	infoFile, err := os.OpenFile("tmp/info.log", os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Fatal(err)
	}

	err = infoFile.Truncate(0)
	if err != nil {
		log.Fatal(err)
	}

	errFile, err := os.OpenFile("tmp/error.log", os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Fatal(err)
	}

	err = errFile.Truncate(0)
	if err != nil {
		log.Fatal(err)
	}

	defer infoFile.Close()
	defer errFile.Close()
	infoLog = log.New(infoFile, "INFO\t", log.Ldate|log.Ltime)
	errorLog = log.New(errFile, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	
	go monitorChanges()

	go func() {
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			data := map[string]string{
				"Region": os.Getenv("FLY_REGION"),
			}
	
			t.ExecuteTemplate(w, "index.html.tmpl", data)
		})
	
		http.HandleFunc("/info", infoHandler)
		http.HandleFunc("/error", errorHandler)
		log.Fatal(http.ListenAndServe(":"+goDotEnvVariable("PORT"), nil))
	}()

	select {}
}

func monitorChanges() {
	// previousHTML := `Latest Snippets`
	previousHTML := `Registration will open soon<i class="fa fa-arrow-right ml-1.5 text-sm"></i>`

	for {
		html, err := fetchHTML(targetURL)
		if err != nil {
			errorLog.Printf("Error fetching HTML: %s", err)
			continue
		}

		doc, err := goquery.NewDocumentFromReader(html)
		if err != nil {
			errorLog.Printf("Error creating document: %s", err)
			continue
		}

		div := doc.Find(targetClass)
		if div.Length() == 0 {
			errorLog.Printf("Target div not found")
			// continue
		}

		currentHTML, _ := div.Html()
		if currentHTML == previousHTML {
			infoLog.Println(time.Now(), "Target div has not changed!")
		} else {
			infoLog.Println(time.Now(), "Target div has changed!")
			infoLog.Println("Div content:", currentHTML)
			newSms := message{
				To:      cfg.RecieverPhone,
				Content: "Brooo devconf ticket is selling now. Go and grab it!!",
				Medium:  "whatsapp",
			}
			twilio := gotwilio.NewTwilioClient(cfg.TwilioAccountSid, cfg.TwilioAuthToken)
			sendingSms(newSms, twilio)
			os.Exit(0)
		}
		time.Sleep(1 * time.Minute)
	}
}

func fetchHTML(url string) (io.Reader, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(resp.Body)
	if err != nil {
		return nil, err
	}

	return buf, nil
}

func sendingSms(smsParam message, twilio *gotwilio.Twilio) {
	if smsParam.Medium == "whatsapp" {
		if _, err := smsParam.SendWhatsapp(twilio); err != nil {
			errorLog.Println("whatsapp sending error: ...", err.Error())
		}
		infoLog.Println("whatsapp message sent... ")
	} else {
		if _, err := smsParam.SendSms(twilio); err != nil {
			errorLog.Println("sms sending error: ...", err.Error())
		}
		
		infoLog.Println("sms message sent... ")
	}
}

// SendWhatsapp - Sends Whatsapp text message
func (m *message) SendWhatsapp(credentials *gotwilio.Twilio) (*gotwilio.SmsResponse, error) {
	resp, _, err := credentials.SendWhatsApp(cfg.SenderPhone, m.To, m.Content, "", "")
	if err != nil {
		return nil, err
	}
	return resp, nil
} //SendSms - Sends Sms text message
func (m *message) SendSms(credentials *gotwilio.Twilio) (*gotwilio.SmsResponse, error) {
	resp, _, err := credentials.SendSMS(cfg.SenderPhone, m.To, m.Content, "", "")
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func goDotEnvVariable(key string) string {

	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}

func infoHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "tmp/info.log")
}		

func errorHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "tmp/error.log")
}