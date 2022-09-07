package logger

import(
	"io"
	"os"
	"log"
	"bytes"
  "github.com/fatih/color"
)
 
 
var (
    Warning *log.Logger
    Info   *log.Logger
		Debug    *log.Logger
    Error   *log.Logger
    ErrorColor = color.New(color.Bold, color.FgRed).SprintFunc()
    InfoColor = color.New(color.Bold, color.FgWhite).SprintFunc()
    DebugColor = color.New(color.Bold, color.FgGreen).SprintFunc()
    WarningColor = color.New(color.Bold, color.FgYellow).SprintFunc()
    Wrt io.Writer
)
func InitLog(debugEnabled bool, route, version string){
  
	f, err := os.OpenFile(route, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	Wrt = io.MultiWriter(os.Stdout, f)
	log.SetOutput(Wrt)
  
	Info = log.New(Wrt, InfoColor("\n[INFO] "), log.Ldate|log.Ltime|log.Lshortfile)
  Info.SetOutput(Wrt)
	
	Warning = log.New(Wrt, WarningColor("\n[WARNING] "), log.Ldate|log.Ltime|log.Lshortfile)
  Warning.SetOutput(Wrt)
  
	Debug = log.New(Wrt, DebugColor("\n[DEBUG] "), log.Ldate|log.Ltime|log.Lshortfile)
  Debug.SetOutput(Wrt)
	if !debugEnabled {
		var buff bytes.Buffer
		Debug.SetOutput(&buff)
	}
	
	Error = log.New(Wrt, ErrorColor("\n[ERROR] "), log.Ldate|log.Ltime|log.Lshortfile)
  Error.SetOutput(Wrt)

  PrintVersion(version)
}

func PrintVersion(version string) {
  f, err := os.OpenFile("SDK.ver", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
  wrt := io.MultiWriter(os.Stdout, f)
  
  versionLog := log.New(wrt, DebugColor("\n[VERSION] "), 0)
  versionLog.SetOutput(wrt)

  versionLog.Println("COBRA ", version)
  
}