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
- [ ] Flesh out the data store for the local and uploaded graph fragments


##### Next steps

- [ ]  Modify pingback URL to also accepts RDF based on content type (ttl)
- [ ]  Build "provocite"


##### Thoughts and notes

##### sharing
- You really don't want to give prove a DOI or other UUID since you can never say you 
"are" the prov.  There is always someone else who has done something with a thing and recorded it.  
There activity is just as valid.   So maybe it's more a HASHID for prov?  I have a prov record, I made a hash of it 
and stored it at my domain.  So I will now expose a hash of it and other can refernce that hash. ???

- What is prov was a shared thing.  Using things like [1] or [2].

[1] https://developers.google.com/web/updates/2016/09/navigator-share
[2] https://paul.kinlan.me/navigator.share/ 
[3] http://caniuse.com/#feat=web-share

```javascript
if (navigator.share) {
  navigator.share({
      title: 'Web Fundamentals',
      text: 'Check out Web Fundamentals — it rocks!',
      url: 'https://developers.google.com/web',
  })
    .then(() => console.log('Successful share'))
    .catch((error) => console.log('Error sharing', error));
}
```
