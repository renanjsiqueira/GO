package main

import (
  "database/sql"
  "crypto/sha256"
  _"github.com/lib/pq"
  "encoding/hex"
  "io/ioutil"
	"fmt"
	"os"
	"encoding/json"
)

var (
  host     = "192.168.0.57"
  port     = 5432
  user     = "postgres"
  password = "postgres"
  dbname   = "yourdb"
  PassKey = "YourEndKey";
  PassKeyInit = "YourInitKey";
  limit = 1000;
)

type Serviceconfigs struct {
	Serviceconfigs []Serviceconfig `json:"Serviceconfigs"`
}

type Serviceconfig struct {
	Jhost string `json:"jhost"`
	Jport int `json:"jport"`
	Juser  string `json:"juser"`
	Jpassword string `json:"jpassword"`
	Jdbname string `json:"jdbname"`
	JPassKey string `json:"jPassKey"`
  JPassKeyInit string `json:"jPassKeyInit"`
  Jlimit int `json:"jlimit"`
}


  
func CToGoString(c [32]byte) string {
    n := -1
    for i, b := range c {
        if b == 0 {
            break
        }
        n = i
    }
    return string(c[:n+1])
}

func main() {

  jsonFile, err := os.Open("services.json")
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
	}


	fmt.Println("Successfully Opened services.json")
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()


	// read our opened xmlFile as a byte array.
	byteValue, _ := ioutil.ReadAll(jsonFile)


	// we initialize our Users array
	var serviceconfigs Serviceconfigs


	// we unmarshal our byteArray which contains our
	// jsonFile's content into 'users' which we defined above
	json.Unmarshal(byteValue, &serviceconfigs)





	
	for i := 0; i < len(serviceconfigs.Serviceconfigs); i++ {
		fmt.Println("host: " , serviceconfigs.Serviceconfigs[i].Jhost)
    fmt.Println("port: " , serviceconfigs.Serviceconfigs[i].Jport)
    fmt.Println("user: " , serviceconfigs.Serviceconfigs[i].Juser)
		fmt.Println("password: " , serviceconfigs.Serviceconfigs[i].Jpassword)
		fmt.Println("dbname: " , serviceconfigs.Serviceconfigs[i].Jdbname)
    fmt.Println("PassKey: " , serviceconfigs.Serviceconfigs[i].JPassKey)
    fmt.Println("PassKeyInit: " , serviceconfigs.Serviceconfigs[i].JPassKeyInit)
    fmt.Println("limitquery: " , serviceconfigs.Serviceconfigs[i].Jlimit)
    host = serviceconfigs.Serviceconfigs[i].Jhost
    port = serviceconfigs.Serviceconfigs[i].Jport
    user = serviceconfigs.Serviceconfigs[i].Juser
		password = serviceconfigs.Serviceconfigs[i].Jpassword
		dbname= serviceconfigs.Serviceconfigs[i].Jdbname
    PassKey = serviceconfigs.Serviceconfigs[i].JPassKey
    PassKeyInit= serviceconfigs.Serviceconfigs[i].JPassKeyInit
    limit=serviceconfigs.Serviceconfigs[i].Jlimit
		
  }
  
  psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
    "password=%s dbname=%s sslmode=disable",
    host, port, user, password, dbname)
  db, err := sql.Open("postgres", psqlInfo)
  if err != nil {
    panic(err)
  }

  
  rows, err := db.Query("SELECT id, password FROM user where passwordcrypt is null LIMIT $1", limit)
  if err != nil {
    // handle this error better than this
    panic(err)
  }
  defer rows.Close()
  for rows.Next() {
    var id int
    var password string
    err = rows.Scan(&id, &password)
    if err != nil {
      // handle this error
      panic(err)
    }

	sqlStatemenup := `
	  update loja.user u set passwordcrypt=$2 where u.id=$1
	   RETURNING id, passwordcrypt;`
	   var passwordcrypt string
	   var id2 int
	   pass := [] byte(PassKeyInit + password + PassKey);
	   h := sha256.New()
	   h.Write(pass)
		   sha1_hash := hex.EncodeToString(h.Sum(nil))

	   passwordcrypt = sha1_hash
	   err = db.QueryRow(sqlStatemenup, id, passwordcrypt).Scan(&id2, &passwordcrypt)
	   if err != nil {
		panic(err)
	   }
	   fmt.Println(id2, passwordcrypt)
  }
  // get any error encountered during iteration
  err = rows.Err()
  if err != nil {
    panic(err)
  }
  
  defer db.Close()

  err = db.Ping()
  if err != nil {
    panic(err)
  }

  fmt.Println("Successfully connected!")
}
