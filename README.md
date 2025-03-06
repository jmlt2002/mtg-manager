# MTG-Manager 🃏  
A small CLI application for managing Magic: The Gathering (MTG) libraries and creating custom cards.  
Built with Go, featuring a server-side API with JWT authentication and a SQLite database for storage.  

## Features  
- Manage MTG card libraries  
- Create and store custom cards  
- Secure authentication with JWT  

---

## 📦 Installation  

### **Prerequisites**  
Ensure you have **Go** installed:  
```sh
go version
```
If not installed, download it from [Go’s official site](https://go.dev/dl/).  

Ensure you have **SQLite** installed:  
```sh
sqlite3 --version
```
If not installed, download it from:  
- **Linux/macOS:** Install via package manager  
  ```sh
  sudo apt install sqlite3    # Debian/Ubuntu  
  brew install sqlite         # macOS (Homebrew)  
  ```  
- **Windows:** Download from [SQLite Official Site](https://www.sqlite.org/download.html).  


### **Install Dependencies** 
In both the client and the server directories, run:
```sh
go mod tidy
```

---

## Usage  

### **Starting the Server**  
```sh
go run server/main.go
```
The server will run on **http://localhost:8080**. You can change this by editing server/main.go

### **Using the CLI**  
Run the CLI to interact with the server:  
```sh
go run client/main.go
```
