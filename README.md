### Provisium

#### About

We propose to develop software for a server and standalone client that
demonstrates the various approaches to PROV and PROV-AQ as described in the W3C
Working Group specifications. We will demonstrate simple reference
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

Currently there is webapp and datastore.  Datastore will be a RDF triple server
like blazegraph or the like.  Need to resolve if I want to communicate via a
3rd package that could then expose those channels externally.  First thought is
a gRPC tool to all webapp to SPARQL update via a secure channel.  Given the
simplicity of this design, doesn't seem much need though.  So I will need to
resolve some method to update (SPARQL UPDATE I guess) the triple store.  Be
sure to use data volumes or something to persist the updates. 

