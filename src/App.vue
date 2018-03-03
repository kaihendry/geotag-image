<template>
  <div id="app">
<form action="/upload" enctype="multipart/form-data" method="post">
<input type="file" required name="jpeg" />
<br>
<label>Lat: <input type="number" step="any" name=lat required v-model.number.lazy="mapCenter.lat" /></label>
<label>Lng: <input type="number" step="any" name=lng required v-model.number.lazy="mapCenter.lng" /></label>
<input type="submit" value="Report" />
</form>
</div>
</template>

<script>
export default {
  name: 'app',
  data: function() {
    return {
      mapCenter: {
        lat: 1.38,
        lng: 103.8,
      }
    };
  },
  created: function () {
    console.log('a is: ' + this.mapCenter)
    if(navigator.geolocation) {
      navigator.geolocation.getCurrentPosition((position) => {
        console.log(position.coords.latitude, position.coords.longitude)
        this.mapCenter.lat = position.coords.latitude
        this.mapCenter.lng = position.coords.longitude
      })
    }
  },
  methods: {
    updateCenter(newCenter) {
      this.mapCenter = {
        lat: newCenter.lat(),
        lng: newCenter.lng(),
      }
    }
  }
}
</script>

<style>
#app {
  font-family: 'Avenir', Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
}
</style>
