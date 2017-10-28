### Provisium

#### About

We propose to develop software for a server and standalone client that
demonstrates the various approaches to PROV and PROV-AQ as described in the W3C
Working Group specifications (https://www.w3.org/TR/prov-aq/). 
We will demonstrate simple reference
implementations of these approaches in this project. The reference
implementations will 1.) serve as a working examples (with open source code) of
how end-to-end provenance systems could be constructed and 2.) provide a basis
to compare and contrast practical implementations.  

As part of this work, a generic visualization service will be created following
W3C web component guidelines (https://www.w3.org/wiki/WebComponents/). Wed
components allow web page authors to import a web visualization, including all
dependencies and networking code, with a few simple import statements.  This
embedding works across domains as well and provides additional reuse options as
the visualization service can easily be integrated inside of dataset landing
and other pages.



#### Notes
Initial testing so far results in behavior where the response is something like

```
Fils:dataInit dfils$ curl -s -D - http://127.0.0.1:9900/doc/dataset/bcd15975-680c-47db-a062-ac0bb6e66816 -o /dev/null
HTTP/1.1 200 OK
Link: <http://127.0.0.1:9900/doc/dataset/bcd15975-680c-47db-a062-ac0bb6e66816/provenance>; rel="http://www.w3.org/ns/prov#has_provenance"
Link: <http://127.0.0.1:9900/doc/dataset/bcd15975-680c-47db-a062-ac0bb6e66816/pingback>; rel="http://www.w3.org/ns/prov#pingbck"
Date: Fri, 27 Oct 2017 17:58:49 GMT
Content-Type: text/html; charset=utf-8
Transfer-Encoding: chunked
```


