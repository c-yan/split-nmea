package main

import (
	"bufio"
	"log"
	"os"
	"strings"
	"time"
)

func main() {
	var (
		err    error
		file   *os.File
		writer *bufio.Writer
	)

	log.SetFlags(log.Lshortfile)

	f := os.Stdin
	if len(os.Args) == 2 {
		f, err = os.Open(os.Args[1])
		if err != nil {
			log.Fatalf("ファイルがオープンできませんでいた。\n")
		}
	}
	defer f.Close()

	pt := time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)
	reader := bufio.NewReader(f)
	lineNo := 0
	for {
		gpggaLine, err := reader.ReadString('\n')
		if strings.Trim(gpggaLine, "\r\n") == "" {
			break
		}
		if !strings.HasPrefix(gpggaLine, "$GPGGA,") {
			log.Fatalf("データが想定と違います。lineNo: %d\n", lineNo)
		}
		lineNo++
		gprmcLine, err := reader.ReadString('\n')
		if !strings.HasPrefix(gprmcLine, "$GPRMC,") {
			log.Fatalf("データが想定と違います。lineNo: %d\n", lineNo)
		}
		lineNo++

		v := strings.Split(gprmcLine, ",")
		if len(v) != 13 {
			log.Fatalf("$GPRMC の要素数が13ではありません. %v\n", gprmcLine)
		}

		t, err := time.Parse("020106150405", (v[9] + v[1])[:12])
		if err != nil {
			log.Fatalf("$GPRMC の日時がおかしいです。\n")
		}
		d := t.Sub(pt)
		if d.Minutes() > 20 || d.Minutes() < 0 {
			if file != nil {
				writer.Flush()
				file.Close()
			}
			file, err = os.OpenFile(t.Format("20060102-1504")+".nma", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
			if err != nil {
				log.Fatalf("ファイルのオープンに失敗しました。\n")
			}
			writer = bufio.NewWriter(file)
		}
		pt = t

		writer.WriteString(gpggaLine)
		writer.WriteString(gprmcLine)
	}
	if file != nil {
		writer.Flush()
		file.Close()
	}
}
