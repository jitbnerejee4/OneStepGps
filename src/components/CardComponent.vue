<template>
    <div>
        <div class="data-container">
        <div class="card mb-4 top-card info-card">
          <div class="card-body">
            <div class="row">
              <div class="col col-md-4 col-sm-4">
                <div class="card-text" style="font-weight: bolder;">Map</div>
              </div>
              <div class="col col-md-2 col-sm-2">
                <div class="card-text" style="font-weight: bolder;">Device Name</div>
              </div>
              <div class="col col-md-2 col-sm-2">
                <div class="card-text" style="font-weight: bolder;">Active State</div>
              </div>
              <div class="col col-md-2 col-sm-2">
                <div class="card-text" style="font-weight: bolder;">Current Location</div>
              </div>
              <div class="col col-md-2 col-sm-2">
                <div class="card-text" style="font-weight: bolder;">Speed</div>
              </div>
            </div>
          </div>
        </div>
        <div class="card mb-4 info-card" v-for="device in devices" :key="device.DeviceId">
          <div class="card-body">
            <div class="row">
              <div class="col col-md-4 col-sm-4">
                <GoogleMap :api-key="googleMapsApiKey" style="width: 100%; height: 200px" :center="{ lat: device.latest_device_point.lat, lng: device.latest_device_point.lng }" :zoom="15">
                  <Marker @click="showMap(device)" :options="{ position: { lat: device.latest_device_point.lat, lng: device.latest_device_point.lng }, icon: faIcon  }" />
                </GoogleMap>
              </div>
              <div class="col col-md-2 col-sm-2" style="display: flex; align-items: center; justify-content: center;">
                <div class="card-text">{{ device.display_name }}</div>
              </div>
              <div class="col col-md-2 col-sm-2" style="display: flex; align-items: center; justify-content: center;">
                <div class="card-text">{{device.active_state}}</div>
              </div>
              <div class="col col-md-2 col-sm-2" style="display: flex; align-items: center; justify-content: center;">
                <div class="card-text">{{ device.address }}</div>
              </div>
              <div class="col col-md-2 col-sm-2" style="display: flex; align-items: center; justify-content: center;">
                <div class="card-text"><div :style="{ 'background-color': device.latest_device_point.speed === 0 ? 'red' : 'orange', padding: '10px', 'border-radius': '20px', letterSpacing: '2px' }">{{ formattedSpeed(device.latest_device_point.speed) }}</div></div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
    <div class="modal fade modal-xl" id="exampleModal" tabindex="-1" aria-labelledby="exampleModalLabel" aria-hidden="true">
        <div class="modal-dialog">
            <div class="modal-content">
                <div class="modal-header">
                    <h1 class="modal-title fs-5" id="exampleModalLabel">{{ selectedDeviceName }}, {{ selectedDeviceAddress }}</h1>
                    <button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close"></button>
                </div>
                <div class="modal-body">
                  <div id="map" style="height: 600px;">
                  </div>
                </div>
            </div>
        </div>
    </div>
      
</template>


<script>
/*global google*/ 
import axios from 'axios'
import { GoogleMap, Marker } from "vue3-google-map";
import * as bootstrap from 'bootstrap';

export default {
  name: 'CardComponent',
  components: {
    GoogleMap,
    Marker
  },
  props:['selectedSort'],
  emits:['devices-from-card'],
  data(){
    const faIcon = `
      data:image/svg+xml;utf-8,
      <svg xmlns="http://www.w3.org/2000/svg" height="2em" viewBox="0 0 448 512"><!--! Font Awesome Free 6.4.2 by @fontawesome - https://fontawesome.com License - https://fontawesome.com/license (Commercial License) Copyright 2023 Fonticons, Inc. --><path d="M429.6 92.1c4.9-11.9 2.1-25.6-7-34.7s-22.8-11.9-34.7-7l-352 144c-14.2 5.8-22.2 20.8-19.3 35.8s16.1 25.8 31.4 25.8H224V432c0 15.3 10.8 28.4 25.8 31.4s30-5.1 35.8-19.3l144-352z"/></svg>
      `;
    return{
      devices: [],
      unfilteredDevices: [],
      error: null,
      interval: 100,
      maxDuration: 0,
      faIcon: faIcon,
      apiCallInterval: null,
      latitude: 0,
      longitude: 0,
      selectedDeviceName: null,
      selectedDeviceAddress: null,
      unsortedDevices: [],
      initialSort: null,
      storedPreferences: null,
      googleMapsApiKey: process.env.VUE_APP_GOOGLE_MAPS_API_KEY
    }
  },
  mounted(){
    this.fetchDevices();
    this.apiCallInterval = setInterval(this.fetchDevices, 50000);
  },
  computed: {
    userPreferences() {
      return this.$store.getters.getUserPreference;
    }
  },
  watch: {
    // devices: function(newVal) {
    //   this.maxDuration = this.interval * newVal.length;
    // },
    userPreferences: {
      deep: true,
      handler(newPreferences) {
        if (newPreferences && newPreferences.preference) {
          // this.initialSort = newPreferences.preference.sortPreference;
          this.filterDevices(newPreferences)
          this.sortDevices(newPreferences.preference.sortPreference)
        }
      }
    },
    selectedSort(newSort) {
      this.sortDevices(newSort)
    },
  },  
  
  methods:{
    fetchDevices(){
      axios.get('http://localhost:8000', {
        headers: {
          'Authorization': `Bearer ${localStorage.getItem('token')}`
        }
      })
      .then(response => {
        this.unsortedDevices = response.data;
        this.unfilteredDevices = [...this.unsortedDevices]
        this.$store.commit('setDevices', {devices: this.unfilteredDevices})
        if(this.$store.getters.getUserPreference){
          this.storedPreferences = this.$store.getters.getUserPreference
          this.filterDevices(this.storedPreferences)
          this.sortDevices(this.storedPreferences.preference.sortPreference)
        }
        this.sendDatatoHome()
      })
      .catch(err => {
        this.error = "An error occurred: " + err;
      });
    },
    filterDevices(newPreferences){
      this.devices = []
      this.unsortedDevices = []
      this.unfilteredDevices.forEach(device => {
        if(newPreferences.preference.visibilityPreference[device.device_id] == true){
          this.devices.push(device)
          this.unsortedDevices.push(device)
        }
        this.$store.commit('setDevices', {devices: this.unfilteredDevices})
      });
      this.$nextTick(() => {
        for (let i = 0; i < this.devices.length + 1; i++) {
          setTimeout(() => {
            const cardElements = document.querySelectorAll('.card');
            if (cardElements[i]) {
              cardElements[i].classList.add('enter-to');
            }
          }, i * this.interval);
        }
      });
    },
    sortDevices(newSort){
      if (newSort === 'name') {
        this.devices.sort((a, b) => a.display_name.localeCompare(b.display_name));
      } else if (newSort === 'speed') {
        this.devices.sort((b, a) => a.latest_device_point.speed - b.latest_device_point.speed);
      } else {
        this.devices = this.unsortedDevices
      }
    },
    showMap(device){
        this.latitude = device.latest_device_point.lat
        this.longitude = device.latest_device_point.lng
        this.selectedDeviceName = device.display_name
        this.selectedDeviceAddress = device.address

        // show/hide modal
        var modalElement = document.getElementById('exampleModal');
        var modalInstance = new bootstrap.Modal(modalElement);
        modalInstance.show();

      // creating map 
      modalElement.addEventListener('shown.bs.modal', () => {
        this.initializeMap();
      });

    },
    initializeMap(){
      var mapDiv = document.getElementById('map');
        var mapOptions = {
          center: new google.maps.LatLng(this.latitude, this.longitude),
          zoom: 15,
          // mapTypeId: google.maps.MapTypeId.ROADMAP
        };
        var map = new google.maps.Map(mapDiv, mapOptions);
        var markerOptions = {
          position: new google.maps.LatLng(this.latitude, this.longitude),
          map: map,
          icon: this.faIcon
        };
        new google.maps.Marker(markerOptions);

    },
    toPacificTime(dateTimeString) {
        const dateObj = new Date(dateTimeString);
        return dateObj.toLocaleString("en-US", {
        timeZone: "America/Los_Angeles",
        year: 'numeric',
        month: '2-digit',
        day: '2-digit',
        hour: '2-digit',
        minute: '2-digit',
        second: '2-digit'
        });
    },
    formattedSpeed(speed) {
      if(speed === 0) {
        return "Stopped";
      }
      return `${Math.floor(speed)} km/h`;
    },
    sendDatatoHome(){
      this.$emit('devices-from-card', this.devices)
    }


  },
  beforeUnmount() {
    clearInterval(this.updateInterval);
  }
}
</script>
<style>

.data-container{
    text-align: center;
    padding-bottom: 50px;
    padding-top: 100px;

}
.top-card {
    margin-top: 10px;
}
.info-card {
  margin-left: 10%;
  margin-right: 10%;
  transform: translateY(100%); 
  opacity: 0;
  transition: transform 1s, opacity 1s; 
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);

}

.card.enter-to {
  transform: translateY(0);
  opacity: 1;
}
.card-text{
    font-family: 'Gentium Plus', serif;
    font-family: 'Montserrat', sans-serif;
    font-family: 'Roboto', sans-serif;
    font-size: 18px;
    
}

@media (max-width: 991px) {
  .info-card {
    margin-left: 5%;
    margin-right: 5%;
  }
}

@media (max-width: 767px) {
  .info-card {
    margin-left: 2%;
    margin-right: 2%;
  }
  .card-text {
    font-size: 16px;
  }
  .modal-dialog {
    width: 90%;
    max-width: none;
  }
  #map {
    width: calc(100% - 30px); 
    height: 300px; 
  }
}
@media (max-width: 575px) {
  .card-text {
    font-size: 14px;
  }
  
  .modal-body {
    padding: 15px; 
  }
  
  #map {
    height: 200px; 
  }
}
</style>