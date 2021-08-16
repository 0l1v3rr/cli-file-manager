# CLI File Manager
<img src="./docs/screenshot.png">
<br>
This is a basic file manager that runs inside your terminal. <br>
It's fully responsive and incredibly fast.

## Features
- Browse directories/files
- Disc usage panel
- File and Folder information
- Open files (With the default program of the OS)
- Delete files or folders
- Rename files or folders
- Create files or folders
- Responsive

## Installation
Download [Git](https://git-scm.com/downloads) and [Go](https://golang.org/dl/)
### - Linux!:
```sh
git clone https://github.com/0l1v3rr/cli-file-manager.git
cd cli-file-manager
make build
export PATH="$PATH:/path/to/your/cli-file-manager/bin"
```
Now try to execute the command `cfm`. <br>
In the first argument, you can give the path where you want it to open. (Not necessary)<br>
**For example:** `cfm /home/user/Desktop`<br><br>
*Note: The `cfm` command will only live until you close the terminal.* To resolve this issue, follow these few steps:
```sh
nano ~/.profile
```
Now scroll to the bottom and paste this:
```sh
export PATH="$PATH:/path/to/your/cli-file-manager/bin"
```
Save the changes and reboot.
<br><br>
There's an update and do you want to use it?<br>
Go to your cli-file-manager folder. Open a terminal and type this command:
```sh
make update
```
