http://provisium.io/webslides/static/images/provactivities.png# Provisium

## PROV-AQ 

- https://www.w3.org/TR/prov-aq/

Provisium is not a production system not a recommendation for one.  It is simply
an attempt to present an implementation of the W3C Note on PROV-AQ and to assess 
points of implementation and value add aspects.  

![provenance activities][provact]


## Landing Page Implications

The PROV-AQ note defines methods for exposing information about PROV pingback 
(https://www.w3.org/TR/2013/NOTE-prov-aq-20130430/#provenance-pingback) 

### Header setting

`http://opencoredata.org/rdf/graph/JRSO_deployments_gl.ttl.gz` Example URL

```
fils@xps:~$ curl -D headers.txt http://opencoredata.org/rdf/graph/JRSO_deployments_gl.ttl.gz -o download.dat
  % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
                                 Dload  Upload   Total   Spent    Left  Speed
100  191k  100  191k    0     0   661k      0 --:--:-- --:--:-- --:--:--  660k

fils@xps:~$ cat headers.txt 
HTTP/1.1 200 OK
Accept-Ranges: bytes
Content-Length: 196040
Content-Type: text/plain
Date: Mon, 26 Mar 2018 23:41:10 GMT
Last-Modified: Fri, 10 Jun 2016 04:50:47 GMT
Link: <http://opencoredata.org/id/graph/JRSO_deployments_gl.ttl.gz/provenance>; rel="http://www.w3.org/ns/prov#has_provenance"
Link: <http://opencoredata.org/rdf/graph/JRSO_deployments_gl.ttl.gz/pingback>; rel="http://www.w3.org/ns/prov#pingbck"
```

### has_provenance
In the above exchange we get two `Link` entries that provide us a connection to a provenance record. 

`NOTE:` The following is just an example packet that this server returns..  it is not real prov

```
curl  http://opencoredata.org/rdf/graph/JRSO_deployments_gl.ttl.gz/provenance
@prefix ns0:	<http://www.w3.org/ns/prov#> .
@prefix ns1:	<http://foo.org/> .
ns1:thisSample	a	ns0:entity ;
	ns0:qualifiedAttribution	_:bn1 ;
	ns0:wasAttributedTo	_:bn1 .
@prefix ns2:	<http://www.w3.org/1999/02/22-rdf-syntax-ns#> .
@prefix ns3:	<http://opencoredata.org/> .
ns3:org	ns2:label	"Geoscience Australia" ;
	a	ns0:Agent ,
			ns0:org .
_:bn1	a	ns0:attribution ;
	ns0:agent	ns3:org .
@prefix ns4:	<http://www.aurole.org/> .
_:bn1	ns0:hadRole	ns4:Publisher .
```

### pinback
The pingback URL can be used to send back references to prov.  Note in discussion with Tim Lebo (ref: https://twitter.com/timrdf/status/969273759621402625 ) 

```
fils@xps:~$ curl -D headers.txt --data "http://coyote.example.org/contraption/provenance" \\
 -H "Content-Type: text/uri-list" \\
 -X POST  http://opencoredata.org/rdf/graph/JRSO_deployments_gl.ttl.gz/pingback

fils@xps:~$ cat headers.txt 
HTTP/1.1 204 No Content
Date: Tue, 27 Mar 2018 00:07:05 GMT

```


## schema.org alternative (NOTE: *Not* part of W3C Note)

An alternative approach could leverage JSON-LD packages included by means of `script` tags to include
HTML 5 microdata.  These would follow the schema.org publishing patterns.

```javascript
<script type="appplication/ld+json">
{
  "@context": {
    "@vocab": "http://schema.org/",
    "datacite": "http://purl.org/spar/datacite/",
    "earthcollab": "https://library.ucar.edu/earthcollab/schema#",
    "geolink": "http://schema.geolink.org/1.0/base/main#",
    "vivo": "http://vivoweb.org/ontology/core#",
    "dbpedia": "http://dbpedia.org/resource/",
    "geo-upper": "http://www.geoscienceontology.org/geo-upper#",
    "prov": "http://www.w3.org/ns/prov#"
  },
  "@type": "Dataset",
  "prov:has_anchor": "http://geodex.org/datapackages/p418graph.zip",
  "prov:has_provenance": "http://provohash.xyz/id/hash/adbef8ab-04c6-4cb9-b11b-003124c40004",
  "prov:pingback": "http://provisium.io/pingback",
  "prov:has_query_service": "http://provisium.io/provenance-query-service",
  "additionalType": ["http://schema.geolink.org/1.0/base/main#Dataset", "http://vivoweb.org/ontology/core#Dataset"],
  "name": "Project 418 Graph",
  "description": "This data set includes RDF triples from the harvesting process for P418",
  "url": "https://geodex.org/datasets/p418graph.zip",
  "version": "2013-11-21",
  "keywords": "RDF Graph",
  "license": "CC0-1.0"
}
</script>
```

In this above example the *has_provenance* link points to an instance of prov stored at a 3rd party service.
In this case, the test site provohash.xyz (part of the provisium development)

It is possible to leverage JSON-LD Framing to extract only the prov data with a frame like:

```
{
  "@context": {"@vocab": "http://schema.org/",     
               "prov":"http://www.w3.org/ns/prov#"},
  "@type": "Dataset",
  "@explicit": true,
  "prov:has_anchor": {},
  "prov:has_provenance": {},
  "prov:pingback": {},
  "prov:has_query_service": {}
}
```

resulting in

```
{
  "@context": {
    "@vocab": "http://schema.org/",
    "prov": "http://www.w3.org/ns/prov#"
  },
  "@graph": [
    {
      "@id": "_:b0",
      "@type": "Dataset",
      "prov:has_anchor": "http://geodex.org/datapackages/p418graph.zip",
      "prov:has_provenance": "http://provohash.xyz/id/hash/adbef8ab-04c6-4cb9-b11b-003124c40004",
      "prov:has_query_service": "http://provisium.io/provenance-query-service",
      "prov:pingback": "http://provisium.io/pingback"
    }
  ]
}
```

The framing is just a reliable mechanism to allow extracting a subset of the JSON-LD and then likely 
map it to a native data structure in a processing program.  

Looking at these entries though, we see another topic of interest.  The holders of prov and 
the providers of PROV-AQ services and functions need *not* be the originating domain.  For example we
could have a site like figshare or the Provohash domain set up in Provisium.


## Provohash.xyz

Provohash is a simple site that takes a prov record and returns a hosting URL for it.  This allows
sites to avoid hosting the prov.  A site this can also:

* validate the PROV package for validity and form.  This might be simply ensuring it is RDF or applying 
SHACL rules to asses or impose compliance. 
* support pingback callbacks
* perform aggregation on pingback pass by reference URI lists  (aggregate and or validate these resources)

```
curl -D headers.txt -d @prov.rdf  -H "Content-Type: rdf/nquads" -X POST   http://provohash.xyz/doc/newprov
http://provohash.xyz/id/hash/adbef8ab-04c6-4cb9-b11b-003124c40004

fils@xps:~$ cat headers.txt 
HTTP/1.1 200 OK
Content-Length: 65
Content-Type: text/plain; charset=utf-8
Date: Tue, 27 Mar 2018 00:45:49 GMT
```

`NOTE:` This is a link to a provenance landing page
```
curl http://provohash.xyz/id/hash/d413b743179b19359e8732e506e4346b 
```

`NOTE:` This simply returns the prov via an API call (should be a content negotiation on the above URL?)
```
curl http://provohash.xyz/api/v1/prov/7fa64a66-d94d-430f-9713-c2069f84beaf
```


## Provisium.io 

Headers only curl:  curl -v -o /dev/null http://localhost:9900/doc/dataset/bcd15975-680c-47db-a062-ac0bb6e66816

`https://provisium.io/doc/dataset/bcd15975-680c-47db-a062-ac0bb6e66816`


#### ProvPingBack POST exmaple
```
POST http://127.0.0.1:9900/doc/dataset/bcd15975-680c-47db-a062-ac0bb6e66816/pingback HTTP/1.1
content-type: text/uri-list

http://coyote.example.org/contraption/provenance
http://coyote.example.org/another/provenance
```

#### LINK test no content
```
POST http://127.0.0.1:9900/doc/dataset/bcd15975-680c-47db-a062-ac0bb6e66816/pingback HTTP/1.1
Link: <http://coyote.example.org/sparql>; rel="http://www.w3.org/ns/prov#has_query_service"; anchor="http://127.0.0.1:9900/doc/dataset/bcd15975-680c-47db-a062-ac0bb6e66816"
Content-Type: text/uri-list
Content-Length: 0
```

#### LINK array test no content
```
POST http://127.0.0.1:9900/doc/dataset/bcd15975-680c-47db-a062-ac0bb6e66816/pingback HTTP/1.1
Link: <http://coyote.example.org/sparql>; rel="http://www.w3.org/ns/prov#has_query_service"; anchor="http://127.0.0.1:9900/doc/dataset/bcd15975-680c-47db-a062-ac0bb6e66816", <http://coyote.example.org/sparql>; rel="http://www.w3.org/ns/prov#has_query_service"; anchor="http://127.0.0.1:9900/doc/dataset/bcd15975-680c-47db-a062-ac0bb6e66816"
Content-Type: text/uri-list
Content-Length: 0
```

#### LINK test with content
```
POST http://127.0.0.1:9900/doc/dataset/bcd15975-680c-47db-a062-ac0bb6e66816/pingback HTTP/1.1
Link: <http://coyote.example.org/extra/provenance>;rel="http://www.w3.org/ns/prov#has_provenance";anchor="http://127.0.0.1:9900/doc/dataset/bcd15975-680c-47db-a062-ac0bb6e66816"
Content-Type: text/uri-list

http://coyote.example.org/contraption/provenance
http://coyote.example.org/another/provenance
http://coyote.example.org/extra/provenance
```

#### API call..   get prov record details
```
http://127.0.0.1:9900//api/v1/provenance/service?target=http://provisium.io/id/dataset/bcd15975-680c-47db-a062-ac0bb6e66816
```

### FAIL REQUIRED
#### ProvPingBack GET exaple  (SHOULD FAIL)
```
http://127.0.0.1:9900/doc/dataset/bcd15975-680c-47db-a062-ac0bb6e66816/pingback
```
#### no content type   (SHOULD FAIL)
```
POST http://127.0.0.1:9900/doc/dataset/bcd15975-680c-47db-a062-ac0bb6e66816/pingback HTTP/1.1

http://coyote.example.org/contraption/provenance
http://coyote.example.org/another/provenance
```
#### bad URL format   (SHOULD FAIL)
```
POST http://127.0.0.1:9900/doc/dataset/bcd15975-680c-47db-a062-ac0bb6e66816/pingback HTTP/1.1
content-type: text/uri-list

coyote.example.org//contraption&provenance/
```

## Review and Future work on this package

![provenance activities][provact]

[provact]:provactivities.png "Logo Title Text 2"


* provisium.io hosted fake datasets to *prime the pump*.  With the existence of the http://geodex.org/data/catalog and the http://opencoredata.org/catalog/geolink there are examples of both schema.org and header patterns.  The fake data should be removed.  *remove the mock data from provisium.io*
* The separate provohash and provisium sites are confusing.  The single provisium.io domain can address both hosting prov and providing the services around it.  There is not need to host datasets as this pattern attempts to promote the idea that facilities host the data but outsource the prov elements.  *remove the data and combine the functions of the two domains into provisium.io*
* The service query and sparql are not implemented in manners that can be exposed well
* There is no *roll up* or aggregate and validate function at provisium.io to address the pass by reference pattern
* Connect this into some of the P418 concepts and flows

