package logger

import(
	"io"
	"os"
	"log"
	"bytes"
  "github.com/fatih/color"
)
 
 
var (

    Warning*log.Logger
    Info   *log.Logger
		Debug    *log.Logger
    Error   *log.Logger
    ErrorColor = color.New(color.Bold, color.FgRed).SprintFunc()
    InfoColor = color.New(color.Bold, color.FgWhite).SprintFunc()
    DebugColor = color.New(color.Bold, color.FgGreen).SprintFunc()
    WarningColor = color.New(color.Bold, color.FgYellow).SprintFunc()

)
func InitLog(debugEnabled bool, route string){
  
	f, err := os.OpenFile(route, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	wrt := io.MultiWriter(os.Stdout, f)
	log.SetOutput(wrt)
  
	Info = log.New(wrt, InfoColor("\n[INFO] "), log.Ldate|log.Ltime|log.Lshortfile)
  Info.SetOutput(wrt)
	
	Warning = log.New(wrt, WarningColor("\n[WARNING] "), log.Ldate|log.Ltime|log.Lshortfile)
  Warning.SetOutput(wrt)
  
	Debug = log.New(wrt, DebugColor("\n[DEBUG] "), log.Ldate|log.Ltime|log.Lshortfile)
  Debug.SetOutput(wrt)
	if !debugEnabled {
		var buff bytes.Buffer
		Debug.SetOutput(&buff)
	}
	
	Error = log.New(wrt, ErrorColor("\n[ERROR] "), log.Ldate|log.Ltime|log.Lshortfile)
  Error.SetOutput(wrt)

  Info.Println("version: 1.1.1\nLoading....")

}