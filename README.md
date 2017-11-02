# go-mapzen-valhalla

Minimal viable Go bindings for doing minimal viable things with the Mapzen Valhalla Turn-by-Turn API.

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

You can also pass the `valhalla-route` tool [Who's On First](https://whosonfirst.mapzen.com) IDs to route to and from. For example:

```
./bin/valhalla-route -api-key mapzen-xxxxxx -from-wofid 1108798585 -to-wofid 270097925 | jq '.trip.legs[].shape'
"cajagAryumhF~H}JrVq\\tYm_@jBhClO~R~kAy~A~kAi_Bju@rcAln@p{@zZxa@xa@jk@jG]dELvDl@zE|@tExA|ExBzEfDzFfErFbFpGpG`H`HzG~HhWb\\|DbFvHlJze@nh@|JvN|JtOdJtOnIvNfItOdcAzqBlEzL~C`HnDpGjF~HrF`HbGpHvH~HfI~HpHrG~HpGvIbGfX`RnDfC~CxCnCvCxCdEnCvDjBtEbAdEl@bF`Bb\\l@~HdAxLrKvNzd@bo@hRhWx\\de@p\\dd@`R{VtO_Sli@_r@tn@o}@~NoStEuE|NeO~HsFzLsFhL]bR}@fg@yBnS}@zPk@li@iC~SkAbAfc@dAng@`Bbp@jBdx@pB`|@pB`{@fS{@"
```

Valhalla polylines can then be handed off to the [Mapzen Places `mapzen.places.getByPolyline` API method](https://mapzen.com/documentation/places/methods/#mapzen.places.getByPolyline).

For example, here are the unique [postal codes](https://whosonfirst.mapzen.com/spelunker/placetypes/postalcode) you would pass through traveling from the [Mapzen offices](https://places.mapzen.com/id/1108798585/), in downtown San Francisco, to [Zeitgeist](https://places.mapzen.com/id/270097925), in the Mission:

```
curl -s 'https://places.mapzen.com/v1?api_key=mapzen-******&method=whosonfirst.places.getByPolyline&polyline=cajagAryumhF%7EH%7DJrVq%5CtYm_%40jBhClO%7ER%7EkAy%7EA%7EkAi_Bju%40rcAln%40p%7B%40zZxa%40xa%40jk%40jG%5DdELvDl%40zE%7C%40tExA%7CExBzEfDzFfErFbFpGpG%60H%60HzG%7EHhWb%5C%7CDbFvHlJze%40nh%40%7CJvN%7CJtOdJtOnIvNfItOdcAzqBlEzL%7EC%60HnDpGjF%7EHrF%60HbGpHvH%7EHfI%7EHp%5C+HrG%7EHpGvIbGfX%60RnDfC%7ECxCnCvCxCdEnCvDjBtEbAdEl%40bF%60Bb%5Cl%40%7EHdAxLrKvNzd%40bo%40hRhWx%5Cde%40p%5Cdd%40%60R%7BVtO_Sli%40_r%40tn%40o%7D%40%7ENoStEuE%7CNeO%7EHsFzLsFhL%5DbR%7D%40fg%40yBnS%7D%40zPk%40li%40iC%7ESkAbAfc%40dAng%40%60Bbp%40jBdx%40pB%60%7C%40pB%60%7B%40fS%7B%40&precision=6&unique=1&placetype=postalcode' | python -mjson.tool
{
    "cursor": null,
    "next_query": null,
    "page": 1,
    "pages": 1,
    "per_page": 100,
    "places": [
        [
            {
                "wof:country": "US",
                "wof:id": 554784673,
                "wof:name": "94105",
                "wof:parent_id": "85922583",
                "wof:placetype": "postalcode",
                "wof:repo": "whosonfirst-data-postalcode-us"
            },
            {
                "wof:country": "US",
                "wof:id": 554784667,
                "wof:name": "94103",
                "wof:parent_id": "85922583",
                "wof:placetype": "postalcode",
                "wof:repo": "whosonfirst-data-postalcode-us"
            },
            {
                "wof:country": "US",
                "wof:id": 554784675,
                "wof:name": "94107",
                "wof:parent_id": "85922583",
                "wof:placetype": "postalcode",
                "wof:repo": "whosonfirst-data-postalcode-us"
            },
            {
                "wof:country": "US",
                "wof:id": 554784681,
                "wof:name": "94110",
                "wof:parent_id": "85922583",
                "wof:placetype": "postalcode",
                "wof:repo": "whosonfirst-data-postalcode-us"
            }
        ]
    ],
    "stat": "ok",
    "total": null
}
```

_Note that a) the polyline in the API call is URL-encoded and b) we are explicitly passing a `precision=6` parameter to distinguish the polyline from the usual 5-decimal precision polylines that Google uses._

## See also

* https://mapzen.com/developers
* https://mapzen.com/documentation/mobility/turn-by-turn/overview/
* https://mapzen.com/documentation/places/
* http://whosonfirst.mapzen.com
