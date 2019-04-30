# go-sfomuseum-aircraft

Go package for working with aircraft, in an SFO Museum context.

## Tools

### lookup

```
./bin/lookup -h
Usage of ./bin/lookup:
  -source string
    	Valid sources are: icao, sfomuseum. (default "icao")
```

Lookup an aircraft by it's ICAO aircraft designator or manufacturer code.

```
./bin/lookup B737
BOEING B737 "737-700"
BOEING B737 "737-700 BBJ"
BOEING B737 "C-40"
BOEING B737 "C-40 Clipper"
BOEING B737 "Clipper"
BOEING B737 "BBJ (737-700)"
```

```
./bin/lookup EMBRAER
EMBRAER E170 "170"
EMBRAER E190 "190"
EMBRAER AMX "A-1"
EMBRAER TUCA "A-27 Tucano"
EMBRAER E314 "A-29"
EMBRAER E314 "ALX"
EMBRAER AMX "AMX"
...
EMBRAER E290 "ERJ-190-300"
EMBRAER E290 "E190-E2"
EMBRAER E295 "ERJ-190-400"
EMBRAER E295 "E195-E2"
EMBRAER E275 "ERJ-190-500"
EMBRAER E275 "E175-E2"
EMBRAER KC39 "KC-390"
EMBRAER P28U "EMB-711ST Corisco 2 Turbo"
EMBRAER P28U "Corisco 2 Turbo"
```

If you pass the `-source sfomuseum` flag you can lookup aircraft by their SFO Museum (WOF) ID or, where there is a concordance, their ICAO designator code.

```
./bin/lookup -source sfomuseum B744
1159289915 B744 "Boeing 747-400"
```

```
> ./bin/lookup -source sfomuseum 1159289467
1159289467  "Antonov An-8"
```

```
./bin/lookup -source sfomuseum AMX
2019/04/30 13:31:05 Not found
```

## See also

* https://millsfield.sfomuseum.org/aircraft
* https://millsfield.sfomuseum.org/blog/2018/12/03/airlines/
* https://www.icao.int/publications/DOC8643/Pages/Search.aspx