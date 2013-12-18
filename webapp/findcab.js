var currentLocation = new google.maps.LatLng(37.79, -122.393);
var map = null;
var marker = null;
var host = null;
var operation = "view";
var rest = null

$(function(){

    host = location.hostname
    rest = location.protocol + "//" + host + ":8080"

    function initialize() {
	if (navigator.geolocation) {
	    navigator.geolocation.getCurrentPosition(function (position) {
		currentLocation = new google.maps.LatLng(position.coords.latitude, position.coords.longitude)
		console.log(["got location", currentLocation])
		map.setCenter(currentLocation)
		marker.setPosition(currentLocation)
	    }, function (error) {
		console.log(arguments);
	    });
	}

	var mapOptions = {
	    zoom: 15,
	    center: currentLocation
	}
	map = new google.maps.Map(document.getElementById('map-canvas'), mapOptions);

	var marker = new google.maps.Marker({
	    position: currentLocation,
	    map: map,
	    title: 'Home'
	});
    }
    google.maps.event.addDomListener(window, 'load', initialize);


    $('#add_cab').click(function(evt) {
	operation = "add"
    })

    $('#remove_cab').click(function(evt) {
	operation = "remove"
    })

    $('#remove_all_cabs').click(function(evt) {
	operation = "remove_all"
    })

    $('#query').click(function(evt) {
	operation = "query"
	var radius = $('#radius').val()
	if (radius.length == 0) radius = 1000.

	var limit = $('#limit').val()
	if (limit.length == 0) limit = 8

	var params = [
	    "latitude=" + currentLocation.lat(),
	    "longitude=" + currentLocation.lng(),
	    "radius=" + radius,
	    "limit=" + limit
	]
	var query = [ rest, "/cabs?", params.join('&')].join('')

	jQuery.get(query, function(result) {
	    console.log([ 'matches', result ]);
	});


	console.log([ query ])


    })

})
