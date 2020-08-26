# Value Villages

This is an open-source project for NGO´s that works in the ground with vocational training.
You can see more information seeing the demo project in kike.me

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes. See deployment for notes on how to deploy the project on a live system. You have full manuals and instructions step-by-step in the demo link above.

### Prerequisites

What things you need to install for the software to run

```
golang
git (only Windows)
```

### Installing

With this command you will install the project and all his dependencies

```
go get -u github.com/kiketordera/valuevillages/...
```

Change the directory to cmd/valuevillages, and there you will find the main.go file to run the project.

To run in Windows
```
set VILLAGE=Central
go run main.go
```

To run in macOS/linux

```
VILLAGE=Central go run main.go
```

Then you can to to localhost:8080/hello in your browser to see the project running (same as demo).
The default username is "admin" and the password is "admin".


## Deployment

Take this repository for developing and expanding/improving features, but not to deployment.
For simplicity, deployment is made all in docker. You have full documentation in the demo itself.

## Built With

* [Gin-gonic](https://github.com/gin-gonic/gin) - The web framework used

## Versioning

This is the first version stable of the app, it is ready for production use.

## Authors

* **Damia Poquet Femenia** - *Initial help* - (https://github.com/DamiaPoquet)

* **Kike Tordera** - *Development* - [PurpleBooth](https://github.com/kiketordera)


## License

Shield: [![CC BY 4.0][cc-by-shield]][cc-by]

This work is licensed under a
[Creative Commons Attribution 4.0 International License][cc-by].

[![CC BY 4.0][cc-by-image]][cc-by]

[cc-by]: http://creativecommons.org/licenses/by/4.0/
[cc-by-image]: https://i.creativecommons.org/l/by/4.0/88x31.png
[cc-by-shield]: https://img.shields.io/badge/License-CC%20BY%204.0-lightgrey.svg


