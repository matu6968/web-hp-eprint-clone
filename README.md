# HP ePrint like service clone (web edition) 

This is a HP ePrint clone that has similar functions to normal HP ePrint but its on the web and it is not HP printer specific

## Prerequisites

- Go (1.23.1 or later)

## Installation

1. Clone the repository:
   ```
   git clone https://git.fluffy.pw/matu6968/web-hp-eprint-clone
   ```

2. Go to the project directory:
   ```
   cd web-hp-eprint-clone
   ```

3. Build the binary:
   ```
   go build -o eprintclone
   ```
   
4. Execture the binary:
   ```
   ./eprintclone
   ```

## Configuration

In the .env file this is the only thing you can set

```
PORT=8080
LOG=log # enables file logging (aka copies file to user home directory under the folder imagelog so make a folder first in the home root if you wish to enable it), to disable it replace it with nolog
```
## Autostart with systemd or OpenRC

You can autostart the web server using systemd or OpenRC (init scripts are in the init-scripts folder)
To use it, edit the script accordingly (edit username on what user it is going to run and the path to the binary on where it will run from)

## for systemd edit the following lines:

```
; Don't forget to change the value here!
; There is no reason to run this program as root, just use your username
User=examplename # change this to your username

WorkingDirectory=/path/to/binary # change this to the path where the binary resides
```
### and to add it as a service:

```
sudo cp /path/to/cloned/repo/init-scripts/eprint-clone-web.service /etc/systemd/system/
sudo systemctl daemon-reload
sudo systemctl enable eprint-clone-web.service
sudo systemctl start eprint-clone-web.service
```

## for OpenRC edit the following lines:

```
command="bash -c cd ~/web-hp-eprint-clone && ./eprintclone" # if you have put the eprintclone binary somewhere else change this line
# Don't forget to change the value here!
# There is no reason to run this program as root, just use your username
command_user="userexample" # change this to your usernames
```

### and to add it as a service:

```
sudo cp /path/to/cloned/repo/init-scripts/eprint-clone-web /etc/init.d/
sudo rc-update add eprint-clone-web
sudo rc-service eprint-clone-web start
``` 

### For this to even work
You need a linux distro which has the modern equifelent of the lpr program. To check if you have the newer version, type `man lpr` and look at the program description, if it says `lpr - print files` (present on recent Ubuntu versions and rolling release distros like Arch Linux) then you are good to go otherwaise if it says `lpr - off line print` (present on Debian and older Ubuntu versions) then this won't work as the commands for the older version are different and this targets the newer version of the program.


Then setup a printer on your host by connecting it over USB or WiFi and then adding the printer in your distro
and put the eprintcloned script (ePrint clone daemon aka it handles the printing commands) in the "/usr/bin" folder 

# !IMPORTANT! 

You need to modify the script to include your printer name (it is at the top and is called printer and not the one used in the included configuration in the repository

for example change it from
```
#!/bin/bash

printer=examplename # replace this with your actual printer name before putting it into /usr/bin
```
to
```
#!/bin/bash

printer=Deskjet-2600 # replace this with your actual printer name before putting it into /usr/bin
```

