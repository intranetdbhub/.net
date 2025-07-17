package main

import (
	"database/sql"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"

	_ "github.com/lib/pq"
)

var db *sql.DB
var templates = template.Must(template.ParseGlob("templates/*.html"))

type Device struct {
	ID           int
	Class        string
	Colocation   string
	Description  string
	Hostname     string
	MgmtIP       string
	SerialNumber string
	DeviceType   string
}

func main() {
	var err error
	connStr := "host=localhost port=5432 user=postgres password=cisco dbname=testdb sslmode=disable"
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err = db.Ping(); err != nil {
		log.Fatal("Could not connect to PostgreSQL:", err)
	}

	// Routes
	http.HandleFunc("/", formHandler)
	http.HandleFunc("/devices", listDevices)
	http.HandleFunc("/edit", editDevice)
	http.HandleFunc("/update", updateDevice)
	http.HandleFunc("/delete", deleteDevice)
	http.HandleFunc("/export/csv", exportCSV)
	http.HandleFunc("/export/json", exportJSON)

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	fmt.Println("Device Inventory App running at http://localhost:8080/")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func formHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		device := Device{
			Class:        r.FormValue("class"),
			Colocation:   r.FormValue("colocation"),
			Description:  r.FormValue("description"),
			Hostname:     r.FormValue("hostname"),
			MgmtIP:       r.FormValue("mgmt_ip"),
			SerialNumber: r.FormValue("serial_number"),
			DeviceType:   r.FormValue("device_type"),
		}

		_, err := db.Exec(`
			INSERT INTO devices (class, colocation, description, hostname, mgmt_ip, serial_number, device_type)
			VALUES ($1, $2, $3, $4, $5, $6, $7)`,
			device.Class, device.Colocation, device.Description, device.Hostname, device.MgmtIP, device.SerialNumber, device.DeviceType)

		if err != nil {
			http.Error(w, "Insert failed", http.StatusInternalServerError)
			log.Println(err)
			return
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Show form and device list
	devices := getAllDevices()
	templates.ExecuteTemplate(w, "form.html", devices)
}

func listDevices(w http.ResponseWriter, r *http.Request) {
	devices := getAllDevices()
	templates.ExecuteTemplate(w, "list.html", devices)
}

func getAllDevices() []Device {
	rows, err := db.Query("SELECT id, class, colocation, description, hostname, mgmt_ip, serial_number, device_type FROM devices ORDER BY id")
	if err != nil {
		log.Println("DB query error:", err)
		return nil
	}
	defer rows.Close()

	var devices []Device
	for rows.Next() {
		var d Device
		if err := rows.Scan(&d.ID, &d.Class, &d.Colocation, &d.Description, &d.Hostname, &d.MgmtIP, &d.SerialNumber, &d.DeviceType); err == nil {
			devices = append(devices, d)
		}
	}
	return devices
}

func editDevice(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	row := db.QueryRow("SELECT id, class, colocation, description, hostname, mgmt_ip, serial_number, device_type FROM devices WHERE id = $1", id)

	var d Device
	err := row.Scan(&d.ID, &d.Class, &d.Colocation, &d.Description, &d.Hostname, &d.MgmtIP, &d.SerialNumber, &d.DeviceType)
	if err != nil {
		http.Error(w, "Device not found", http.StatusNotFound)
		return
	}

	templates.ExecuteTemplate(w, "edit.html", d)
}

func updateDevice(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		id, _ := strconv.Atoi(r.FormValue("id"))
		_, err := db.Exec(`
			UPDATE devices
			SET class=$1, colocation=$2, description=$3, hostname=$4, mgmt_ip=$5, serial_number=$6, device_type=$7
			WHERE id=$8`,
			r.FormValue("class"), r.FormValue("colocation"), r.FormValue("description"),
			r.FormValue("hostname"), r.FormValue("mgmt_ip"), r.FormValue("serial_number"),
			r.FormValue("device_type"), id)

		if err != nil {
			http.Error(w, "Update failed", http.StatusInternalServerError)
			log.Println(err)
			return
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func deleteDevice(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	_, err := db.Exec("DELETE FROM devices WHERE id = $1", id)
	if err != nil {
		http.Error(w, "Delete failed", http.StatusInternalServerError)
		log.Println(err)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func exportCSV(w http.ResponseWriter, r *http.Request) {
	devices := getAllDevices()
	w.Header().Set("Content-Disposition", "attachment; filename=devices.csv")
	w.Header().Set("Content-Type", "text/csv")

	writer := csv.NewWriter(w)
	defer writer.Flush()

	writer.Write([]string{"ID", "Class", "Colocation", "Description", "Hostname", "Mgmt IP", "Serial Number", "Type"})

	for _, d := range devices {
		writer.Write([]string{
			strconv.Itoa(d.ID), d.Class, d.Colocation, d.Description,
			d.Hostname, d.MgmtIP, d.SerialNumber, d.DeviceType,
		})
	}
}

func exportJSON(w http.ResponseWriter, r *http.Request) {
	devices := getAllDevices()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(devices)
}
