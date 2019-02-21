## ProvoHash

### About
This is a simple app to register and share prov.  POST your prov and 
get a MD5 hash to share with.   Also will build a simple web ui for 
for doing that too.   There is no issue wit security at this time 
since this is a lab.   I will require the content be check as
valid parsable RDF.  I will only handle turtle at this time.  


### Sharing
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

