package main

import (
	"bufio"
	"fmt"
	"strings"
)

func main() {
	conn, err := tlsConnect()
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)

	var authdata string

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			return
		}

		line = strings.TrimSpace(line)
		args := strings.Split(line, " ")
		cmd := args[0]

		switch cmd {

		case "HELO":
			writer.WriteString("EHLO\n")
			writer.Flush()

		case "ERROR":
			fmt.Println("ERROR:", strings.Join(args[1:], " "))
			return

		case "POW":
			authdata = args[1]
			difficulty := atoi(args[2])
			fmt.Println(authdata, difficulty)
			suffix := solvePOW(authdata, difficulty)
			writer.WriteString(suffix + "\n")
			writer.Flush()

		case "END":
			writer.WriteString("OK\n")
			writer.Flush()
			fmt.Println("DONE!")
			return

		case "NAME":
			writer.WriteString(sha1hex(authdata+args[1]) + " Nikhil Ivannan\n")
			writer.Flush()

		case "MAILNUM":
			writer.WriteString(sha1hex(authdata+args[1]) + " 1\n")
			writer.Flush()

		case "MAIL1":
			writer.WriteString(sha1hex(authdata+args[1]) + " nickonicko779@gmail.com\n")
			writer.Flush()

		case "MAIL2":
			writer.WriteString(sha1hex(authdata+args[1]) + " my.name2@example.com\n")
			writer.Flush()

		case "SKYPE":
			writer.WriteString(sha1hex(authdata+args[1]) + " nickonicko779@gmail.com\n")
			writer.Flush()

		case "BIRTHDATE":
			writer.WriteString(sha1hex(authdata+args[1]) + " 20.08.2003\n")
			writer.Flush()

		case "COUNTRY":
			writer.WriteString(sha1hex(authdata+args[1]) + " India\n")
			writer.Flush()

		case "ADDRNUM":
			writer.WriteString(sha1hex(authdata+args[1]) + " 2\n")
			writer.Flush()

		case "ADDRLINE1":
			writer.WriteString(
				sha1hex(authdata+args[1]) +
					" E105 Isha Gayatri T Ponnambalam Salai Gerugambakkam Main Road\n",
			)
			writer.Flush()

		case "ADDRLINE2":
			writer.WriteString(
				sha1hex(authdata+args[1]) + " Chennai 600128\n",
			)
			writer.Flush()
		}
	}
}
