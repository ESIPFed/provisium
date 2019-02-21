/index.html             GET
/about.html             GET

/doc/                  303 DX exchange for doc to id
/id/event               PUT  (calls md5hash and keyload)
/id/event/{EVENTID}     GET   (303 from /doc/event/{EVENTID} GET)

/list (range & offset)  GET  Pointless?  Why would I want this?

/search?q={term}        (where term might be HASH in sub or obj place  (or just in value))

/api/md5hash
/api/keyload/ID         (IDs are MD5 hashes)
/api/graphload/ID       calls keyget and then loads the results to internal graph
/api/keyget/ID
/api/graph              (return the entire graph....)


Store in Bolt but also as a graph in https://github.com/wallix/triplestore
