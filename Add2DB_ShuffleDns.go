package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var pause chan bool
var logger *log.Logger

func main() {
	setupLogger()
	pause = make(chan bool, 1)
	setupSignalHandling()
	go handlePauseResume()

	domainListPath := "lists.txt"

	for {
		runTasks(domainListPath)
		time.Sleep(4 * time.Hour)
	}
}

func setupLogger() {
	file, err := os.OpenFile(time.Now().Format("20060102_150405")+".log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Failed to open log file:", err)
	}
	logger = log.New(file, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
}

func runTasks(domainListPath string) {
	file, err := os.Open(domainListPath)
	if err != nil {
		logger.Println("Error opening domain list file:", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		domain := scanner.Text()
		if domain != "" {
			logger.Printf("Processing domain: %s\n", domain)
			processDomain(domain)
		}
	}

	if err := scanner.Err(); err != nil {
		logger.Println("Error reading domain list file:", err)
	}
}

func processDomain(domain string) {
	datetime := time.Now().Format("20060102_150405")
	outputFile := fmt.Sprintf("%s_%s.txt", domain, datetime)

	cmd := exec.Command("shuffledns", "-w", "static.txt", "-d", domain, "-r", "r.txt", "-m", "/bin/massdns", "-o", outputFile, "-nc", "-sw")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		logger.Println("Error executing shuffledns for domain", domain, ":", err)
		return
	}
	logger.Printf("shuffledns task completed for domain: %s\n", domain)

	processOutputFile(outputFile)
}

func processOutputFile(outputFile string) {
	db, err := sql.Open("sqlite3", "./domains.db")
	if err != nil {
		logger.Println("Error opening database:", err)
		return
	}
	defer db.Close()

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS domains (id INTEGER PRIMARY KEY, domain TEXT UNIQUE)")
	if err != nil {
		logger.Println("Error creating table:", err)
		return
	}

	var newDomains []string
	file, err := os.Open(outputFile)
	if err != nil {
		logger.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		domain := scanner.Text()
		if domain != "" {
			res, err := db.Exec("INSERT OR IGNORE INTO domains (domain) VALUES (?)", domain)
			if err != nil {
				logger.Printf("Error inserting domain %s: %s\n", domain, err)
			} else {
				id, _ := res.LastInsertId()
				if id != 0 {
					newDomains = append(newDomains, domain)
					logger.Printf("New domain %s added successfully.\n", domain)
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		logger.Println("Error reading file:", err)
	}

	if len(newDomains) > 0 {
		tmpFile, err := os.CreateTemp("", "newDomains_*.txt")
		if err != nil {
			logger.Println("Error creating temp file:", err)
			return
		}
		defer os.Remove(tmpFile.Name()) 

		for _, domain := range newDomains {
			if _, err := tmpFile.WriteString(domain + "\n"); err != nil {
				logger.Println("Error writing to temp file:", err)
				return
			}
		}

		if err := tmpFile.Close(); err != nil {
			logger.Println("Error closing temp file:", err)
			return
		}

		notifyCmd := exec.Command("notify", "-silent", "-data", tmpFile.Name(), "-bulk")
		notifyCmd.Stdout = os.Stdout
		notifyCmd.Stderr = os.Stderr
		err = notifyCmd.Run()
		if err != nil {
			logger.Println("Error sending notification:", err)
		} else {
			logger.Println("New domains have been sent to the Telegram channel.")
		}
	} else {
		logger.Println("No new domains to send to the Telegram channel.")
	}
}

func handlePauseResume() {
	reader := bufio.NewReader(os.Stdin)
	for {
		input, _ := reader.ReadByte()
		if input == 32 { // ASCII for space
			pause <- true
		}
	}
}

func setupSignalHandling() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		logger.Println("Received an interrupt, stopping services...")
		os.Exit(0)
	}()
}
