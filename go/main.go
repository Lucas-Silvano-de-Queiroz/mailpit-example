package main

import (
	"fmt"
	"net/smtp"
)

func main() {
	host := "localhost"
	port := "1025"
	address := host + ":" + port

	from := "sistema-go@exemplo.com"
	to := []string{"usuario-go@exemplo.com"}
	subject := "Teste de E-mail com Mailpit e Go 🐹"
	body := "Olá! Este é um e-mail de teste enviado usando Go e Mailpit."

	header := make(map[string]string)
	header["From"] = from
	header["To"] = to[0]
	header["Subject"] = subject
	header["Content-Type"] = "text/plain; charset=UTF-8"

	message := ""
	for k, v := range header {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + body

	err := smtp.SendMail(address, nil, from, to, []byte(message))
	if err != nil {
		fmt.Printf("Erro ao enviar e-mail: %v\n", err)
		return
	}

	fmt.Println("E-mail enviado com sucesso via Go!")
	fmt.Println("Visualize em: http://localhost:8025")
}
