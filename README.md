# BML
BML is a command-line bookmark management system for quickly managing working directories in your terminal

Currently, only works on Linux. Could potentially work on Mac but I don't have a machine to test with. Windows support is a no no.

## The Problem
Not having the ability to cd into bookmarked working directories has been a bit annoying,  I sometimes found myself forgetting where project directories on my system or found using the cd command constantly when switching between common working directories tedious. So I decided to make my own terminal working directory bookmark management system.

## The Solution
Creating an application that can change the working directory of a terminal session is not as simple as I thought it would. I wanted to write it Golang, but programming an app that could change the working directory of the parent terminal session process was not a trivial thing. In the end, I decided to split the bookmark management system into a Golang app and the change of working directory logic in a bash command stored in my ~/.bashrc. If you know of a better architecture for this, please feel free to share it. Additionally, I chose this approach because I did not want to write the whole thing as one large shell script.

### Images
#### Using the bml
![bml](https://user-images.githubusercontent.com/38519016/80855292-e6484b00-8bf4-11ea-9615-2f9b7dc8beea.gif)
#### Creating new bookmark
![new](https://user-images.githubusercontent.com/38519016/80855299-f8c28480-8bf4-11ea-9df8-f0502e0f3177.gif)
#### Removing a bookmark
![remove](https://user-images.githubusercontent.com/38519016/80855301-feb86580-8bf4-11ea-9f20-71b6d418659b.gif)

## Requirements
All this app requires is an installation of Golang to build and install the Binary.

## Installation
Run the commands below to install bml
```sh
git clone https://github.com/JonAlfaro/BML.git
cd ./BML
go build
./BML install
```

Once install you will have to start a new terminal session. All new terminal sessions after installation will have the `bml` command present.

## Uninstall
```sh
bml uninstall
```

## bml Commands
Supported Commands:
```sh
$ bml
$ bml new
$ bml remove
$ bml uninstall
$ bml help
```
