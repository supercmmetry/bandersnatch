
<a href="https://documenter.getpostman.com/view/7086087/SzS8rjbV?version=latest"><img src="https://img.shields.io/badge/-Documentation-black?logo=postman"></a>
<p align="center"><img src="https://i.ibb.co/5xcNxBK/gopher.png" alt="gopher"></p>
<br>
<p align="center"><img src="https://i.ibb.co/Pr917y2/bandersnatch.png" alt="bandersnatch" border="0"></p>
<br>
<h4> Steps to use </h4>

- Clone the repository

`git clone github.com/supercmmetry/bandersnatch`
<br>

- Change your working directory

`cd bandersnatch`

- Create a .env file with necessary key-value pairs.

| Key          	| Value                                                                                 	|
|--------------	|---------------------------------------------------------------------------------------	|
| DB_URI       	| "postgresql://localhost/bandersnatch?user=postgres&password=postgres&sslmode=disable" 	|
| DEBUG        	| "true"                                                                                	|
| PORT         	| "1729"                                                                                	|
| JWT_PASSWORD 	| "password"                                                                            	|
| NEXUS_FILE   	| "sample.json"                                                                         	|

- Manage your go modules

 `go mod tidy`

- Create a database in postgresql

- Run your go program

`go run main.go`
