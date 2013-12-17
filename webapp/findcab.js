

function initialize() {
  var myLatlng = new google.maps.LatLng(37.79, -122.393);
  var mapOptions = {
    zoom: 15,
    center: myLatlng
  }
  var map = new google.maps.Map(document.getElementById('map-canvas'), mapOptions);

  var marker = new google.maps.Marker({
      position: myLatlng,
      map: map,
      title: 'Uber'
  });
}

google.maps.event.addDomListener(window, 'load', initialize);
