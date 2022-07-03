# Garbage Man

This is a personal simple tool for my work environment that helps me manage 
screenshots. It watches a directory that is populated using the MacOS
Screenshot Application, only keeping the last 10 screenshots and recordings
that were created. The remaining files are moved into Macs default
Trash folder. 

# Dependencies
1. MacOS - Big Sur or Monterey should suffice
2. Recent version of Go - Currently on 1.18

# How to Run
The disclaimer is that this project was developed for personal use. Your
mileage may vary depending on your OS, and version of Go. I would encourage
all to fork this repository and play around with the dependencies to fit your
own needs. To go a step further, you could write it in your favorite language!
Otherwise, first fork this github repository, and in your terminal
navigate to this directory running the following commands.

```bash
go mod init github.com/your-handle/garbageman
go mod tidy
go run main.go -from /Valid/Directory
```
This will install the 3rd party pages and begin the process watching for 
file changes on the specified directory path.

## DISCLAIMER
Make sure that the contents of your `FROM` directory does
**NOT** contain any valuable files or folders, otherwise
this process _will_ move the majority of the contents to
the trashbin.

You should be set up now. Test this by setting up your Screenshot
app to save screenshots specified in your `-from` flag. Make
sure the folder has more than 10 files. After saving a screenshot,
check to make sure that only 10 files are in the folder and 
the overflow files can be found in the Trash bin.

If all works well, you can run the following in the command line to have
the module work globally in your environment. 

```bash
go install
cd ~/
garbageman -from /Your/Valid/Directory
```