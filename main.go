package main

import (
  // "fmt"
  "encoding/csv"
  "io"
  "log"
  "os"
  "sort"
  "time"
)

type Record struct {
  activation_date   time.Time
  deactivation_date time.Time
}

type Records []Record

func (l Records) Len() int {
  return len(l)
}

func (l Records) Less(i, j int) bool {
  return l[i].activation_date.Before(l[j].activation_date)
}

func (l Records) Swap(i, j int) {
  l[i], l[j] = l[j], l[i]
}

func parseTime(str string) time.Time {
  value, err := time.Parse("2006-01-02", str)
  if err != nil {
    log.Fatal("Error: %v, %v", str, err)
  }
  return value
}

func lastMonth(time time.Time) time.Time {
  return time.AddDate(0, -1, 0)
}

func dateDiff(from_date, to_date time.Time) bool {
  last_month := lastMonth(from_date)

  if to_date.Before(last_month) {
    return true
  } else {
    return false
  }
}

func exportCsv(filename string, l map[string]time.Time) {
  csvOut, err := os.Create(filename)
  if err != nil {
    log.Fatal("Unable to open output: %v", filename)
  }
  w := csv.NewWriter(csvOut)
  defer csvOut.Close()

  w.Write([]string{ "PHONE_NUMBER", "REAL_ACTIVATION_DATE" })
  w.Flush()

  for phone, time := range l {
    if err = w.Write([]string{ phone, time.Format("2006-01-02") }); err != nil {
      log.Fatal(err)
    }
    w.Flush()
  }
}

func main() {
  input := ""
  output := ""

  if len(os.Args) > 0 {
    input = os.Args[1]
    output = os.Args[2]
  }

  if len(input) == 0 {
    input = "./input.csv"
  }

  if len(output) == 0 {
    output = "./output.csv"
  }

  // A map of List of Records
  // key = phone
  // value = list of Record
  m_data := make(map[string]Records)

  // A map of phone and the activation date
  // key = phone
  // value = activation date
  m_result := make(map[string]time.Time)

  // Read from CSV file
  csvIn, err := os.Open(input)
  if err != nil {
    log.Fatal(err)
  }
  r := csv.NewReader(csvIn)

  // Read the 1st line (header) and ignore it
  rec, err := r.Read()
  if err != nil {
    log.Fatal(err)
  }

  // Loop row by row
  for {
    rec, err = r.Read()
    if err != nil {
      if err == io.EOF {
        break
      }
      log.Fatal(err)
    }

    phone := rec[0]
    activation_date := parseTime(rec[1])
    deactivation_date := parseTime(rec[2])

    m_data[phone] = append(m_data[phone], Record{ activation_date, deactivation_date })
  }

  for phone, records := range m_data {
    length := len(records)

    if length == 1 {
      m_result[phone] = records[0].activation_date
    } else {
      sort.Sort(records)
      m_result[phone] = records[0].activation_date

      j := length - 1
      for j >= 1 {
        if dateDiff(records[j].activation_date, records[j - 1].deactivation_date) {
          m_result[phone] = records[j].activation_date
          break
        }
        j--
      }
    }
  }

  exportCsv(output, m_result)
}
