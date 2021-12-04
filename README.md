# HTTP Deleter

*Sends HTTP DELETE requests to a list of URLs*

Each URL is parsed with STRFTIME, so can include placeholders for date or time e.g. `http://elasticsearch.example.com:9200/my_index-%Y.%m.%d`

There are various options to influence `time`.

(This is a potentially simpler replacement for Elasticsearch `curator`)

```
  -backcount int
    	Iterate back in time <backcount> times 
    	[HTTP_DELETER_BACKCOUNT] (default 30)
    	
  -backstep int
    	Each <backcount> iteration goes back in time <backstep> hours 
    	[HTTP_DELETER_BACKSTEP] (default 24)
    	
  -concurrent int
    	Number of HTTP requests to send in parallel 
    	[HTTP_DELETER_CONCURRENT] (default 2)
    	
  -dryrun
    	Show parsed target list without sending HTTP Delete request 
    	[HTTP_DELETER_DRYRUN]
    	
  -httptimeout int
    	HTTP request timeout in seconds 
    	[HTTP_DELETER_HTTPTIMEOUT] (default 30)
    	
  -loopdelay int
    	If specified: run continuously with this number of hours between actions.
    	If unspecified: run once and exit 
    	[HTTP_DELETER_LOOPDELAY]
    	
  -startupdelay int
    	Seconds to wait before doing anything (used to facilitate easier testing) 
    	[HTTP_DELETER_STARTUPDELAY]
    	
  -timeoffset int
    	STRFTIME will be <utc-now> minus this number of hours 
    	[HTTP_DELETER_TIMEOFFSET] (default 8760)
    	
  -urls string
    	JSON array of URLs to DELETE. Each will be parsed with STRFTIME 
    	[HTTP_DELETER_URLS]
    	
    	For example:
    	[ "http://host.name:9200/my_index-%Y.%m.%d", "http://other.name:9200/some_index-%Y.%m.%d" ]

```



