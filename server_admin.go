package main

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"text/template"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

var (
	adress  string       = "https://schedule-cloud.cfuv.ru/index.php/s/YLoTDF3GqDjnDbR/download/09.03.01%20Информатика%20и%20вычислительная%20техника,09.03.04%20Программная%20инженерия%20%281-4%29.xlsx"
	ticker  *time.Ticker = time.NewTicker(30 * time.Minute)
	store                = sessions.NewCookieStore([]byte("0Z0llZ15JFVryM8OmXfh95dfBfqrpvif"))
	letters              = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
)

func randString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func download_latest_rasp() {
	resp, err := http.Get(adress)
	if err != nil {
		fmt.Println(time.Now().Format("02-01-06 15:04:05"), err)
		return
	}
	parse_xslx_rasp(resp, nil, nil)
	return
}

func upload_click(w http.ResponseWriter, r *http.Request) {

	if r.Method == "POST" {
		adress = r.FormValue("adress")
		minut, err := strconv.Atoi(r.FormValue("time"))
		if err == nil {
			ticker.Stop()
			ticker.Reset(time.Duration(minut) * time.Minute)
			logMessage(r, w, fmt.Sprintf("%s Новый адресс и период успешно передан : %s <br> Период в минутах : %d \n", time.Now().Format("02-01-06 15:04:05"), adress, minut))
		}
	}

	resp, err := http.Get(adress)
	if err != nil {
		logMessage(r, w, fmt.Sprintf("%s Ошибка получения файла из :<br> %s \n", time.Now().Format("02-01-06 15:04:05"), adress))
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	logMessage(r, w, fmt.Sprintf("%s Успешно получен файл из :<br> %s \n", time.Now().Format("02-01-06 15:04:05"), adress))
	parse_xslx_rasp(resp, r, w)
	http.Redirect(w, r, "/", http.StatusSeeOther)
	return
}

func uploadFile(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(10 << 20)

	file, err := r.MultipartForm.File["file"][0].Open()

	if err != nil {
		logMessage(r, w, fmt.Sprintf("%s Ошибка передачи файла \n", time.Now().Format("02-01-06 15:04:05")))
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	defer file.Close()

	resp := &http.Response{
		Body: file,
	}
	logMessage(r, w, fmt.Sprintf("%s Файл успешно получен \n", time.Now().Format("02-01-06 15:04:05")))
	parse_xslx_rasp(resp, r, w)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func main_page(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "cookie-user")
	data, _ := os.ReadFile("SITE.html")
	tmpl, _ := template.New("main").Parse(string(data))
	tmpl.Execute(w, session.Values["message"])
	return
}

func exit(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "cookie-user")
	session.Options.MaxAge = -1
	session.Save(r, w)
	w.Write([]byte("Сессия успешно закончена"))
	return
}

type User struct {
	Username string `json:"Username"`
	Tg_id    string `json:"Tg_id"`
	Id       string `json:"Id"`
	Role     string `json:"Role"`
	Group    string `json:"Group"`
}

type Data struct {
	Roles []User
	Error bool
}

func role_st(w http.ResponseWriter, r *http.Request) {
	data, _ := os.ReadFile("role.html")
	var dattest Data
	resp, err := http.Get("http://localhost:8080/get_all_users")
	if err != nil {
		dattest.Error = true
	} else {
		dattest.Error = false
		defer resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)
		json.Unmarshal(body, &dattest.Roles)
	}
	tmpl, _ := template.New("role").Parse(string(data))
	tmpl.Execute(w, dattest)
	return
}

func roleChangeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Write([]byte("Ошибка"))
		return
	}

	r.ParseForm()

	Id := r.FormValue("Id")
	role := r.FormValue("Role")

	resp, err := http.Get(fmt.Sprintf("http://localhost:8080/update?Id=%s&Key=Role&Value=%s", Id, role))
	if err != nil {
		logMessage(r, w, fmt.Sprintf("%s Ошибка выполнения запроса: %s", time.Now().Format("02-01-06 15:04:05"), err))
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	if resp.StatusCode == http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return
		}
		if string(body) == "Updated successfully" {
			logMessage(r, w, fmt.Sprintf("%s Успешно изменена роль пользователя %s на %s", time.Now().Format("02-01-06 15:04:05"), Id, role))
			http.Redirect(w, r, "/role", http.StatusSeeOther)
			return
		}
	} else {
		logMessage(r, w, fmt.Sprintf("%s Ошибка изменения ролей: %s", time.Now().Format("02-01-06 15:04:05"), err))
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
}

func delete_user(w http.ResponseWriter, r *http.Request) {

	tg := r.URL.Query().Get("tg_id")

	resp, err := http.Get(fmt.Sprintf("http://localhost:8080/del?tg_id=%s", tg))
	if err != nil {
		logMessage(r, w, fmt.Sprintf("%s *1 Ошибка удаления пользователя по tg_id : %s", time.Now().Format("02-01-06 15:04:05"), tg))
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	if resp.StatusCode != http.StatusOK {
		logMessage(r, w, fmt.Sprintf("%s *2 Ошибка удаления пользователя: %s", time.Now().Format("02-01-06 15:04:05"), tg))
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logMessage(r, w, fmt.Sprintf("%s *3 Ошибка удаления %s пользователя: %s", time.Now().Format("02-01-06 15:04:05"), tg, err))
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	if string(body) == "Deleted successfully" {
		logMessage(r, w, fmt.Sprintf("%s Пользователь по tg_id : %s успешно удален", time.Now().Format("02-01-06 15:04:05"), tg))
		http.Redirect(w, r, "/role", http.StatusSeeOther)
		return
	} else {
		logMessage(r, w, fmt.Sprintf("%s *4 Ошибка удаления пользователя по tg_id : %s", time.Now().Format("02-01-06 15:04:05"), tg))
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
}

var New_Sessions = make(map[string]time.Time)

var SECRET string = "123"

func new_session(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Пришел запрос новой сессии")
	sessi, _ := store.Get(r, "cookie-user")
	id := r.URL.Query().Get("id")
	if id == "" {
		id, _ = sessi.Values["token"].(string)
		fmt.Println(id)
	}
	expire, ok := New_Sessions[id]
	tokenString := r.URL.Query().Get("jwt_token")
	if tokenString != "" {
		fmt.Println("*")
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			fmt.Println("1")
			return []byte(SECRET), nil
		})
		if err != nil {
			fmt.Println("2")
			return
		}
		payload, okk := token.Claims.(jwt.MapClaims)
		if !(okk && token.Valid) {
			w.Write([]byte("Токен не валидный"))
			fmt.Println("3")
			return
		}
		if payload["expires_at"].(float64) < float64(time.Now().Unix()) {
			fmt.Println(payload["expires_at"].(float64), float64(time.Now().Unix()))
			w.Write([]byte("Токен просрочен"))
			fmt.Println("4")
			return
		}
		if payload["role"] != "admin" {
			w.Write([]byte("Недостаточно прав"))
			fmt.Println("5")
			return
		}
		if id == "" || len(id) != 4 {
			id = randString(4)
		}
		expire = time.Now().Add(time.Duration(30) * time.Second)
		New_Sessions[id] = expire
		w.Write([]byte(fmt.Sprintf("http://127.0.0.1:8082/new_session?id=%s", id)))
		fmt.Println(fmt.Sprintf("http://127.0.0.1:8082/new_session?id=%s", id))
		return
	} else if ok {
		fmt.Println("/")
		if expire.Before(time.Now()) {
			delete(New_Sessions, id)
			sessi.Options.MaxAge = -1
			sessi.Save(r, w)
			w.Write([]byte("Сессия просрочена"))
			fmt.Println("1")
			return
		}
		tokeExpiresAt := time.Now().Add(time.Minute * time.Duration(15))
		fmt.Println("2")
		payload_new := jwt.MapClaims{
			"Rights":     "Admin",
			"expires_at": tokeExpiresAt.Unix(),
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload_new)

		tokenStr, _ := token.SignedString([]byte(SECRET))

		sessi.Values["token"] = tokenStr
		sessi.Save(r, w)
		delete(New_Sessions, id)
		http.Redirect(w, r, "/", http.StatusSeeOther)
		fmt.Println("3")
		return

	}
	if r.URL.Query().Get("id") != "" {
		http.Redirect(w, r, "/new_session", http.StatusSeeOther)
		return
	}
	id = randString(4)
	sessi.Values["token"] = id
	fmt.Println(sessi.Values["token"])
	sessi.Save(r, w)
	fmt.Println("****")
	fmt.Fprintf(w, "Введите этот токен в телеграмм бота, а после перезагрузите страницу :\n \t %s", id)
	return
}

func Loggin_check(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/new_session" || r.URL.Path == "/Exit" {
			next.ServeHTTP(w, r)
			return
		} else {
			sessi, _ := store.Get(r, "cookie-user")
			tokenString, ok := sessi.Values["token"].(string)
			if ok {
				token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
					return []byte(SECRET), nil
				})

				if err != nil {
					http.Redirect(w, r, "/new_session", http.StatusSeeOther)
					return
				}

				payload, ok := token.Claims.(jwt.MapClaims)
				if !(ok && token.Valid) {
					sessi.Options.MaxAge = -1
					sessi.Save(r, w)
					w.Write([]byte("Токен не валидный"))
					return
				}
				if payload["expires_at"].(float64) < float64(time.Now().Unix()) {
					sessi.Options.MaxAge = -1
					sessi.Save(r, w)
					w.Write([]byte("Токен просрочен"))
					return
				}
				if payload["Rights"] != "Admin" {
					sessi.Options.MaxAge = -1
					sessi.Save(r, w)
					w.Write([]byte("Недостаточно прав"))
					return
				}

				tokeExpiresAt := time.Now().Add(time.Minute * time.Duration(15))
				payload_new := jwt.MapClaims{
					"Rights":     "Admin",
					"expires_at": tokeExpiresAt.Unix(),
				}

				token = jwt.NewWithClaims(jwt.SigningMethodHS256, payload_new)
				tokenStr, _ := token.SignedString([]byte(SECRET))

				sessi.Values["token"] = tokenStr
				sessi.Save(r, w)

				next.ServeHTTP(w, r)
				return
			}
			sessi.Options.MaxAge = -1
			sessi.Save(r, w)
			http.Redirect(w, r, "/new_session", http.StatusSeeOther)
			return
		}
	})
}

func main() {
	go func() {
		for {
			select {
			case <-ticker.C:
				fmt.Println("tick")
				download_latest_rasp()
			}
		}
	}()

	router := mux.NewRouter()
	router.HandleFunc("/", main_page)
	router.HandleFunc("/upload", uploadFile)
	router.HandleFunc("/role", role_st)
	router.HandleFunc("/Exit", exit)
	router.HandleFunc("/upload_latest", upload_click)
	router.HandleFunc("/new_session", new_session)
	router.HandleFunc("/role_change", roleChangeHandler)
	router.HandleFunc("/delete_user", delete_user)
	router.Use(Loggin_check)
	http.ListenAndServe(":8082", router)
}
