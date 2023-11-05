### How to run the application
 ```sh
 go build .
 ```

 ```
 edit the config file , and configure the server address and port appropriately , set automigration to true to run  the initial migration then run ./xmserver. After sucessfully migration set the automigrate to false and continue testing. make sure the server address and port , and database address and database port are set appropriately.

 you can also configure the configuration file located using the flags ie ./xmserver path/to/config.toml.

 ```

 ```
 ./xmserver 
 ```