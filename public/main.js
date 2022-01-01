Vue.use(GmapVue, {
  load: {
    key: "AIzaSyD4VHBovJ2dnHSZpS-Y46hheA_JL6mtwZI",
    libraries: "places", // This is required if you use the Autocomplete plugin
  },
});

new Vue({
  el: "#app",
  data: function () {
    return {
      reportedMapCenter: {
        lat: 1.38,
        lng: 103.8,
      },
      mapCenter: {},
    };
  },
  created: function () {
    if (navigator.geolocation) {
      navigator.geolocation.getCurrentPosition((position) => {
        // console.log(position.coords.latitude, position.coords.longitude)
        this.reportedMapCenter.lat = position.coords.latitude;
        this.reportedMapCenter.lng = position.coords.longitude;
        this.sync();
      });
    }
  },
  methods: {
    updateCenter(latLng) {
      this.reportedMapCenter = {
        lat: latLng.lat(),
        lng: latLng.lng(),
      };
    },
    sync() {
      this.mapCenter = this.reportedMapCenter;
    },
  },
});
