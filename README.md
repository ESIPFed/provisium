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
Provisium is only just starting.  Leveraging off work done in the Open Core Data project an initial implementation of the basic [PROV-AQ pingback](https://www.w3.org/TR/2013/NOTE-prov-aq-20130430/#provenance-pingback) approach has been implemented.  

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


