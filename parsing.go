package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/xuri/excelize/v2"
)

type lession struct {
	Lec  string `json:"Название пары"`
	Form string `json:"Тип"`
	Prep string `json:"Преподаватель"`
	Room string `json:"Помещение"`
	Comm string `json:"Комментарий"`
}

type Group map[string]map[string]map[string][]lession

func NewLession() lession {
	return lession{
		Lec:  "",
		Form: "",
		Prep: "",
		Room: "",
		Comm: "",
	}
}

func NewGroup() Group {
	g :=
		make(map[string]map[string]map[string][]lession)
	for i := 231; i <= 233; i++ {
		for j := 1; j <= 2; j++ {
			// Преобразуем два целых числа в строку в виде "232(1)"
			key := fmt.Sprintf("%d(%d)", i, j)
			g[key] = make(map[string]map[string][]lession)
			g[key]["Нечётная неделя"] = make(map[string][]lession, 7)
			g[key]["Чётная неделя"] = make(map[string][]lession, 7)
			for _, k := range [6]string{"Понедельник", "Вторник", "Среда", "Четверг", "Пятница", "Суббота"} {
				g[key]["Нечётная неделя"][k] = make([]lession, 7)
				g[key]["Чётная неделя"][k] = make([]lession, 7)
				for ii := 0; ii < 7; ii++ {
					g[key]["Нечётная неделя"][k][ii] = NewLession()
					g[key]["Чётная неделя"][k][ii] = NewLession()
				}
			}
		}
	}
	return g
}

func parse_xslx_rasp(data *http.Response, r *http.Request, w http.ResponseWriter) {
	f, err := excelize.OpenReader(data.Body)
	if err != nil {
		logMessage(r, w, fmt.Sprintf("%s Ошибка чтения файла: %s \n", time.Now().Format("2006-01-02 15:04:05"), err.Error()))
		return
	}

	sheet := "курс 1 ПИ "

	rows, err := f.GetRows(sheet)
	if err != nil {
		logMessage(r, w, fmt.Sprintf("%s Ошибка получения строк: %s \n", time.Now().Format("2006-01-02 15:04:05"), err.Error()))
		return
	}

	sheetData := make([]map[string]string, 0)

	for _, row := range rows {
		rowData := make(map[string]string)
		for i, colCell := range row {
			rowData[fmt.Sprintf("col%d", i+1)] = colCell
		}
		sheetData = append(sheetData, rowData)
	}

	spisok := make(map[string]bool)
	for i := 4; i < len(sheetData); i += 5 {
		spisok[sheetData[i]["col5"]] = true
		spisok[sheetData[i]["col17"]] = true
	}
	raspisanie := NewGroup()
	var k, coun int
	var prepod, lecs, romm, count string

	day := make(map[int]string)
	day[1] = "Понедельник"
	day[2] = "Вторник"
	day[3] = "Среда"
	day[4] = "Четверг"
	day[5] = "Пятница"
	day[6] = "Суббота"
	coun = 0
	count = "Понедельник"

	for i := 4; i < len(sheetData); i += 5 {
		if sheetData[i]["col2"] != "0" {
			k, _ = strconv.Atoi(sheetData[i]["col2"])
			k--
		}
		if sheetData[i]["col2"] == "1" {
			coun++
			count = day[coun]
		}
		if sheetData[i]["col3"] == "ЛК" {
			lecs = sheetData[i]["col4"]
			prepod = sheetData[i+1]["col4"]
			romm = sheetData[i+2]["col4"]
			raspisanie["231(1)"]["Нечётная неделя"][count][k].Lec = lecs
			raspisanie["231(1)"]["Нечётная неделя"][count][k].Prep = prepod
			raspisanie["231(1)"]["Нечётная неделя"][count][k].Room = romm
			raspisanie["231(1)"]["Нечётная неделя"][count][k].Form = "Лекция"
			raspisanie["231(2)"]["Нечётная неделя"][count][k].Lec = lecs
			raspisanie["231(2)"]["Нечётная неделя"][count][k].Prep = prepod
			raspisanie["231(2)"]["Нечётная неделя"][count][k].Room = romm
			raspisanie["231(2)"]["Нечётная неделя"][count][k].Form = "Лекция"
			raspisanie["232(1)"]["Нечётная неделя"][count][k].Lec = lecs
			raspisanie["232(1)"]["Нечётная неделя"][count][k].Prep = prepod
			raspisanie["232(1)"]["Нечётная неделя"][count][k].Room = romm
			raspisanie["232(1)"]["Нечётная неделя"][count][k].Form = "Лекция"
			raspisanie["232(2)"]["Нечётная неделя"][count][k].Lec = lecs
			raspisanie["232(2)"]["Нечётная неделя"][count][k].Prep = prepod
			raspisanie["232(2)"]["Нечётная неделя"][count][k].Room = romm
			raspisanie["232(2)"]["Нечётная неделя"][count][k].Form = "Лекция"
			raspisanie["233(1)"]["Нечётная неделя"][count][k].Lec = lecs
			raspisanie["233(1)"]["Нечётная неделя"][count][k].Prep = prepod
			raspisanie["233(1)"]["Нечётная неделя"][count][k].Room = romm
			raspisanie["233(1)"]["Нечётная неделя"][count][k].Form = "Лекция"
			raspisanie["233(2)"]["Нечётная неделя"][count][k].Lec = lecs
			raspisanie["233(2)"]["Нечётная неделя"][count][k].Prep = prepod
			raspisanie["233(2)"]["Нечётная неделя"][count][k].Room = romm
			raspisanie["233(2)"]["Нечётная неделя"][count][k].Form = "Лекция"
		} else {
			if sheetData[i]["col3"] == "ПЗ" {
				if spisok[sheetData[i]["col4"]] {
					raspisanie["231(1)"]["Нечётная неделя"][count][k].Lec = sheetData[i]["col4"]
					raspisanie["231(1)"]["Нечётная неделя"][count][k].Prep = sheetData[i+1]["col4"]
					raspisanie["231(1)"]["Нечётная неделя"][count][k].Room = sheetData[i+2]["col4"]
					raspisanie["231(1)"]["Нечётная неделя"][count][k].Form = "Практичекое занятие"
					raspisanie["231(2)"]["Нечётная неделя"][count][k].Lec = sheetData[i]["col5"]
					raspisanie["231(2)"]["Нечётная неделя"][count][k].Prep = sheetData[i+1]["col5"]
					raspisanie["231(2)"]["Нечётная неделя"][count][k].Room = sheetData[i+2]["col5"]
					raspisanie["231(2)"]["Нечётная неделя"][count][k].Form = "Практичекое занятие"
				} else {
					lecs = sheetData[i]["col4"]
					prepod = sheetData[i+1]["col4"]
					romm = sheetData[i+2]["col4"]
					raspisanie["231(1)"]["Нечётная неделя"][count][k].Lec = lecs
					raspisanie["231(1)"]["Нечётная неделя"][count][k].Prep = prepod
					raspisanie["231(1)"]["Нечётная неделя"][count][k].Room = romm
					raspisanie["231(1)"]["Нечётная неделя"][count][k].Form = "Практичекое занятие"
					raspisanie["231(2)"]["Нечётная неделя"][count][k].Lec = lecs
					raspisanie["231(2)"]["Нечётная неделя"][count][k].Prep = prepod
					raspisanie["231(2)"]["Нечётная неделя"][count][k].Room = romm
					raspisanie["231(2)"]["Нечётная неделя"][count][k].Form = "Практичекое занятие"
				}
			}
			if sheetData[i]["col6"] == "ПЗ" {
				if spisok[sheetData[i]["col7"]] {
					raspisanie["232(1)"]["Нечётная неделя"][count][k].Lec = sheetData[i]["col7"]
					raspisanie["232(1)"]["Нечётная неделя"][count][k].Prep = sheetData[i+1]["col7"]
					raspisanie["232(1)"]["Нечётная неделя"][count][k].Room = sheetData[i+2]["col7"]
					raspisanie["232(1)"]["Нечётная неделя"][count][k].Form = "Практичекое занятие"
					raspisanie["232(2)"]["Нечётная неделя"][count][k].Lec = sheetData[i]["col8"]
					raspisanie["232(2)"]["Нечётная неделя"][count][k].Prep = sheetData[i+1]["col8"]
					raspisanie["232(2)"]["Нечётная неделя"][count][k].Room = sheetData[i+2]["col8"]
					raspisanie["232(2)"]["Нечётная неделя"][count][k].Form = "Практичекое занятие"
				} else {
					lecs = sheetData[i]["col7"]
					prepod = sheetData[i+1]["col7"]
					romm = sheetData[i+2]["col7"]
					raspisanie["232(1)"]["Нечётная неделя"][count][k].Lec = lecs
					raspisanie["232(1)"]["Нечётная неделя"][count][k].Prep = prepod
					raspisanie["232(1)"]["Нечётная неделя"][count][k].Room = romm
					raspisanie["232(1)"]["Нечётная неделя"][count][k].Form = "Практичекое занятие"
					raspisanie["232(2)"]["Нечётная неделя"][count][k].Lec = lecs
					raspisanie["232(2)"]["Нечётная неделя"][count][k].Prep = prepod
					raspisanie["232(2)"]["Нечётная неделя"][count][k].Room = romm
					raspisanie["232(2)"]["Нечётная неделя"][count][k].Form = "Практичекое занятие"
				}
			}
			if sheetData[i]["col9"] == "ПЗ" {
				if spisok[sheetData[i]["col10"]] {
					raspisanie["233(1)"]["Нечётная неделя"][count][k].Lec = sheetData[i]["col10"]
					raspisanie["233(1)"]["Нечётная неделя"][count][k].Prep = sheetData[i+1]["col10"]
					raspisanie["233(1)"]["Нечётная неделя"][count][k].Room = sheetData[i+2]["col10"]
					raspisanie["233(1)"]["Нечётная неделя"][count][k].Form = "Практичекое занятие"
					raspisanie["233(2)"]["Нечётная неделя"][count][k].Lec = sheetData[i]["col11"]
					raspisanie["233(2)"]["Нечётная неделя"][count][k].Prep = sheetData[i+1]["col11"]
					raspisanie["233(2)"]["Нечётная неделя"][count][k].Room = sheetData[i+2]["col11"]
					raspisanie["233(2)"]["Нечётная неделя"][count][k].Form = "Практичекое занятие"
				} else {
					lecs = sheetData[i]["col10"]
					prepod = sheetData[i+1]["col10"]
					romm = sheetData[i+2]["col10"]
					raspisanie["233(1)"]["Нечётная неделя"][count][k].Lec = lecs
					raspisanie["233(1)"]["Нечётная неделя"][count][k].Prep = prepod
					raspisanie["233(1)"]["Нечётная неделя"][count][k].Room = romm
					raspisanie["233(1)"]["Нечётная неделя"][count][k].Form = "Практичекое занятие"
					raspisanie["233(2)"]["Нечётная неделя"][count][k].Lec = lecs
					raspisanie["233(2)"]["Нечётная неделя"][count][k].Prep = prepod
					raspisanie["233(2)"]["Нечётная неделя"][count][k].Room = romm
					raspisanie["233(2)"]["Нечётная неделя"][count][k].Form = "Практичекое занятие"
				}
			}
		}
		if sheetData[i]["col15"] == "ЛК" {
			lecs = sheetData[i]["col16"]
			prepod = sheetData[i+1]["col16"]
			romm = sheetData[i+2]["col16"]
			raspisanie["231(1)"]["Чётная неделя"][count][k].Lec = lecs
			raspisanie["231(1)"]["Чётная неделя"][count][k].Prep = prepod
			raspisanie["231(1)"]["Чётная неделя"][count][k].Room = romm
			raspisanie["231(1)"]["Чётная неделя"][count][k].Form = "Лекция"
			raspisanie["231(2)"]["Чётная неделя"][count][k].Lec = lecs
			raspisanie["231(2)"]["Чётная неделя"][count][k].Prep = prepod
			raspisanie["231(2)"]["Чётная неделя"][count][k].Room = romm
			raspisanie["231(2)"]["Чётная неделя"][count][k].Form = "Лекция"
			raspisanie["232(1)"]["Чётная неделя"][count][k].Lec = lecs
			raspisanie["232(1)"]["Чётная неделя"][count][k].Prep = prepod
			raspisanie["232(1)"]["Чётная неделя"][count][k].Room = romm
			raspisanie["232(1)"]["Чётная неделя"][count][k].Form = "Лекция"
			raspisanie["232(2)"]["Чётная неделя"][count][k].Lec = lecs
			raspisanie["232(2)"]["Чётная неделя"][count][k].Prep = prepod
			raspisanie["232(2)"]["Чётная неделя"][count][k].Room = romm
			raspisanie["232(2)"]["Чётная неделя"][count][k].Form = "Лекция"
			raspisanie["233(1)"]["Чётная неделя"][count][k].Lec = lecs
			raspisanie["233(1)"]["Чётная неделя"][count][k].Prep = prepod
			raspisanie["233(1)"]["Чётная неделя"][count][k].Room = romm
			raspisanie["233(1)"]["Чётная неделя"][count][k].Form = "Лекция"
			raspisanie["233(2)"]["Чётная неделя"][count][k].Lec = lecs
			raspisanie["233(2)"]["Чётная неделя"][count][k].Prep = prepod
			raspisanie["233(2)"]["Чётная неделя"][count][k].Room = romm
			raspisanie["233(2)"]["Чётная неделя"][count][k].Form = "Лекция"
		} else {
			if sheetData[i]["col15"] == "ПЗ" {
				if spisok[sheetData[i]["col16"]] {
					raspisanie["231(1)"]["Чётная неделя"][count][k].Lec = sheetData[i]["col16"]
					raspisanie["231(1)"]["Чётная неделя"][count][k].Prep = sheetData[i+1]["col16"]
					raspisanie["231(1)"]["Чётная неделя"][count][k].Room = sheetData[i+2]["col16"]
					raspisanie["231(1)"]["Чётная неделя"][count][k].Form = "Практичекое занятие"
					raspisanie["231(2)"]["Чётная неделя"][count][k].Lec = sheetData[i]["col17"]
					raspisanie["231(2)"]["Чётная неделя"][count][k].Prep = sheetData[i+1]["col17"]
					raspisanie["231(2)"]["Чётная неделя"][count][k].Room = sheetData[i+2]["col17"]
					raspisanie["231(2)"]["Чётная неделя"][count][k].Form = "Практичекое занятие"
				} else {
					lecs = sheetData[i]["col16"]
					prepod = sheetData[i+1]["col16"]
					romm = sheetData[i+2]["col16"]
					raspisanie["231(1)"]["Чётная неделя"][count][k].Lec = lecs
					raspisanie["231(1)"]["Чётная неделя"][count][k].Prep = prepod
					raspisanie["231(1)"]["Чётная неделя"][count][k].Room = romm
					raspisanie["231(1)"]["Чётная неделя"][count][k].Form = "Практичекое занятие"
					raspisanie["231(2)"]["Чётная неделя"][count][k].Lec = lecs
					raspisanie["231(2)"]["Чётная неделя"][count][k].Prep = prepod
					raspisanie["231(2)"]["Чётная неделя"][count][k].Room = romm
					raspisanie["231(2)"]["Чётная неделя"][count][k].Form = "Практичекое занятие"
				}
			}
			if sheetData[i]["col18"] == "ПЗ" {
				if spisok[sheetData[i]["col19"]] {
					raspisanie["232(1)"]["Чётная неделя"][count][k].Lec = sheetData[i]["col19"]
					raspisanie["232(1)"]["Чётная неделя"][count][k].Prep = sheetData[i+1]["col19"]
					raspisanie["232(1)"]["Чётная неделя"][count][k].Room = sheetData[i+2]["col19"]
					raspisanie["232(1)"]["Чётная неделя"][count][k].Form = "Практичекое занятие"
					raspisanie["232(2)"]["Чётная неделя"][count][k].Lec = sheetData[i]["col20"]
					raspisanie["232(2)"]["Чётная неделя"][count][k].Prep = sheetData[i+1]["col20"]
					raspisanie["232(2)"]["Чётная неделя"][count][k].Room = sheetData[i+2]["col20"]
					raspisanie["232(2)"]["Чётная неделя"][count][k].Form = "Практичекое занятие"
				} else {
					lecs = sheetData[i]["col19"]
					prepod = sheetData[i+1]["col19"]
					romm = sheetData[i+2]["col19"]
					raspisanie["232(1)"]["Чётная неделя"][count][k].Lec = lecs
					raspisanie["232(1)"]["Чётная неделя"][count][k].Prep = prepod
					raspisanie["232(1)"]["Чётная неделя"][count][k].Room = romm
					raspisanie["232(1)"]["Чётная неделя"][count][k].Form = "Практичекое занятие"
					raspisanie["232(2)"]["Чётная неделя"][count][k].Lec = lecs
					raspisanie["232(2)"]["Чётная неделя"][count][k].Prep = prepod
					raspisanie["232(2)"]["Чётная неделя"][count][k].Room = romm
					raspisanie["232(2)"]["Чётная неделя"][count][k].Form = "Практичекое занятие"
				}
			}
			if sheetData[i]["col21"] == "ПЗ" {
				if spisok[sheetData[i]["col22"]] {
					raspisanie["233(1)"]["Чётная неделя"][count][k].Lec = sheetData[i]["col22"]
					raspisanie["233(1)"]["Чётная неделя"][count][k].Prep = sheetData[i+1]["col22"]
					raspisanie["233(1)"]["Чётная неделя"][count][k].Room = sheetData[i+2]["col22"]
					raspisanie["233(1)"]["Чётная неделя"][count][k].Form = "Практичекое занятие"
					raspisanie["233(2)"]["Чётная неделя"][count][k].Lec = sheetData[i]["col23"]
					raspisanie["233(2)"]["Чётная неделя"][count][k].Prep = sheetData[i+1]["col23"]
					raspisanie["233(2)"]["Чётная неделя"][count][k].Room = sheetData[i+2]["col23"]
					raspisanie["233(2)"]["Чётная неделя"][count][k].Form = "Практичекое занятие"
				} else {
					lecs = sheetData[i]["col22"]
					prepod = sheetData[i+1]["col22"]
					romm = sheetData[i+2]["col22"]
					raspisanie["233(1)"]["Чётная неделя"][count][k].Lec = lecs
					raspisanie["233(1)"]["Чётная неделя"][count][k].Prep = prepod
					raspisanie["233(1)"]["Чётная неделя"][count][k].Room = romm
					raspisanie["233(1)"]["Чётная неделя"][count][k].Form = "Практичекое занятие"
					raspisanie["233(2)"]["Чётная неделя"][count][k].Lec = lecs
					raspisanie["233(2)"]["Чётная неделя"][count][k].Prep = prepod
					raspisanie["233(2)"]["Чётная неделя"][count][k].Room = romm
					raspisanie["233(2)"]["Чётная неделя"][count][k].Form = "Практичекое занятие"
				}
			}
		}
	}

	jsonData, err := json.MarshalIndent(raspisanie, " ", "  ")
	if err != nil {
		logMessage(r, w, fmt.Sprintf("%s Ошибка: %s", time.Now().Format("2006-01-02 15:04:05"), err))
		return
	}
	resp, err := http.Post("https://599c-139-28-177-252.ngrok-free.app/UpdateSchedule", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		logMessage(r, w, fmt.Sprintf("%s Ошибка отправки расписания: %s", time.Now().Format("2006-01-02 15:04:05"), err))
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		logMessage(r, w, fmt.Sprintf("%s Расписание успешно передано", time.Now().Format("2006-01-02 15:04:05")))
	} else {
		logMessage(r, w, fmt.Sprintf("%s Ошибка при передаче на расписания : %s", time.Now().Format("2006-01-02 15:04:05"), resp.Status))
	}
}

func logMessage(r *http.Request, w http.ResponseWriter, message string) {
	if r != nil {
		session, _ := store.Get(r, "cookie-user")
		msg, ok := session.Values["message"].([]string)
		if !ok {
			session.Values["message"] = make([]string, 1)
		}
		session.Values["message"] = append(msg, message)
		session.Save(r, w)
	} else {
		fmt.Println(message)
	}
}
