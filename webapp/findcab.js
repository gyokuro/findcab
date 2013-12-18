var currentLocation = new google.maps.LatLng(37.79, -122.393);
var map = null;
var homeMarker = null;
var host = null;
var operation = "view";
var rest = null
var cabs = {}
var queryCircle = null;


$(function(){

    host = location.hostname
    rest = location.protocol + "//" + host + ":8080"

    function setInRange(marker) {
	marker.setIcon("/assets/ico/cab_in_range.png")
    }

    function setOutOfRange(marker) {
	marker.setIcon("/assets/ico/cab.png")
    }

    function resetCabs() {
	for (key in cabs) {
	    setOutOfRange(cabs[key])
	}
    }

    function upsertCab(id, loc) {
	var cab = {
	    id : id,
	    latitude: loc.lat(),
	    longitude: loc.lng()
	}
	var url = [ rest, "/cabs/" + id ].join('')
	jQuery.post(url, JSON.stringify(cab)).always(function(e) {
	    if (e.status == 200) {
		m = new google.maps.Marker({
		    position: loc,
		    map: map,
		    draggable: true,
		    clickable: true,
		    title: 'cab ' + id
		});
		setOutOfRange(m)
		google.maps.event.addListener(m, 'click', function(evt) {
		    console.log(['click on', id, m])
		    switch (operation) {
		    case "remove":
			removeCab(id, m)
			break;
		    }
		})
		google.maps.event.addListener(m, 'dragend', function(evt) {
		    console.log(['dragend', id, m, evt.latLng])
		    upsertCab(id, evt.latLng)
		})

		cabs[id] = m // store the marker
	    }
	});
	console.log(['add_cab', cab, url])
    }

    function removeCab(id, marker) {
	var url = [ rest, "/cabs/", id ].join('')
	console.log(['remove', url])
	marker.setMap(null)
	// In Chrome this somehow ends up with a OPTIONS http method when DELETE is set
	// $.ajax({
        //     url: url,
        //     type: "DELETE"
        //     contentType: 'application/json',
        //     always: function(result) {
	// 	marker.setMap(null)
        //     }
        // });

	// Here we implement a workaround with a new endpoint that supports POST
	jQuery.post([rest, "/delete/", id].join(''), "").always(function(e) {
	    if (e.status == 200) {
		marker.setMap(null)
	    }
	});
    }

    function removeAll() {
	var url = [ rest, "/cabs" ].join('')
	console.log(['remove_all', url])

	// In Chrome this somehow ends up with a OPTIONS http method, which is incorrect!
	// $.ajax({
        //     url: url,
        //     type: "DELETE",
        //     contentType: 'application/json',
        //     always: function(result) {
	// 	for (key in cabs) {
	// 	    cabs[key].setMap(null)
	// 	}
        //     }
        // });

	// Here we implement a workaround with a new endpoint that supports POST
	jQuery.post([rest, "/deleteAll"].join(''), "").always(function(e) {
	    if (e.status == 200) {
		for (key in cabs) {
		    cabs[key].setMap(null)
		}
	    }
	});
    }

    function drawQueryCircle(circle) {
	// clean the last query circle
	if (queryCircle != null) {
	    queryCircle.setMap(null)
	}
	var circOptions = {
	    strokeColor: '#FF0000',
	    strokeOpacity: 0.8,
	    strokeWeight: 2,
	    fillColor: '#FF0000',
	    fillOpacity: 0.35,
	    map: map,
	    center: circle.center,
	    radius: circle.radius
	};
	// Add the circle for this city to the map.
	queryCircle = new google.maps.Circle(circOptions);
	console.log(['drawCircle', queryCircle])
    }

    function getQueryCircle(loc) {
	var radius = $('#radius').val()
	radius = (radius.length == 0)? 100. : parseFloat(radius)

	var limit = $('#limit').val()
	limit = (limit.length == 0) ? 8 : parseInt(limit)
	return {
	    center:loc,
	    radius:radius,
	    limit:limit
	}
    }

    // c is the query circle
    function query(c) {
	var params = [
	    "latitude=" + c.center.lat(),
	    "longitude=" + c.center.lng(),
	    "radius=" + c.radius,
	    "limit=" + c.limit
	]
	var query = [ rest, "/cabs?", params.join('&')].join('')

	jQuery.get(query, function(result) {
	    console.log([ 'matches', result ]);

	    for (var i=0; i < result.length; i++) {
		if (cabs[result[i].id] !== undefined) {
		    setInRange(cabs[result[i].id])
		}
	    }
	});
	console.log([ query ])
    }

    function initialize() {
	if (navigator.geolocation) {
	    navigator.geolocation.getCurrentPosition(function (position) {
		currentLocation = new google.maps.LatLng(position.coords.latitude, position.coords.longitude)
		console.log(["got location", currentLocation])
		map.setCenter(currentLocation)
		homeMarker.setPosition(currentLocation)
	    }, function (error) {
		console.log(arguments);
	    });
	}

	var mapOptions = {
	    zoom: 15,
	    center: currentLocation
	}
	map = new google.maps.Map(document.getElementById('map-canvas'), mapOptions);

	homeMarker = new google.maps.Marker({
	    position: currentLocation,
	    map: map,
	    draggable: true,
	    title: 'Home'
	});
	google.maps.event.addListener(homeMarker, 'click', function() {
	    map.setCenter(homeMarker.getPosition());
	});
	google.maps.event.addListener(homeMarker, 'dragstart', function() {
	    if (queryCircle != null) {
		queryCircle.setMap(null)
	    }
	});

	google.maps.event.addListener(map, 'click', function(evt) {
	    switch (operation) {
	    case "add":
		upsertCab(new Date().getTime(), evt.latLng)
		break;
	    }
	});

    }
    google.maps.event.addDomListener(window, 'load', initialize);


    $('#add_cab').click(function(evt) {
	operation = operation == "add"? "done" : "add" // toggle
	switch (operation) {
	case "add":
	    $(evt.target).html('Done')
	    $('#remove_cab').html('Remove Cabs')
	    break;
	default:
	    $(evt.target).html('Add Cabs')
	    break;
	}
    })

    $('#remove_cab').click(function(evt) {
	operation = operation == "remove"? "done" : "remove" // toggle
	switch (operation) {
	case "remove":
	    $(evt.target).html('Done')
	    $('#add_cab').html('Add Cabs')
	    break;
	default:
	    $(evt.target).html('Remove Cabs')
	    break;
	}
    })

    $('#remove_all_cabs').click(function(evt) {
	operation = "remove_all"
	if (confirm("Really remove all cabs from the roads?")) {
	    removeAll()
	}
    })

    $('#query').click(function(evt) {
	operation = "query"
	c = getQueryCircle(homeMarker.getPosition())
	drawQueryCircle(c)
	resetCabs()
	query(c)
    })

})
