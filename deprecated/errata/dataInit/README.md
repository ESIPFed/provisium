### Data Init

#### About
There are just some simple tools to load up some of the mock 
data into the provisium application/kv stores.

For some of these exampls we will use a Frictionless data approach 
denoted with the mimetype ```application/vnd.dataresource+json```.  
For the CSV document this will be ```text/csv```.

#### Framing
During the loading we will process the JSON-LD and using framing
to extract the key elements we need in an easy manner to facilitate 
loading the data sets.   We will really only need the UUID
and the file name.  Frames like the following will help.

```
{
  "@context": {"@vocab": "http://schema.org/"},
  "@type": "Dataset",
  "@explicit": true,
  "name": {},
  "identifier": {
     "@explicit": true,
    "@type": "PropertyValue",
    "value":{}
  }
}
```

This will result in something like
```
{
  "@context": {
    "@vocab": "http://schema.org/"
  },
  "@graph": [
    {
      "@id": "http://provisium.io/id/dataset/bcd15975-680c-47db-a062-ac0bb6e66816",
      "@type": "Dataset",
      "identifier": {
        "@id": "_:b1",
        "@type": "PropertyValue",
        "value": "bcd15975-680c-47db-a062-ac0bb6e66816"
      },
      "name": "208_1262A_JanusThermalConductivity_VyaMsepM.csv"
    }
  ]
}
```


You can also go directly to a branch or branches  of a specific type at any level.
```
{
  "@context": {"@vocab": "http://schema.org/"},
   "@type": "PropertyValue",
   "@explicit": true,
    "value":{}
}
``` 
