package main

import (
	"bufio"
	"bytes"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"log"
	"os"
    "os/user"
	"strings"
    "crypto/aes"
    "crypto/cipher"
    "crypto/sha256"
    "path"
)

func main() {

    usr, err := user.Current()
    if err != nil {
        log.Fatal( err )
    }
    
    database_dir := path.Join(usr.HomeDir, ".secret_storage")

    _ = os.Mkdir(database_dir, 0755)


	reader := bufio.NewReader(os.Stdin)
	var directories []string

	key_value_pairs := make(map[string]string)
    nonce := []byte("eeeeeeeeeeee")

	for {

		fmt.Printf(">>")
		input_data, _, err := reader.ReadLine()

		if err != nil {
			panic("fail to read user input.")
		}

		command := strings.Fields(string(input_data))

		var comm string
		var args []string

		if len(command) > 0 {

			args = command[1:]
			comm = command[0]

		}
		switch comm {
		case "ls":
			{
				for _, dir := range directories {
					fmt.Println(dir)
				}

			}

		case "new":
			{
				if len(args) == 1 {
					dir_name := args[0]
					directories = append(directories, dir_name)
				} else {
					fmt.Println("new command expects 1 argument")
				}

			}

		case "add":
			if len(args) == 2 {
				key := args[0]
				value := args[1]
				key_value_pairs[key] = value
			} else {
				fmt.Println("add command expects 2 argument")
			}

		case "list":
			for key, value := range key_value_pairs {
				fmt.Printf("[%s:%s]\n", key, value)
			}

		case "save":
			var buffer bytes.Buffer
			enc := gob.NewEncoder(&buffer)

			if err := enc.Encode(key_value_pairs); err != nil {
				log.Fatal(err)
			}


            password, _, err := reader.ReadLine() 


            _key := sha256.Sum256(password)
            key := _key[:]

			block, err := aes.NewCipher(key)
			if err != nil {
				panic(err.Error())
			}

			aesgcm, err := cipher.NewGCM(block)
			if err != nil {
				panic(err.Error())
			}

			ciphertext := aesgcm.Seal(nil, nonce, buffer.Bytes(), nil)
                
 
            file_path := path.Join(database_dir, "tmp")


			if err := ioutil.WriteFile(file_path, ciphertext, 0644); err != nil {
				log.Fatal(err)
			}

		case "load":
 
            file_path := path.Join(database_dir, "tmp")


			data, err := ioutil.ReadFile(file_path)

            password, _, err := reader.ReadLine() 


            _key := sha256.Sum256(password)
            key := _key[:]

			block, err := aes.NewCipher(key)
			if err != nil {
				panic(err.Error())
			}

			aesgcm, err := cipher.NewGCM(block)
			if err != nil {
				panic(err.Error())
			}

			plaintext, err := aesgcm.Open(nil, nonce, data, nil)
			if err != nil {
                fmt.Println("password is not correct.")
                continue
			}

			buffer := bytes.NewBuffer(plaintext)
			if err != nil {
			}
			dec := gob.NewDecoder(buffer)

			if err := dec.Decode(&key_value_pairs); err != nil {
				log.Fatal(err)
			}

		default:
			{
				fmt.Println("unknown command")
			}

		}

	}
}
