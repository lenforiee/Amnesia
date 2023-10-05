<p align="center">
  <img src="https://github.com/lenforiee/Amnesia/raw/main/assets/logo_readme.png" alt="Amnesia for Passbolt" />
</p>
<hr>

[![Go Report Card](https://goreportcard.com/badge/github.com/lenforiee/Amnesia)](https://goreportcard.com/report/github.com/lenforiee/Amnesia)
![Version](https://img.shields.io/badge/Version-v0.1.0-blue)

**Not maintained anymore**

## What is Amnesia

Amnesia is a desktop application for the Passbolt password manager! It is:

- 🚀**Simple yet beautiful** with an intuitive user interface.
- ⚡️**Lightweight and performant**, being written in the Go language.
- 🛠**Fully open source** and actively maintained.

![Example Screenshots](https://cdn.discordapp.com/attachments/769679895298310154/1098390339255025724/Untitled.png)

<hr>

## How do I run Amnesia?

We offer multiple ways to get Amnesia running on your system!

### Using Releases

We upload a new release every time there is a major update to the project. These are pre-compiled executables that you can [download from the releases page](https://github.com/lenforiee/Amnesia/releases) and just simply run.

This is the method we recommend for the majority of users.

### Building From Source

Since Amnesia is an open source project, you may also compile the executables from the source code yourself.
To do so, you must first download and install the Go compiler. Then it is a matter of running the following commands:

```sh
# NOTE: before you start any compilation, make sure 
# to follow (for fyne compiler): https://developer.fyne.io/started/#prerequisites 

# Clone the repository
git clone https://github.com/lenforiee/Amnesia

# Go into directory
cd Amnesia

# Download the required modules
go mod download

# Download fyne compiler
go install fyne.io/fyne/v2/cmd/fyne@latest

# Install GTK+ 3 dependencies
# required by gcc compiler (linux only)
sudo apt-get install -y libgtk-3-dev libwebkit2gtk-4.1-dev libappindicator3-dev librsvg2-dev patchelf

# Build to release executable!
# `windows` compiles for Windows
# `darwin` compiles for MacOS
# `linux` compiles for any linux distro
fyne package -os <windows/darwin/linux> -appID "com.lenforiee.amnesia" -icon "assets/logo.png" -name "amnesia" -release

# OR 

# You can build it to normal golang executable
go build

# ON WINDOWS (removes debug console)
go build -ldflags -H=windowsgui
```

We only recommend this method for users interested in contributing to the project, or making small adjustments for themselves.

<hr>

## Credits

- Thank you [@adrplays](https://github.com/adrplays) for the application logo!

## Licence

Amnesia is licenced under the permissive [MIT Licence](https://github.com/lenforiee/Amnesia/blob/main/LICENSE)
