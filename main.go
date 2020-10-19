package main

import (
	"flag"
	"fmt"
	"log"
	"mvault/encryption"
	"mvault/filesys"
	"mvault/network"
	"os"
	"path/filepath"
	"runtime"
	"runtime/pprof"
	"strings"
	"syscall"

	"golang.org/x/crypto/ssh/terminal"
)

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to `file`")
var memprofile = flag.String("memprofile", "", "write memory profile to `file`")
var path = flag.String("file", "", "pass the `file`")
var local = flag.Bool("local", false, "do not interact with server")
var encr = flag.Bool("encrypt", false, "encrypt file")
var decr = flag.Bool("decrypt", false, "decrypt file")

func retrievesecret(intro string) string {
	fmt.Print(intro)
	bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
	if err == nil {
		fmt.Println("\nPassword typed: " + string(bytePassword))
	}
	password := string(bytePassword)
	return strings.TrimSpace(password)
}

func getLocalSecrets() (string, string, error) {
	pwd := retrievesecret("Enter password: ")
	log.Println(pwd)
	pin := retrievesecret("Enter pin: ")
	log.Println(pin)
	return pwd, pin, nil
}

func getSecrets(local bool) (string, string, error) {
	if local {
		return getLocalSecrets()
	}
	return network.GetCreds("", "")
}

func main() {
	flag.BoolVar(local, "l", *local, "alias for `-local`")
	flag.BoolVar(encr, "e", *encr, "alias for `-encrypt`")
	flag.BoolVar(decr, "d", *decr, "alias for `-decrypt`")
	flag.Parse()

	// === CPUPROFILE START ===
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal("could not create CPU profile: ", err)
		}
		defer f.Close() // error handling omitted for example
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal("could not start CPU profile: ", err)
		}
		defer pprof.StopCPUProfile()
	}
	// === CPUPROFILE END ===

	data, err := filesys.ReadFile(*path)
	if err != nil {
		log.Fatal(err)
	}

	passwd, pin, err := getSecrets(*local) // fix bullshit with string,string,err
	if err != nil {
		log.Fatal(err)
	}

	var newdata string

	if *encr {
		ciphertext, nonce := encryption.Encrypt(data, passwd, pin)
		newdata = string(append(nonce, ciphertext...))
	}
	if *decr {
		nonce, ciphertext := data[:12], data[12:]
		plaintext, err := encryption.Decrypt(ciphertext, nonce, passwd, pin)
		if err != nil {
			log.Fatal(err)
		}
		newdata = plaintext
	}

	result := []byte(newdata)

	// ZIP ? :)

	newpath := filepath.Dir(*path) + "_" + filepath.Base(*path)
	err = filesys.WriteFile(result, newpath)
	if err != nil {
		log.Fatal(err)
	}

	err = os.Remove(*path)
	if err != nil {
		log.Fatal(err)
	}

	err = os.Rename(newpath, *path)
	if err != nil {
		log.Fatal(err)
	}

	// === MEMPROFILE START ===
	if *memprofile != "" {
		f, err := os.Create(*memprofile)
		if err != nil {
			log.Fatal("could not create memory profile: ", err)
		}
		defer f.Close() // error handling omitted for example
		runtime.GC()    // get up-to-date statistics
		if err := pprof.WriteHeapProfile(f); err != nil {
			log.Fatal("could not write memory profile: ", err)
		}
	}
	// === MEMPROFILE END ===
}
