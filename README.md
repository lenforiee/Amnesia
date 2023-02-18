<p align="center">
  <img src="https://github.com/lenforiee/AmnesiaGUI/raw/main/assets/logo_readme.png" alt="Amnesia for Passbolt" />
</p>
<hr>

[![Go Report Card](https://goreportcard.com/badge/github.com/lenforiee/AmnesiaGUI)](https://goreportcard.com/report/github.com/lenforiee/AmnesiaGUI)
![Version](https://img.shields.io/badge/Version-v0.0.1-blue)

## What is Amnesia
Amnesia is a desktop application for the Passbolt password manager! It is:
- 🚀**Simple yet beautiful** with an intuitive user interface.
- ⚡️**Lightweight and performant**, being written in the Go language.
- 🛠**Fully open source** and actively maintained.

![Example Screenshots](https://user-images.githubusercontent.com/36131887/219877620-d3c0d9a3-335a-4fc9-ae3d-ff4e72516cd1.png)

<hr>

## How do I run Amnesia?
We offer multiple ways to get Amnesia running on your system!

### Using Releases
We upload a new release every time there is a major update to the project. These are pre-compiled executables that you can [download from the releases page](https://github.com/lenforiee/AmnesiaGUI/releases) and just simply run.

This is the method we recommend for the majority of users.

### Building From Source
Since Amnesia is an open source project, you may also compile the executables from the source code yourself.
To do so, you must first download and install the Go compiler. Then it is a matter of running the following commands:

```sh
# Clone the repository
git clone https://github.com/lenforiee/AmnesiaGUI

# Go into directory
cd AmnesiaGUI

# Download the required modules
go get

# Compile the project executable
go build

# If you building on windows use this command to make console not appear
go build -ldflags -H=windowsgui
```

We only recommend this method for users interested in contributing to the project, or making small adjustments for themselves.

<hr>

## Credits
- Thank you @adrplays for a application logo!

## Licence
Amnesia is licenced under the permissive [MIT Licence](https://github.com/lenforiee/AmnesiaGUI/blob/main/LICENSE)
