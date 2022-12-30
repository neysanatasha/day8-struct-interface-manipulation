package main

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/gosimple/slug"
)

func main() {
	route := mux.NewRouter()

	// for public folder
	// ex: localhost:port/public/ +../path/to/file
	route.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("./public"))))

	route.HandleFunc("/", home).Methods("GET")
	// CRUD Project
	route.HandleFunc("/add-project", addProject).Methods("GET")
	route.HandleFunc("/store-project", storeProject).Methods("POST")
	route.HandleFunc("/detail-project/{id}", detailProject).Methods("GET")
	route.HandleFunc("/edit-project/{id}", editProject).Methods("GET")
	route.HandleFunc("/update-project/{id}", updateProject).Methods("POST")
	route.HandleFunc("/delete-project/{id}", deleteProject).Methods("GET")
	route.HandleFunc("/contact", contact).Methods("GET")

	fmt.Println("Server berjalan pada port 8000")
	http.ListenAndServe("localhost:8000", route)
}

// home
func home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html; charset=utf-8")
	tmpt, err := template.ParseFiles("views/index.html")

	if err != nil {
		w.Write([]byte("Message: " + err.Error()))
		return
	}

	listProject := map[string]interface{}{
		"Projects": projects,
	}

	tmpt.Execute(w, listProject)
}

//
// CRUD Project
//

// Project Struct
type Project struct {
	ID           int
	ProjectName  string
	StartDate    string
	EndDate      string
	Duration     string
	Description  string
	Technologies []string
	Image        string
}

var projects = []Project{
	{
		ProjectName:  "Project Title",
		StartDate:    "2022-11-03",
		EndDate:      "2022-11-18",
		Duration:     "3 month",
		Description:  "Pulu pulu pulu pulu pulu pulu pulu, pulu pulu pulu pulu pulu pulu pulu! Pulu pulu pulu pulu pulu pulu pulu, pulu pulu pulu pulu pulu pulu pulu. ",
		Technologies: []string{"nodejs", "vuejs", "reactjs", "nextjs"},
		Image:        "public/img/iconlogo.jpg",
	},
}

// createProject
func addProject(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html; charset=utf-8")
	tmpt, err := template.ParseFiles("views/add-project.html")

	if err != nil {
		w.Write([]byte("Message: " + err.Error()))
		return
	}

	tmpt.Execute(w, nil)
}

// storeProject
func storeProject(w http.ResponseWriter, r *http.Request) {
	// left shift 32 << 20 which results in 32*2^20 = 33554432
	// x << y, results in x*2^y
	err := r.ParseMultipartForm(32 << 20)

	if err != nil {
		log.Fatal(err)
	}

	project_name := r.PostForm.Get("project_name")
	technologies := r.Form["technologies"]
	description := r.PostForm.Get("description")

	// Image
	// Retrieve the image from form data
	uploadedFile, handler, err := r.FormFile("image")
	if err != nil {
		w.Write([]byte("Message: " + err.Error()))
		return
	}
	defer uploadedFile.Close()
	fileLocation := "public/uploads/"
	imageName := slug.Make(project_name)
	_ = os.MkdirAll(fileLocation, os.ModePerm)
	fullPath := fileLocation + imageName + filepath.Ext(handler.Filename)
	targetFile, err := os.OpenFile(fullPath, os.O_WRONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		w.Write([]byte("Message: " + err.Error()))
		return
	}
	defer targetFile.Close()
	// Copy the file to the destination path
	_, err = io.Copy(targetFile, uploadedFile)
	if err != nil {
		w.Write([]byte("Message: " + err.Error()))
		return
	}
	// Image Path
	imagePath := fileLocation + imageName + filepath.Ext(handler.Filename)
	// End For Image

	// Duration
	// Date
	startDate := r.PostForm.Get("start_date")
	endDate := r.PostForm.Get("end_date")
	const (
		layoutISO = "2006-01-02"
	)
	tStartDate, _ := time.Parse(layoutISO, startDate)
	tEndDate, _ := time.Parse(layoutISO, endDate)
	diff := tEndDate.Sub(tStartDate)

	months := int64(diff.Hours() / 24 / 30)
	days := int64(diff.Hours() / 24)

	if days%30 >= 0 {
		days = days % 30
	}

	var duration string

	if months >= 1 && days >= 1 {
		duration = strconv.FormatInt(months, 10) + " month " + strconv.FormatInt(days, 10) + " days"
	} else if months >= 1 && days <= 0 {
		duration = strconv.FormatInt(months, 10) + " month"
	} else if months < 1 && days >= 0 {
		duration = strconv.FormatInt(days, 10) + " days"
	} else {
		duration = "0 days"
	}
	// End for Duration

	var newProject = Project{
		ProjectName:  project_name,
		StartDate:    startDate,
		EndDate:      endDate,
		Duration:     duration,
		Description:  description,
		Technologies: technologies,
		Image:        imagePath,
	}

	projects = append(projects, newProject)

	// fmt.Println(projects)

	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

// detailProject
func detailProject(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html; charset=utf-8")
	tmpt, err := template.ParseFiles("views/detail-project.html")

	if err != nil {
		w.Write([]byte("Message: " + err.Error()))
		return
	}
	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	var DataProject = Project{}

	for index, data := range projects {
		if index == id {
			DataProject = Project{
				ID:           id,
				ProjectName:  data.ProjectName,
				StartDate:    data.StartDate,
				EndDate:      data.EndDate,
				Duration:     data.Duration,
				Description:  data.Description,
				Technologies: data.Technologies,
				Image:        data.Image,
			}
		}
	}

	EditProject := map[string]interface{}{
		"Project": DataProject,
	}
	// fmt.Println(EditProject)
	tmpt.Execute(w, EditProject)
}

// editProject
func editProject(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html; charset=utf-8")
	tmpt, err := template.ParseFiles("views/edit-project.html")

	if err != nil {
		w.Write([]byte("Message: " + err.Error()))
		return
	}
	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	var DataProject = Project{}

	for index, data := range projects {
		if index == id {
			DataProject = Project{
				ID:           id,
				ProjectName:  data.ProjectName,
				StartDate:    data.StartDate,
				EndDate:      data.EndDate,
				Description:  data.Description,
				Technologies: data.Technologies,
				Image:        data.Image,
			}
		}
	}

	EditProject := map[string]interface{}{
		"Project": DataProject,
	}
	// fmt.Println(EditProject)
	tmpt.Execute(w, EditProject)
}

// updateProject
func updateProject(w http.ResponseWriter, r *http.Request) {

	// left shift 32 << 20 which results in 32*2^20 = 33554432
	// x << y, results in x*2^y
	err := r.ParseMultipartForm(32 << 20)

	if err != nil {
		log.Fatal(err)
	}

	project_name := r.PostForm.Get("project_name")
	technologies := r.Form["technologies"]
	description := r.PostForm.Get("description")

	// Image
	// Retrieve the image from form data
	uploadedFile, handler, err := r.FormFile("image")
	if err != nil {
		w.Write([]byte("Error message upload file: " + err.Error()))
		return
	}
	defer uploadedFile.Close()
	fileLocation := "public/uploads/"
	imageName := slug.Make(project_name)
	_ = os.MkdirAll(fileLocation, os.ModePerm)
	fullPath := fileLocation + imageName + filepath.Ext(handler.Filename)
	targetFile, err := os.OpenFile(fullPath, os.O_WRONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		w.Write([]byte("Error message target file: " + err.Error()))
		return
	}
	defer targetFile.Close()
	// Copy the file to the destination path
	_, err = io.Copy(targetFile, uploadedFile)
	if err != nil {
		w.Write([]byte("Error message copy file: " + err.Error()))
		return
	}
	// Image Path
	imagePath := fileLocation + imageName + filepath.Ext(handler.Filename)
	// End For Image

	// Duration
	// Date
	startDate := r.PostForm.Get("start_date")
	endDate := r.PostForm.Get("end_date")
	const (
		layoutISO = "2006-01-02"
	)
	tStartDate, _ := time.Parse(layoutISO, startDate)
	tEndDate, _ := time.Parse(layoutISO, endDate)
	diff := tEndDate.Sub(tStartDate)

	months := int64(diff.Hours() / 24 / 30)
	days := int64(diff.Hours() / 24)

	if days%30 >= 0 {
		days = days % 30
	}

	var duration string

	if months >= 1 && days >= 1 {
		duration = strconv.FormatInt(months, 10) + " month " + strconv.FormatInt(days, 10) + " days"
	} else if months >= 1 && days <= 0 {
		duration = strconv.FormatInt(months, 10) + " month"
	} else if months < 1 && days >= 0 {
		duration = strconv.FormatInt(days, 10) + " days"
	} else {
		duration = "0 days"
	}
	// End for Duration

	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	projects[id].ProjectName = project_name
	projects[id].ProjectName = project_name
	projects[id].StartDate = startDate
	projects[id].EndDate = endDate
	projects[id].Duration = duration
	projects[id].Description = description
	projects[id].Technologies = technologies
	projects[id].Image = imagePath

	// fmt.Println(projects)

	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

// deleteProject
func deleteProject(w http.ResponseWriter, r *http.Request) {

	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	projects = append(projects[:id], projects[id+1:]...)

	http.Redirect(w, r, "/", http.StatusFound)
}

// contact
func contact(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html; charset=utf-8")
	tmpt, err := template.ParseFiles("views/contact.html")

	if err != nil {
		w.Write([]byte("Message: " + err.Error()))
		return
	}

	tmpt.Execute(w, nil)
}
