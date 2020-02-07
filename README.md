# gotemplate
template for future projects or projects that will be developed during hackathons

# package structure
    app                            - main application directory 
    ├── cmd                        - everything related to cli commands
    │   ├── command.go             - general parameters, general commander
    │   ├── server.go              - serve command, its description and execution
    │   └── db.go                  - all commands to operate the database
    ├── rest
    │   ├── public                 - public controllers
    │   ├── private                - private controllers
    │   ├── admin                  - admin controllers
    │   └── server.go              - builder for web-server
    ├── store                      - describes everything, which is related to storage functions
    │   └── user                   
    │       └── user.go            - user and his structs, interface for working directly with 
    │                                 database and service methodss
    └── main.go                    - application entrypoint, processes commands and cli-arguments
