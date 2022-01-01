Vue.use(GmapVue, {
  load: {
    key: "AIzaSyBKfytmIV5szDPgIaTzMMd7xLKM2Aa1ivc",
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
