```
██████╗ ███████╗██╗██╗     ██╗   ██╗
██╔══██╗██╔════╝██║██║     ██║   ██║
██████╔╝█████╗  ██║██║     ██║   ██║
██╔══██╗██╔══╝  ██║██║     ██║   ██║
██║  ██║███████╗██║███████╗╚██████╔╝
╚═╝  ╚═╝╚══════╝╚═╝╚══════╝ ╚═════╝ 
                                            
                                       
```

# Reilu Terminal Project Management Suite
## Introduction
Reilu Terminal Project Management Suite is a terminal application for managing projects. It doesn't use a database, rather it uses Git (via go-git) to track and manage changes in data. The initial premise of the application was to create a terminal application which imitates traditional software development ticketing and workflow systems. 

### Why use Git instead of database?
_Using Git rather than a database adds a layer of complexity to the development. It's more unituitive to use, and less transparent than a database. So why?_
Why not? YOLO.

## How to run?
No binaries exists yet for this application. 

### Prerequisites
You will need to have Golang installed to run the application. You will also need to set a working directory as a environment variable "CLIRA_ROOTDIR". Otherwise the application will try to initialize a git repo in your current directory and automatically commits file changes, which would ruin the repo structure. A dotenv file is provided for setting this value.

To run the application:

1. Install dependencies. Run: 

   `go mod tidy`

2. Run main.go. Run:

   `go run main.go`


