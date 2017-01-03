# transfer.sh GUI
This is a little project which adds a context menu item to the Windows explorer to conveniently upload files to transfer.sh. It copies the URL to the clipboard afterwards.

# What is transfer.sh?
[Transfer.sh](https://github.com/dutchcoders/transfer.sh "Github Repo") is an easy and fast file sharing service designed to be used from the command-line.

# What is this?
Sometimes I need to quickly transfer a file from my Windows working machine to a server without having to set up an FTP server. That's why I created this project.
Please don't be too harsh on me since I began coding Go just recently, so this code may be a bit messy and unprofessional :)
Also, I'm not responsible for anything you do using my code.

# Notes
- Make sure you compile it using `go build -ldflags="-H windowsgui" github.com/PlatzhalterDE/transfer.sh-gui` to get rid of the cmd prompt
