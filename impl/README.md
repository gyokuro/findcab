
# Implementations

## Backends

Currently there are two implementations of the `CabService` interface as defined in `service.go`:

    + Simple/ naive implementation: all cabs are stored in memory in a hashmap. Computation of the
nearest cab within location of radium M requires O(N) computations of the haversine distance.
    + MongoDb implementation: this uses mongodb as the database backend and spatial index. The
`2dsphere` index is used on a GeoJSON representation of the Cab struct in the database and
proximity queries are used for the 'within' computations.

MongoDb is used for the following reasons:
    + This application is actually write heavy because each cab is expected to send an update of
its locations at frequent intervals.  Because it's write heavy, backend datastores that also support
spatial indexing (e.g. CouchDb, Lucene, ElasticSearch) are not ideal candidates:
    ++ CouchDb with lots of updates, would require frequent database truncation of the append-only write log.
    ++ Lucene, ElasticSearch are optimized for reads of infrequently updated documents.


## Alternative implementations

Unlike