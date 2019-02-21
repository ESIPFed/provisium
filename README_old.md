### Provisium

#### About

A project to develop a server and standalone client that
demonstrates the various approaches to PROV and PROV-AQ as described in the W3C
Working Group specifications (https://www.w3.org/TR/prov-aq/). 
Provisium will demonstrate simple reference
implementations of these approaches . The reference
implementations will:

1) serve as a working examples (with open source code) of
how end-to-end provenance systems could be constructed and 
2) provide a basis
to compare and contrast practical implementations.  

As part of this work, a generic visualization service will be created following
W3C web component guidelines (https://www.w3.org/wiki/WebComponents/). Web
components allow web page authors to import a web visualization, including all
dependencies and networking code, with a few simple import statements.  This
embedding works across domains as well and provides additional reuse options as
the visualization service can easily be integrated inside of dataset landing
and other pages.


#### Resources
Some resources for those interested in web architecture based PROV.

1) https://www.w3.org/TR/prov-aq/
2) https://www.rd-alliance.org/groups/provenance-patterns-wg 
3) http://wiki.esipfed.org/index.php/Semantic_Technologies 


#### Notes
Provisium is only just starting.  Leveraging off work done in the Open Core Data project an initial implementation of the 
basic [PROV-AQ pingback](https://www.w3.org/TR/2013/NOTE-prov-aq-20130430/#provenance-pingback) approach has been implemented.  

An example of this looks when accessing a resource is shown below with the link elments in the HEAD.

```bash
Fils:dataInit dfils$ curl -s -D - http://127.0.0.1:9900/doc/dataset/bcd15975-680c-47db-a062-ac0bb6e66816 -o /dev/null

HTTP/1.1 200 OK
Link: <http://127.0.0.1:9900/doc/dataset/bcd15975-680c-47db-a062-ac0bb6e66816/provenance>; rel="http://www.w3.org/ns/prov#has_provenance"
Link: <http://127.0.0.1:9900/doc/dataset/bcd15975-680c-47db-a062-ac0bb6e66816/pingback>; rel="http://www.w3.org/ns/prov#pingbck"

Date: Fri, 27 Oct 2017 17:58:49 GMT
Content-Type: text/html; charset=utf-8
Transfer-Encoding: chunked
```

#### Issues

##### hasProv 
We can have many ``` hasProv ``` header items.  The question is should we link to non-origin prov traces with this?  It would 
seem to imply a certain level of vetting that is not there.

##### query service
It's likely a site would not have both a hasProv and query service.  For this test we will.  It should be noted that the URI of
a prov record can then not contain # or & or ? as this would effect the structure below.  So, no prov record can be part of a 
RESTful API with query parameters so defined.  

from the docs
```
A server should not offer a template containing {+uri} or other non-simple variable expansion options [URI-template] unless all 
valid target-URIs for which it can provide provenance do not contain problematic characters like '#' or '&'.
```

query service example
```- Link: <service-URI>;
  rel="http://www.w3.org/ns/prov#has_query_service";
  anchor="target-URI"
```


##### hasAnchor
The  link element (#has_anchor) specifies an identifier for the document that may be used within the provenance record when referring to the document.
This is preferably the PID of the resource.

examples for has_query_service and has_anchor
```
<html xmlns="http://www.w3.org/1999/xhtml">
     <head>
        <link rel="http://www.w3.org/ns/prov#has_query_service" href="service-URI">
        <link rel="http://www.w3.org/ns/prov#has_anchor" href="target-URI">
        <title>Welcome to example.com</title>
     </head>
     <body>
       <!-- HTML content here... -->
     </body>
  </html>
```

The query service points to a dereferenced query description assumed options include
*  opensearch description document
*  swagger description document
*  schema.org JSON-LD with provides service  (explore this one)


#### Questions

* should I content negotiate as RDF (JSON-LD) and put in this namespace..
* extend the JSON-LD to include the various prov namespace hashed URI terms
* make a security card on provohash (and provisium) to talk about the security done and NOT done
 (limits, RDF check, etc) and not done (auth, same-origin, etc)


RDF example
```
@prefix prov: <http://www.w3.org/ns/prov#>.

<> dcterms:title        "Welcome to example.com" ;
   prov:has_anchor       <http://example.com/data/resource.rdf> ;
   prov:has_provenance   <http://example.com/provenance/resource.rdf> ;
   prov:has_query_service <http://example.com/provenance-query-service/> .

   # (More RDF data ...)
```


The direct HTTP query service may return provenance in any available format. For interoperable provenance publication, 
use of PROV represented in any of its specified formats is recommended. Where alternative formats are available, selection 
may be made by content negotiation, using Accept: header fields in the HTTP request. Services must identify the Content-Type of the provenance returned.



