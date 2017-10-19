# go-mapzen-valhalla

Minimal viable Go bindings for doing minimal viable things with the Mapzen Valhalla API.

## Install

You will need to have both `Go` (specifically a version of Go more recent than 1.6 so let's just assume you need [Go 1.9](https://golang.org/dl/) or higher) and the `make` programs installed on your computer. Assuming you do just type:

```
make bin
```

All of this package's dependencies are bundled with the code in the `vendor` directory.

## Tools

### valhalla-route

```
./bin/valhalla-route -h
Usage of ./bin/valhalla-route:
  -api-key string
    	A valid Mapzen API key. (default "mapzen-xxxxxx")
  -costing string
    	A valid Valhalla costing. (default "auto")
  -endpoint string
    	A valid Valhalla API endpoint. (default "valhalla.mapzen.com")
  -from string
    	Starting latitude,longitude position.
  -from-wofid int
    	Starting Who's On First ID.
  -to string
    	Destination latitude,longitude position.
  -to-wofid int
    	Destination Who's On First ID.
```

For example:

```
$> valhalla-route -api-key mapzen-xxxx -from 40.759220,-73.987126 -to 40.765852,-73.968329 | jq '.trip.legs[].shape'
"eewvlAnzxblChRwl@nSeo@Nm@~BoHzAuEdAyCd@iB`BcF|s@m|B~Mwb@xBqGjBuFjUkt@xv@qcChBcGbBsFfc@ytAzAeFyCiBsPkL_OmJyf@a\\af@s[i\\}SmEyCkBkAaC{AmTuNyLaIqb@wXiBkAaCiBqQkLaHuEaCyAqB{AiC{Ao^kViBkAyB{AoSwMeJcGsBkA_CkBqRgMyGgEuEwCoXaRuO}Jyb@eYcFgD{e@c[yGuE_TyMaHuEye@q[{B{Am^{UqB{Aqb@wXqB{AyBkAu^kVyB{AyBkAo^kVyByAjAwD|c@gvAbf@r[ze@b[oIfX"
```

Or:

```
$> valhalla-route  -api-key mapzen-xxxxx -from 40.759220,-73.987126 -to 40.765852,-73.968329 | jq '.trip.legs[].maneuvers[].instruction'

"Drive southeast on West 46th Street."
"Turn left onto Madison Avenue."
"Turn right onto East 65th Street."
"Turn right onto Park Avenue."
"Turn right onto East 63rd Street."
"Your destination is on the right."
```
