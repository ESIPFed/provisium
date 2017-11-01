## Status document

### About
Some of the current development goals for the Provisium project



##### Digester
The *Digester* is a remote service (gRPC or FN) that takes a URI LIST and attempts to access the URL and read teh graph.  It should also try and vet the graph as at least well formed.  

Need to resolve how the web server can communicate to something like Digester to keep a graph current.  
Is it a cron job like event, push, pull, other?  Also, need to keep track of what is indexed and record of 
that. 

It then could load it into a triple store or keep a local full graph to work from like some libraries such 
as wallix/triplestore.  The local full graph could then be loaded into a triplestore at any time by anyone.  

> NOTE:  Do make the full graph accessible by anyone


##### Questions and notes

- [ ] Are blank nodes ever going to be present in a valid Prov record?  If so can we just make
it policy they are not allowed?
- [ ] Need some example Prov  (reference RDA examples)
- [ ] 