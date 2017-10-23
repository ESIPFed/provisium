### Datastore readme

The provisum data store will be a triple store of some sort.  Likely this will 
be a blazegraph instance with read/write capacity.  Since the sparql end point
needs to be open in some cases for PROV-AQ a ledger of RDF events will be 
kept.  This may be a Mongo store or s3 based approach.


