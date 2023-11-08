<template>
    <div class="flip-container">
        <div class="flip-card" v-for="device in devices"  :key="device.DeviceId">
            <input type="file" :ref="`fileInput-${device.device_id}`" style="display: none" @change="handleFileUpload" accept="image/*" />
            <div class="flip-card-inner" :class="{ 'flipped': device.isFlipped }">
                <div class="flip-card-front" @click="flipCard(device)">
                    <div class="row">
                        <div class="col col-md-12 col-sm-12 image-container">
                            <img :src="getIcons(device)" alt="icon" @click.stop="triggerFileInput(device)" style="border: 2px solid green; border-radius: 50%; width: 200px; height: 200px;">
                            <div class="online-icon" title="Device is online"  v-if="device.active_state == 'active'"><i class="fa-solid fa-circle" style="color: #3cf609; font-size: 25px;"></i></div>
                        </div>
                        <div class="col col-md-12 col-sm-12">
                            <h2>{{ device.display_name }}</h2>
                            <p>{{ device.address }}</p>
                            <div class="speed-div">
                                <div  :style="{ 'background-color': device.latest_device_point.speed === 0 ? 'red' : 'orange', padding: '10px', 'border-radius': '10px', letterSpacing: '2px' }">{{ formattedSpeed(device.latest_device_point.speed) }}</div>
                            </div>
                        </div>
                    </div>
                </div>
                <div class="flip-card-back" >
                    <div class="card border-dark" style="min-width: 100%;">
                        <GoogleMap :api-key="googleMapsApiKey" style="width: 100%; height: 500px" :center="{ lat: device.latest_device_point.lat, lng: device.latest_device_point.lng }" :zoom="15">
                            <Marker @click="flipCard(device)" :options="{ position: { lat: device.latest_device_point.lat, lng: device.latest_device_point.lng }, icon: faIcon  }" />
                        </GoogleMap>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>
<script>
import { GoogleMap, Marker } from "vue3-google-map";
import axios from "axios";
import AuthService from "@/services/AuthService";

export default{
    name: 'FlipperComponent',
    props: ['selectedSort'],
    components: {
        GoogleMap,
        Marker
    },
    data() {
        return {
            devices: [],
            unfilteredDevices: [],
            interval: 100,
            visibleDevices: [],
            icons: {},
            storedPreferences: null,
            selectedDevice: null,
            apiCallInterval: null,
            googleMapsApiKey: process.env.VUE_APP_GOOGLE_MAPS_API_KEY,
            faIcon: `
                data:image/svg+xml;utf-8,
                <svg xmlns="http://www.w3.org/2000/svg" height="2em" viewBox="0 0 448 512"><!--! Font Awesome Free 6.4.2 by @fontawesome - https://fontawesome.com License - https://fontawesome.com/license (Commercial License) Copyright 2023 Fonticons, Inc. --><path d="M429.6 92.1c4.9-11.9 2.1-25.6-7-34.7s-22.8-11.9-34.7-7l-352 144c-14.2 5.8-22.2 20.8-19.3 35.8s16.1 25.8 31.4 25.8H224V432c0 15.3 10.8 28.4 25.8 31.4s30-5.1 35.8-19.3l144-352z"/></svg>
                `
        }
    },
    mounted(){
        this.fetchDevices()
        this.apiCallInterval = setInterval(this.fetchDevices, 50000);

    },
    computed: {
        userPreferences() {
            return this.$store.getters.getUserPreference;
        }
    },
    watch:{
        selectedSort(newSort) {
            this.sortDevices(newSort)
        },
        userPreferences: {
        deep: true,
            handler(newPreferences) {
                if (newPreferences && newPreferences.preference) {
                    this.filterDevices(newPreferences)
                    this.sortDevices(newPreferences.preference.sortPreference)
                    this.setIcons()
                }
            }
        }
    },
    methods: {
        flipCard(device) {
            device.isFlipped = !device.isFlipped;
        },
        formattedSpeed(speed){
            if(speed == 0){
                return 'Stopped'
            }
            return `${Math.floor(speed)} km/h`;
        },
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
                    this.setIcons()
                    this.filterDevices(this.storedPreferences)
                    this.sortDevices(this.storedPreferences.preference.sortPreference)
                }
            })
            .catch(err => {
                this.error = "An error occurred: " + err;
            });
        },
        setIcons(){
            const allIcons = this.$store.getters.getUserPreference.preference.deviceIcons
            allIcons.forEach(icon=>{
                this.icons[icon.name] = icon.iconUrl
            })
        },
        getIcons(device){
            return this.icons[device.display_name]
        },
        filterDevices(newPreferences){
            this.devices = []
            this.unsortedDevices = []
            this.unfilteredDevices.forEach(device => {
                if(newPreferences.preference.visibilityPreference[device.device_id] == true){
                    this.devices.push(device)
                    this.unsortedDevices.push(device)
                }
            });
            this.$nextTick(() => {
                    for (let i = 0; i < this.unfilteredDevices.length + 1; i++) {
                        setTimeout(() => {
                        const cardElements = document.querySelectorAll('.flip-card');
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
        triggerFileInput(device) {
            this.$refs[`fileInput-${device.device_id}`][0].click();
            this.selectedDevice = device.display_name
        },

        async handleFileUpload(event) {
            const file = event.target.files[0];
            if (file) {
                let formData = new FormData();
                formData.append('image', file); 
                formData.append('device', this.selectedDevice)
                const response = await AuthService.uploadImage(formData)
                if(response.status == 200){
                    this.$store.commit('setUserPreference', {preference: response.data.userPreferences})
                    this.setIcons()
                }
            }
        },
    }
}
</script>
<style scoped>
.flip-container {
    position: relative; 
    display: flex;
    flex-wrap: wrap;
    justify-content: flex-start;
    padding: 20px;       
    /* border: 1px solid #ccc; 
    overflow: hidden;    */
}

.flip-card {
  perspective: 1000px;
  /* width: calc(33.333% - 20px); */
  height: 500px;
  width: 100%; 
  margin: 10px 0;
  /* margin: 10px; */
  transition: transform 0.6s, opacity 0.6s;
  transform: translateY(100px);
  opacity: 0;
}
.flip-card.enter-to {
  transform: translateY(0);
  opacity: 1;
}


@media (min-width: 768px) {
    .flip-card {
        width: calc(50% - 20px); /* 2 cards per row */
        margin: 10px; /* Add left/right margin back */
    }
}


@media (min-width: 992px) {
    .flip-card {
        width: calc(33.333% - 20px); /* 3 cards per row */
    }
}

.flip-card-inner {
  position: relative;
  width: 100%;
  height: 100%;
  text-align: center;
  transition: transform 0.6s;
  transform-style: preserve-3d;
  box-shadow: 0 4px 8px 0 rgba(0,0,0,0.2);
}

.flip-card-front, .flip-card-back {
  position: absolute;
  width: 100%;
  height: 100%;
  backface-visibility: hidden;
  display: flex;
  align-items: center;
  justify-content: center;
  border: 1px solid #ccc; 
}

.flip-card-front {
  background-color: #fff;
}

.flip-card-back {
  background-color: #eee; 
  transform: rotateY(180deg);
}

.flipped {
  transform: rotateY(180deg);
}

.speed-div {
    display: flex;
    justify-content: center; 
    align-items: center;     
    width: 100%;
}

.image-container {
    position: relative;
    text-align: center;
    justify-content: center;
    align-items: center;
    display: flex;
    margin-bottom: 20px;
    
}

.online-icon {
    position: absolute; 
    top: 93%;

}
</style>