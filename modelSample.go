package main

import "fmt"


type Pessoa struct{
	name string
	idade int
}

func getName(p Pessoa) string {
	return p.name;
}

func getIdade(p Pessoa) int {
	return p.idade;
}

func main(){
	 p := Pessoa{"renan testes",20}

	fmt.Printf(getName(p), " " ,getIdade(p))

}