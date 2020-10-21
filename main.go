package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"mvault/credentials"
	"mvault/encryption"
	"mvault/filesys"
	"mvault/network"
	"os"
	"runtime"
	"runtime/pprof"
	"strings"
	"syscall"
	"time"

	"github.com/gosuri/uiprogress"
	"golang.org/x/crypto/ssh/terminal"
)

// Credentials is used to store user specific password, salt
type Credentials = credentials.Credentials

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

func getLocalSecrets() (Credentials, error) {
	var creds Credentials
	creds.Password = retrievesecret("Enter password: ")
	creds.Salt = retrievesecret("Enter pin: ")
	return creds, nil
}

func getSecrets(local bool) (Credentials, error) {
	if local {
		return getLocalSecrets()
	}
	return network.GetCreds("", "")
}

func randBetw(min int, max int) int {
	return rand.Intn(max-min) + min
}

func main() {
	rand.Seed(time.Now().UnixNano())
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

	// ADD HELP

	if len(*path) < 1 || (!*encr && !*decr) {
		fmt.Println("For usage information type `mvault -help`")
	}

	uiprogress.Start()
	bar := uiprogress.AddBar(100)

	data, readErr := filesys.ReadFile(*path)
	if readErr != nil {
		log.Fatal(readErr)
	}

	bar.Set(randBetw(7, 15))

	var creds Credentials

	creds, getErr := getSecrets(*local)
	if getErr != nil {
		log.Fatal(getErr)
	}

	bar.Set(randBetw(23, 40))

	var newdata []byte

	if *encr {
		ciphertext, encryptErr := encryption.Encrypt(data, creds)
		if encryptErr != nil {
			log.Fatal(encryptErr)
		}
		newdata = ciphertext
	}
	if *decr {
		plaintext, decryptErr := encryption.Decrypt(data, creds)
		if decryptErr != nil {
			log.Fatal(decryptErr)
		}
		newdata = plaintext
	}

	bar.Set(randBetw(75, 88))

	// ZIP ? :)

	replaceErr := filesys.Replace(*path, newdata)
	if replaceErr != nil {
		log.Fatal(replaceErr)
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

	bar.AppendCompleted()
}

/* IN CASE OF FAILURE ??? */
