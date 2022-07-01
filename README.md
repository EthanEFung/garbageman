# Garbage Man

This is simple tool for my current work environment that helps me manage the
screenshots I take. In my environment, it watches a directory that is populated
using the MacOS Screenshot Application, only keeping the last 10 screenshots
and recordings that were created. The remaining files are moved into Macs
default Trash folder. 

# Dependencies
1. MacOS - Big Sur or Monterey should suffice
2. Recent version of Go - Currently on 1.18

# How to Run
The disclaimer is that this project was developed for personal use. Your
mileage may vary depending on your OS, and version of Go. I would encourage
all to fork this repository and play around with the dependencies to fit your
own needs. To go a step further, you could write it in your favorite language! Otherwise, first fork this github repository, and in your terminal
navigate to this directory running the following commands.

```bash
go mod init github.com/your-handle/garbageman
go mod tidy
cp ./.env.example ./.env
```
This will install the 3rd party pages and create a copy of the
example .env file. Be sure to replace the FROM variable with
the filepath of the directory that your Screenshot app automatically saves to.

## DISCLAIMER
Make sure that the contents of your `FROM` directory does
**NOT** contain any valuable files or folders, otherwise
this process _will_ move the majority of the contents to
the trashbin.

Once the .env file has been set up run the script in
the terminal
```
go run . // Watching for Changes...
```
You should be set up now. Test this by taking a screenshot
waiting for the image to be saved. You should then hear the
Trash sound a couple times assuming that the contents of
the Screenshot directory has more than 10 files.

