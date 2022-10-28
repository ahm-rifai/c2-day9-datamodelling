package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"personal-web/connection"
	"strconv"
	"text/template"
	"time"

	"github.com/gorilla/mux"
)

func main() {

	route := mux.NewRouter()

	connection.DatabaseConnect()
	

	route.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("./public"))))

	route.HandleFunc("/", home).Methods("GET")
	route.HandleFunc("/contact", contact).Methods("GET")

	route.HandleFunc("/project-detail/{index}", projectDetail).Methods("GET")
	
	route.HandleFunc("/add-project", addProject).Methods("GET")	
	route.HandleFunc("/create-project", createProject).Methods("POST") // CREATE PROJECT
	
	route.HandleFunc("/edit-project/{index}", editProject).Methods("GET") // EDIT PROJECT

	route.HandleFunc("/delete-project/{index}", deleteProject).Methods("GET") // DELETE PROJECT

	fmt.Println("Server running on port 3000")
	http.ListenAndServe("localhost:3000", route)

}


type Project struct {
	ID 				int
	ProjectName 	string
	StartDate 		string
	EndDate 		string
	Duration 		string
	Description 	string
	Technologies 	[]string
	Image 			string
}

var dataProject = []Project {}



func home(w http.ResponseWriter, r *http.Request)  {
	w.Header().Set("Content-Type", "Text/html; charset=utp-8")
	var tmpl, err = template.ParseFiles("views/index.html")

	if err != nil {
		w.Write([]byte("Message : " + err.Error()))
		return
	}

	data, _ := connection.Conn.Query(context.Background(), `SELECT "ID", "ProjectName", "StartDate", "EndDate", "Duration", "Description", "Technologies", "Image" FROM tb_projects`)
	var result []Project
	for data.Next() {

		var item = Project{}

		err := data.Scan(&item.ID, &item.ProjectName, &item.StartDate, &item.EndDate, &item.Duration, &item.Description, &item.Technologies, &item.Image)

		if err != nil {
			fmt.Println(err.Error())
			return
		}

		result = append(result, item)
	}

	response := map[string]interface{}{
		"DataProject": result,
	}

	

	tmpl.Execute(w, response)
}

func contact(w http.ResponseWriter, r *http.Request)  {
	w.Header().Set("Content-Type", "Text/html; charset=utp-8")
	var tmpl, err = template.ParseFiles("views/contact.html")

	if err != nil {
		w.Write([]byte("Message : " + err.Error()))
		return
	}

	tmpl.Execute(w, nil)
}

func projectDetail(w http.ResponseWriter, r *http.Request)  {
	w.Header().Set("Content-Type", "Text/html; charset=utp-8")
	var tmpl, err = template.ParseFiles("views/project-detail.html")

	if err != nil {
		w.Write([]byte("Message : " + err.Error()))
		return
	}

	var RenderProjectDetail = Project{}

	index, _ := strconv.Atoi(mux.Vars(r)["index"])
	
	for i, data := range dataProject {
		if i == index {
			RenderProjectDetail = Project{
				ProjectName: data.ProjectName,
				StartDate: data.StartDate,
				EndDate: data.EndDate,
				Duration: data.Duration,
				Description: data.Description,
				Technologies: data.Technologies,
			}
		}
	}

	data := map[string]interface{}{
		"RenderProjectDetail": RenderProjectDetail,
	}

	tmpl.Execute(w, data)
}

func addProject(w http.ResponseWriter, r *http.Request)  {
	w.Header().Set("Content-Type", "Text/html; charset=utp-8")
	var tmpl, err = template.ParseFiles("views/add-project.html")

	if err != nil {
		w.Write([]byte("Message : " + err.Error()))
		return
	}

	tmpl.Execute(w, nil)
}

func createProject(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}

	projectName := r.PostForm.Get("input-project")
	startDate  := r.PostForm.Get("input-start")
	endDate  := r.PostForm.Get("input-end")
	description := r.PostForm.Get("input-desc")
	iconNodeJS := r.PostForm.Get("node")
	iconReactJS := r.PostForm.Get("react")
	iconNextJS := r.PostForm.Get("next")
	iconTypescript := r.PostForm.Get("type")
	

	
	var newProject = Project {
		ProjectName: projectName,
		StartDate: formatDate(startDate),
		EndDate: formatDate(endDate),
		Duration: getDuration(startDate, endDate),
		Description: description,
		Technologies: []string{iconNodeJS, iconReactJS, iconNextJS, iconTypescript},
	}

	dataProject = append(dataProject, newProject)

	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

func editProject(w http.ResponseWriter, r *http.Request)  {
	w.Header().Set("Content-Type", "Text/html; charset=utp-8")
	var tmpl, err = template.ParseFiles("views/edit-project.html")

	if err != nil {
		w.Write([]byte("Message : " + err.Error()))
		return
	}

	var updateProject = Project{}
	index, _ := strconv.Atoi(mux.Vars(r)["index"])

	for i, data := range dataProject {
		if index == i {
			updateProject = Project {
				ProjectName: data.ProjectName,
				StartDate: returnDate(data.StartDate),
				EndDate: returnDate(data.EndDate),
				Description: data.Description,
				Technologies: data.Technologies,
			}
			dataProject = append(dataProject[:index], dataProject[index+1:]...)
		}
	}

	data := map[string]interface{} {
		"updateProject" : updateProject,
	}

	tmpl.Execute(w, data)
}

func deleteProject(w http.ResponseWriter, r *http.Request)  {
	index, _ := strconv.Atoi(mux.Vars(r)["index"])

	dataProject = append(dataProject[:index], dataProject[index+1:]...)

	http.Redirect(w, r, "/", http.StatusFound)
}










func getDuration(startDate string, endDate string) string {

	layout := "2006-01-02"

	projectPost, _ := time.Parse(layout, startDate)
	currenTime, _ := time.Parse(layout, endDate)

	distance := currenTime.Sub(projectPost).Hours() / 24
	var duration string

	if distance > 30 {
		if (distance / 30) <= 1 {
			duration = "1 Month"
		} else {
			duration = strconv.Itoa(int(distance)/30) + " Months"
		}
	} else {
		if distance <= 1 {
			duration = "1 Day"
		} else {
			duration = strconv.Itoa(int(distance)) + " Days"
		}
	}

	return duration
}

func formatDate(InputDate string) string {

	layout := "2006-01-02"
	t, _ := time.Parse(layout, InputDate)

	formated := t.Format("02 January 2006")

	return formated
}

func returnDate(InputDate string) string {

	layout := "02 January 2006"
	t, _ := time.Parse(layout, InputDate)

	formated := t.Format("2006-01-02")

	return formated
}
