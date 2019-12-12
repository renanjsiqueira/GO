package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"net/smtp"
	"encoding/json"
	"log"
	"time"
)


var cfrom = ""
var cto =""
var cpassword =""
var chost=""
var cport="" 
var cuser="" 



type SMTPconfigs struct {
	SMTPconfigs []SMTPconfig `json:"SMTPconfigs"`
}

type SMTPconfig struct {
	Fromc string `json:"fromc"`
	Toc string `json:"toc"`
	Passwordc  string `json:"passwordc"`
	Hostc string `json:"hostc"`
	Portc string `json:"portc"`
	Userc string `json:"userc"`
}


type Configs struct {
	Configs []Config `json:"config"`
}

var Sendemailconf=true
var Logconf=false
var Logdownconf=true
var Filelogconf="local"
var Filenameconf="monitoring.log"
var delayBreakconf = 3
var delayIntervalconf = 30
var CountDownconf = 1


type Config struct {
	Sendemailc bool `json:"sendemail"`
	Logc bool `json:"log"`
	Logdownc bool `json:"logdown"`
	Filelogc string `json:"filelog"`
	Filenamec string `json:"filename"`
	Intervalc int `json:"interval"`
	Delayc  int `json:"delay"`
	CountDown int `json:"countdown"`
}

type Siteconfigs struct {
	Siteconfigs []Siteconfig `json:"Siteconfig"`
}

type Siteconfig struct {
	Namesv string `json:"name"`
	URLsv string `json:"url"`

}




func main() {

	jsonFile, err := os.Open("smtp.json")
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
	}


	fmt.Println("Successfully Opened smtp.json")
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()


	// read our opened xmlFile as a byte array.
	byteValue, _ := ioutil.ReadAll(jsonFile)


	// we initialize our Users array
	var smtps SMTPconfigs


	// we unmarshal our byteArray which contains our
	// jsonFile's content into 'users' which we defined above
	json.Unmarshal(byteValue, &smtps)

	
	for i := 0; i < len(smtps.SMTPconfigs); i++ {
		fmt.Println("host: " + smtps.SMTPconfigs[i].Hostc)
		fmt.Println("port: " + smtps.SMTPconfigs[i].Portc)
		fmt.Println("password: " + smtps.SMTPconfigs[i].Passwordc)
		fmt.Println("from: " + smtps.SMTPconfigs[i].Fromc)
		fmt.Println("to: " + smtps.SMTPconfigs[i].Toc)
		cfrom =  smtps.SMTPconfigs[i].Fromc
		cto =smtps.SMTPconfigs[i].Toc
		cpassword =smtps.SMTPconfigs[i].Passwordc
		chost=smtps.SMTPconfigs[i].Hostc
		cport=smtps.SMTPconfigs[i].Portc
		cuser=smtps.SMTPconfigs[i].Userc
		
	}


	jsonFile2, err2 := os.Open("config.json")
	// if we os.Open returns an error then handle it
	if err2 != nil {
		fmt.Println(err2)
	}


	fmt.Println("Successfully Opened config.json")
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile2.Close()


	// read our opened xmlFile as a byte array.
	byteValue2, _ := ioutil.ReadAll(jsonFile2)


	// we initialize our Users array
	var configs Configs


	json.Unmarshal(byteValue2, &configs)


	for i := 0; i < len(configs.Configs); i++ {
		
		Sendemailconf =  configs.Configs[i].Sendemailc
		Logconf =  configs.Configs[i].Logc
		Logdownconf= configs.Configs[i].Logdownc
		Filelogconf= configs.Configs[i].Filelogc
		Filenameconf= configs.Configs[i].Filenamec
		delayBreakconf = configs.Configs[i].Delayc
		delayIntervalconf = configs.Configs[i].Intervalc
		CountDownconf = configs.Configs[i].CountDown
		fmt.Println("File Name LOG: ", Filenameconf )
		fmt.Println("Interval: ", delayIntervalconf )
		fmt.Println("Delay Break: " , delayBreakconf )
		fmt.Println("CountDown: " , CountDownconf )
	if Sendemailconf == true {
		fmt.Println("Send Email: true " )
	} else {
			fmt.Println("Send Email: false " )	
		}
		if Logconf == true {
			fmt.Println("Log Online: true " )
		} else {
				fmt.Println("Log Online: false " )	
			}
			if Logdownconf == true {
				fmt.Println("Log Down Service: true " )
			} else {
					fmt.Println("Log Down Service: false " )	
				}


		
		
	}
	

	for {
		go MonitoramentoJson()
	time.Sleep(time.Duration(delayIntervalconf) * time.Second)
	}
	
	
	//exibeIntroducao()
/*
	for {
		exibeMenu()

		comando := leComando()

		switch comando {
		case 1:
			for {
				go MonitoramentoJson()
			time.Sleep(10 * time.Second)
			}
		case 2:
			fmt.Println("Exibindo Logs...")
			imprimeLogs()
		case 0:
			fmt.Println("Saindo do programa")
			os.Exit(0)
		default:
			fmt.Println("Não conheço este comando")
			os.Exit(-1)
		}
	}
*/
}








func sendduo(body string,nameService string) {
	from := cfrom
	pass := cpassword
	to := cto

	msg := "From: " + from + "\n" +
		"To: " + to + "\n" +
		"Subject: Status Services Down "+nameService+"\n\n" +
		body

	err := smtp.SendMail(chost+":"+cport,
		smtp.PlainAuth("", cuser, pass, chost),
		from, []string{to}, []byte(msg))

	if err != nil {
		log.Printf("smtp error: %s", err)
		return
	}
	
	log.Print("sent, emails")
}



	

func exibeIntroducao() {
	nome := "Renan"
	versao := 1.0
	fmt.Println("Olá, sr.", nome)
	fmt.Println("Este programa está na versão", versao)
}

func exibeMenu() {
	fmt.Println("1- Iniciar Monitoramento")
	fmt.Println("2- Exibir Logs")
	fmt.Println("0- Sair do Programa")
}

func leComando() int {
	var comandoLido int
	fmt.Scan(&comandoLido)
	fmt.Println("O comando escolhido foi", comandoLido)
	fmt.Println("Comando..")

	return comandoLido
}



func MonitoramentoJson() {

	fmt.Println("Monitorando...")


	jsonFile, err := os.Open("site.json")
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
	}


	fmt.Println("Successfully Opened site.json")
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()


	// read our opened xmlFile as a byte array.
	byteValue, _ := ioutil.ReadAll(jsonFile)


	// we initialize our Users array
	var sites Siteconfigs


	// we unmarshal our byteArray which contains our
	// jsonFile's content into 'users' which we defined above
	json.Unmarshal(byteValue, &sites)

	


		for i := 0; i < len(sites.Siteconfigs); i++ {
			fmt.Println("Testando Serviço", sites.Siteconfigs[i].Namesv, ":", sites.Siteconfigs[i].URLsv)
			testaSite(sites.Siteconfigs[i].URLsv,sites.Siteconfigs[i].Namesv)
			time.Sleep(time.Duration(delayBreakconf) * time.Second)
		   
	   }
	


}

func testaSiteCountDown(site string,servicename string,erro string) {
	errors  := 1
	errormsg := erro
	for i := 1; i < CountDownconf; i++ {
		resp, err := http.Get(site)
		
			if err != nil {
				fmt.Println("Ocorreu um erro:", err)
			}
		
			if resp.StatusCode == 200 {
				fmt.Println("Site:", site, "foi carregado com sucesso!")
		
		
				if Logconf == true {
					registraLog(site, true)
				}
			} else {
				
				erro := time.Now().Format("02/01/2006 15:04:05") + " - " + "Site: "+ site +" esta com problemas. Status Code:" + strconv.Itoa(resp.StatusCode)
				errormsg += "\n"+erro
				fmt.Println("Site: ", site, "esta com problemas. Status Code:", resp.StatusCode)
				errors++;
				if Logdownconf == true {
					registraLog(site, true)
				}
		
			}
			time.Sleep(time.Duration(delayBreakconf) * time.Second)
	}
	if errors >= CountDownconf && Logdownconf == true {
		sendduo(errormsg,servicename)
	}
	
}

func testaSite(site string,servicename string) {
	resp, err := http.Get(site)

	if err != nil {
		fmt.Println("Ocorreu um erro:", err)
	}

	if resp.StatusCode == 200 {
		fmt.Println("Site:", site, "foi carregado com sucesso!")


		if Logconf == true {
			registraLog(site, true)
		}
	} else {
		erro := time.Now().Format("02/01/2006 15:04:05") + " - " + "Site: "+ site +" esta com problemas. Status Code:" + strconv.Itoa(resp.StatusCode)
		fmt.Println("Site: ", site, "esta com problemas. Status Code:", resp.StatusCode)
		if Sendemailconf == true {
			testaSiteCountDown(site,servicename,erro)
		}
		if Logdownconf == true {
			registraLog(site, true)
		}
	}
}



func leSitesDoArquivo() []string {
	var sites []string
	arquivo, err := os.Open("sites.txt")

	if err != nil {
		fmt.Println("Ocorreu um erro:", err)
	}

	leitor := bufio.NewReader(arquivo)
	
	for {
		linha, err := leitor.ReadString('\n')
		linha = strings.TrimSpace(linha)

		sites = append(sites, linha)

		if err == io.EOF {
			break
		}

	}
	
	arquivo.Close()
	return sites
}

func registraLog(site string, status bool) {

	arquivo, err := os.OpenFile(Filenameconf, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

	if err != nil {
		fmt.Println(err)
	}
	

	arquivo.WriteString(time.Now().Format("02/01/2006 15:04:05") + " - " + site + " - online: " + strconv.FormatBool(status) + "\n")

	arquivo.Close()
}

func imprimeLogs() {

	arquivo, err := ioutil.ReadFile(Filenameconf)

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(string(arquivo))

}
