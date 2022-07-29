package handler

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"time"
)

type Handler struct {
	DB   *sql.DB
	Tmpl *template.Template
}

type User struct {
	Id       int
	Username string
	Password string
}

type Item struct {
	Title string
	Link  string
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	err := h.Tmpl.ExecuteTemplate(w, "login.html", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) Auth(w http.ResponseWriter, r *http.Request) {
	var body string
	inLogin := r.FormValue("login")
	inPassword := r.FormValue("password")
	row := h.DB.QueryRow("SELECT id, login FROM users WHERE login = ? and password = ? LIMIT 1", inLogin, inPassword)
	err := row.Scan()
	if err == sql.ErrNoRows {
		body += fmt.Sprintln("Placeholders case: NOT FOUND")
	} else {
		expiration := time.Now().Add(10 * time.Hour)
		cookie := http.Cookie{
			Name:    "session_id",
			Value:   inLogin,
			Expires: expiration,
		}
		http.SetCookie(w, &cookie)
		http.Redirect(w, r, "/", http.StatusFound)
		body += fmt.Sprintln("Placeholders id:", inLogin, "login:", inPassword)
	}
	w.Write([]byte(body))
}
func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	session, err := r.Cookie("session_id")
	if err == http.ErrNoCookie {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	session.Expires = time.Now().AddDate(0, 0, -1)
	http.SetCookie(w, session)
	http.Redirect(w, r, "/", http.StatusFound)
}

func (h *Handler) AdminPage(w http.ResponseWriter, r *http.Request) {
	users := []*User{}

	rows, err := h.DB.Query("SELECT id, login, password updated FROM users")
	if err != nil {
		panic(err)
	}
	for rows.Next() {
		user := &User{}
		err = rows.Scan(&user.Id, &user.Username, &user.Password)
		if err != nil {
			panic(err)
		}
		users = append(users, user)
	}
	rows.Close()

	err = h.Tmpl.ExecuteTemplate(w, "admin.html", struct {
		Users []*User
	}{
		Users: users,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) AddAdmin(w http.ResponseWriter, r *http.Request) {
	result, err := h.DB.Exec(
		"INSERT INTO users (`login`, `password`) VALUES (?, ?)",
		r.FormValue("username"),
		r.FormValue("password"),
	)
	if err != nil {
		panic(err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		panic(err)
	}
	lastID, err := result.LastInsertId()
	if err != nil {
		panic(err)
	}

	fmt.Println("Insert - RowsAffected", affected, "LastInsertId: ", lastID)

	http.Redirect(w, r, "/admin", http.StatusFound)

}

func (h *Handler) Head(w http.ResponseWriter, r *http.Request) {
	session, err := r.Cookie("session_id")

	var (
		id    int
		login string
	)
	row := h.DB.QueryRow("SELECT id, login FROM users WHERE login = ?  LIMIT 1", session.Value)
	row.Scan(&id, &login)
	items := []*Item{}

	rows, err := h.DB.Query("SELECT title, link FROM items")
	if err != nil {
		panic(err)
	}
	for rows.Next() {
		post := &Item{}
		err = rows.Scan(&post.Title, &post.Link)
		if err != nil {
			panic(err)
		}
		items = append(items, post)
	}
	rows.Close()

	err = h.Tmpl.ExecuteTemplate(w, "index.html", struct {
		Items []*Item
		Login string
	}{
		Items: items,
		Login: login,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}
