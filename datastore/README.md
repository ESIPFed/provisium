### Datastore readme

The provisum data store will be a triple store of some sort at some point.  Likely this will 
be a blazegraph instance with read/write capacity.  Since the sparql end point
needs to be open in some cases for PROV-AQ a ledger of RDF events will be 
kept.  This may be a Mongo store or s3 based approach.

The web apps uses a KV store (BoltDB) for most of its data needs.  This is partly due to 
the need in a lab project to keep all events as unique things.  We will roll data into the 
triple store though to support some of the PROV-AQ patterns that we said we would support.

It will also give us a chance to review triple stores and SPARQL in comparison to other
patterns in this project. 


