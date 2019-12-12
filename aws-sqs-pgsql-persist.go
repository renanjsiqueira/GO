package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	_ "github.com/lib/pq"
	uuid "github.com/satori/go.uuid"
)

const (
	QueueUrl    = "https://sqs.us-east-1.amazonaws.com/URL/Webhook.fifo"
	Region      = "us-east-1"
	CredPath    = "/home/administrator/credentials"
	CredProfile = "default"
)

type NFEHook struct {
	CnpjEmitente         string `json:"cnpj_emitente"`
	Ref                  string `json:"ref"`
	Status               string `json:"status"`
	Statussefaz          string `json:"status_sefaz"`
	Mensagemsefaz        string `json:"mensagem_sefaz"`
	Chavenfe             string `json:"chave_nfe"`
	Numero               string `json:"numero"`
	Serie                string `json:"serie"`
	Caminhoxmlnotafiscal string `json:"caminho_xml_nota_fiscal"`
	Caminhodanfe         string `json:"caminho_danfe"`
}

var (
	host     = "192.168.0.55"
	port     = 5432
	user     = "postgres"
	password = "postgres"
	dbname   = "postgres"
)

var delayIntervalconf = 30



func main() {

	fmt.Println("Sqs Read Start!!")
	for {
		go sqsRead()
		time.Sleep(time.Duration(delayIntervalconf) * time.Second)
	}

}
func sqsRead() {
	u1 := uuid.Must(uuid.NewV4())
	str := fmt.Sprintf("%s", u1)
	idHASH := strings.ToUpper(strings.Replace(str, "-", "", -1))

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	fmt.Printf("ID HASH: %s\n", idHASH)

	sess := session.New(&aws.Config{
		Region:      aws.String(Region),
		Credentials: credentials.NewSharedCredentials(CredPath, CredProfile),
		MaxRetries:  aws.Int(5),
	})

	svc := sqs.New(sess)

	// Send message
	/*send_params := &sqs.SendMessageInput{
		MessageBody:  aws.String("message body"), // Required
		QueueUrl:     aws.String(QueueUrl),       // Required
		DelaySeconds: aws.Int64(3),               // (optional) 傳進去的 message 延遲 n 秒才會被取出, 0 ~ 900s (15 minutes)
	}
	send_resp, err := svc.SendMessage(send_params)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("[Send message] \n%v \n\n", send_resp)
	*/
	// Receive message
	receive_params := &sqs.ReceiveMessageInput{
		QueueUrl:            aws.String(QueueUrl),
		MaxNumberOfMessages: aws.Int64(3),  // 一次最多取幾個 message
		VisibilityTimeout:   aws.Int64(30), // 如果這個 message 沒刪除，下次再被取出來的時間
		WaitTimeSeconds:     aws.Int64(20), // long polling 方式取，會建立一條長連線並且等在那邊，直到 SQS 收到新 message 回傳給這條連線才中斷
	}
	receive_resp, err := svc.ReceiveMessage(receive_params)
	if err != nil {
		log.Println(err)
	}
	fmt.Printf("[Receive message] \n%v \n\n", receive_resp)

	for _, message := range receive_resp.Messages {
		delete_params := &sqs.DeleteMessageInput{
			QueueUrl:      aws.String(QueueUrl),  // Required
			ReceiptHandle: message.ReceiptHandle, // Required

		}
		var nfe NFEHook
		json.Unmarshal([]byte(*message.Body), &nfe)

		sqlStatement := `
	INSERT INTO NFE (nfe_id, status, status_sefaz, mensagem_sefaz, serie, numero, cnpj_emitente, ref, chave_nfe, caminho_xml_nota_fiscal, caminho_danfe, caminho_xml_carta_correcao, numero_carta_correcao, caminho_xml_cancelamento) 
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, null, null, null);`
		_, err = db.Exec(sqlStatement, idHASH, nfe.Status, nfe.Statussefaz, nfe.Mensagemsefaz, nfe.Serie, nfe.Numero, nfe.CnpjEmitente, nfe.Ref, nfe.Chavenfe, nfe.Caminhoxmlnotafiscal, nfe.Caminhodanfe)
		if err != nil {
			panic(err)
		}

		fmt.Printf("Caminho danfe: %s, caminho xml: %s", nfe.Caminhodanfe, nfe.Caminhoxmlnotafiscal)

		_, err := svc.DeleteMessage(delete_params) // No response returned when successed.
		if err != nil {
			log.Println(err)
		}
		fmt.Printf("[Delete message] \nMessage ID: %s has beed deleted.\n\n", *message.MessageId)
	}
}
