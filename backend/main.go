package main

import (
    "encoding/json"
    "net/http"
    "strconv"
)

type Menu struct {
    ID    int     `json:"id"`
    Name  string  `json:"name"`
    Price float64 `json:"price"`
    Type  string  `json:"type"`
}

var menus = []Menu{
    {ID: 1, Name: "ต้มยำกุ้ง", Price: 120, Type: "soup"},
    {ID: 2, Name: "พิซซ่าฮาวายเอี้ยน", Price: 199, Type: "pizza"},
    {ID: 3, Name: "พิซซ่าเห็ด", Price: 179, Type: "pizza"},
}

var nextID = 4

func listMenu(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    t := r.URL.Query().Get("type")
    if t == "" {
    json.NewEncoder(w).Encode(menus)
    return
    }
    filtered := []Menu{}
    for _, m := range menus {
        if m.Type == t {
            filtered = append(filtered, m)
        }
    }
    json.NewEncoder(w).Encode(filtered)
}


func main() {
    mux := http.NewServeMux()
    mux.HandleFunc("GET /menu", listMenu)
    mux.HandleFunc("GET /menu/{id}", getMenu)
    mux.HandleFunc("POST /menu", createMenu)
    mux.HandleFunc("DELETE /menu/{id}", deleteMenu)
    
    http.ListenAndServe(":8080", withCORS(mux))
}

func getMenu(w http.ResponseWriter, r *http.Request) {
    id, err := strconv.Atoi(r.PathValue("id"))
    if err != nil {
       writeError(w, http.StatusBadRequest, "BAD_ID", "เลขจานต้องเป็นตัวเลข")
        return
    }
    for _, m := range menus {
        if m.ID == id {
            w.Header().Set("Content-Type", "application/json")
            json.NewEncoder(w).Encode(m)
            return
        }
    }
    writeError(w, http.StatusNotFound, "MENU_NOT_FOUND", "ไม่พบเมนูหมายเลขนี้")
}
func createMenu(w http.ResponseWriter, r *http.Request) {
    var m Menu
    if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
        writeError(w, http.StatusBadRequest, "BAD_JSON", "อ่านกล่อง JSON ไม่ออก")
        return
    }
    if m.Name == "" || m.Price <= 0 {
        writeError(w, http.StatusBadRequest, "MISSING_FIELD", "ต้องมีชื่อเมนู และราคาต้องมากกว่าศูนย์")
        return
    }
    m.ID = nextID
    nextID++
    menus = append(menus, m)
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(m)
}
func writeError(w http.ResponseWriter, status int, code, message string) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    json.NewEncoder(w).Encode(map[string]any{
        "error": map[string]string{
            "code":    code,
            "message": message,
        },
    })
}
func deleteMenu(w http.ResponseWriter, r *http.Request) {
    id, err := strconv.Atoi(r.PathValue("id"))
    if err != nil {
        writeError(w, http.StatusBadRequest, "BAD_ID", "เลขจานต้องเป็นตัวเลข")
        return
    }
    for i, m := range menus {
        if m.ID == id {
            menus = append(menus[:i], menus[i+1:]...)
            w.WriteHeader(http.StatusNoContent)
            return
        }
    }
    writeError(w, http.StatusNotFound, "MENU_NOT_FOUND", "ไม่พบเมนูหมายเลขนี้")
}
// พนักงานต้อนรับหน้าประตู (CORS Middleware)
func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		// ตอบกลับ preflight request (OPTIONS) ของเบราว์เซอร์ทันที
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}
