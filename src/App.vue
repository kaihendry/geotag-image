<template>
<div id="app">
<form action="/upload" enctype="multipart/form-data" method="post">
<input type="file" required name="jpeg" />
<br>

<label>
Lat: <input type="number"
  name="lat"
  v-model.number.lazy="reportedMapCenter.lat"
  @change="sync"
  step="any" />
</label>

<label>
Lng: <input type="number"
  name="lng"
  v-model.number.lazy="reportedMapCenter.lng"
  @change="sync"
  step="any" />
</label>

<input type="submit" />
<p>Please report issues to <a href="https://github.com/kaihendry/geotag-image/issues">Geotag Github</a></p>
<p><a target="_blank" :href="`http://maps.google.com/maps?q=${mapCenter.lat},${mapCenter.lng}`">Open maps</a></p>
</form>


<gmap-map :center="mapCenter" :zoom="12"
ref="map"
@center_changed="updateCenter($refs.map.$mapObject.getCenter())"
:options="{disableDefaultUI : true}"
@dragend="sync"
@zoom_changed="sync"
class="map-container">

<div slot="visible">
<div style="height: 100%; width: 100%; display: block; position: absolute; ">
    <img
        id="crosshairs"
        src="./assets/crosshairs.gif"
        style="
          width: 19px;
          height: 19px;
          border: 0;
          zIndex: 100;
        "
    />
  </div>
</div>
</gmap-map>

</div>
</template>

<script>
export default {
  name: 'app',
  data: function() {
return {
      reportedMapCenter: {
        lat: 1.38,
        lng: 103.8,
      },
  mapCenter: {}
  };
},
  created: function () {
    if(navigator.geolocation) {
      navigator.geolocation.getCurrentPosition((position) => {
        // console.log(position.coords.latitude, position.coords.longitude)
        this.reportedMapCenter.lat = position.coords.latitude
        this.reportedMapCenter.lng = position.coords.longitude
        this.sync()
      })
    }
  },
  methods: {
    updateCenter(latLng) {
      this.reportedMapCenter = {
        lat: latLng.lat(),
        lng: latLng.lng()
      }
    },
    sync () {
      this.mapCenter = this.reportedMapCenter
    }
  }
}
</script>

<style>
html { width: 100%; height: 100% }
body { height: 100%; margin: 0px; padding: 0px }
#app {
  font-family: 'Avenir', Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  height: 100%; margin: 0px; padding: 0px
}
.map-container {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
}
form { display: flex; justify-content: space-evenly }
</style>
