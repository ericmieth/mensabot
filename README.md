# build

* ```cd $GOPATH```
* ```git clone https://github.com/ericmieth/mensabot```
* ```cd mensabot```
* ```go build mensabot.go```
* ```cp config.json.example config.json``` (update the data)

# run
* ```./mensabot```

# build with docker

* ```export CGO_ENABLED=0```
* ```go build mensabot.go```
* ```docker build -t mensabot .```

# id of canteens and cafeterias

* Cafeteria Dittrichring: 153
* Cafeteria Phillipp-Rosenthal-Straße: 127
* Mensaria Liebigstraße: 162
* Mensa Academica: 118
* Mensa am Park: 106
* Mensa am Elsterbecken: 115
* Mensa Peterssteinweg: 111
* Mensa Schönauer Straße: 140
* Mensa Tierklinik: 170
