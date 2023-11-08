<template>
    <button type="button" class="btn btn-secondary" @click="showPreferences" data-bs-toggle="modal" data-bs-target="#preferencesModal" style="margin-right: 20px;">
        Edit Preferences
    </button>
    <div class="modal fade" id="preferencesModal" tabindex="-1" aria-labelledby="preferencesModalLabel" aria-hidden="true">
        <div class="modal-dialog modal-dialog-centered modal-lg">
            <div class="modal-content">
                <div class="modal-header">
                    <h5 class="modal-title" id="preferencesModalLabel">User Preferences</h5>
                    <button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close"></button>
                </div>
                <div class="modal-body">
                    <form id="preferencesForm" @submit.prevent="handlePreference">
                        <!-- View Preference -->
                            <div class="mb-3">
                                <label for="viewPreference" class="form-label" style="display: flex; justify-content: start; font-weight: bold;">View Preference</label>
                                <select class="form-select" id="viewPreference" v-model="preference.viewPreference">
                                    <option value="list">List View</option>
                                    <option value="card">Card View</option>
                                </select>
                            </div>
                        <!-- Sort Preference -->
                            <div class="mb-3">
                                <label class="form-label" style="display: flex; justify-content: start; font-weight: bold;">Sort Preference</label>
                                <div style="display: flex; justify-content: start;">
                                    <div class="form-check">
                                        <input class="form-check-input" type="radio" value="unsorted" name="flexRadioDefault" id="flexRadioDefault2" checked v-model="preference.sortPreference">
                                        <label class="form-check-label me-3" for="flexRadioDefault2">
                                            Default
                                        </label>
                                    </div>
                                    <div class="form-check">
                                        <input class="form-check-input" type="radio" value="name" name="flexRadioDefault" id="flexRadioDefault3" v-model="preference.sortPreference">
                                        <label class="form-check-label me-3" for="flexRadioDefault3">
                                            By Name
                                        </label>
                                    </div>
                                    <div class="form-check">
                                        <input class="form-check-input" type="radio" value="speed" name="flexRadioDefault" id="flexRadioDefault1" v-model="preference.sortPreference">
                                        <label class="form-check-label me-3" for="flexRadioDefault1">
                                            By Speed
                                        </label>
                                    </div>
                                </div>
                            </div>
                        <!-- Device Visibility -->
                            <div class="mb-3">
                                <label class="form-label" style="display: flex; justify-content: start; font-weight: bold;">Device Visibility</label>
                                <!-- Dynamically add checkboxes for each device -->
                                <div class="visibility-container">
                                    <div class="form-check me-3" v-for="device in devices" :key="device.DeviceId">
                                        <input class="form-check-input" type="checkbox" value="" :id="device.device_id" v-model="preference.visibilityPreference[device.device_id]">
                                        <label class="form-check-label" :for="device.device_id">
                                            {{ device.display_name }}
                                        </label>
                                    </div>
                                </div>
                            </div>
                        <div class="modal-footer">
                            <button type="button" class="btn btn-secondary" data-bs-dismiss="modal">Close</button>
                            <button type="submit" class="btn btn-primary" data-bs-dismiss="modal">Save Preferences</button>
                        </div>
                    </form>
                </div>
            </div>
        </div>
    </div>

</template>
<script>
import AuthService from '@/services/AuthService';
export default{

    data(){
        return{
            devices: this.allDevices,
            preference:{
                viewPreference: 'list',
                sortPreference: 'unsorted',
                visibilityPreference: {}
            },
            storedPreferences: null
        }
    },
    methods:{
        async handlePreference(){
            const response = await AuthService.savePreference(this.preference)
            if (response.status == 200) {
                this.$store.commit('setUserPreference', {preference: response.data.userPreferences})
            }
        },
        showPreferences(){
            this.devices = this.$store.getters.getDevices.devices

            if(this.$store.getters.getUserPreference){
                this.storedPreferences = this.$store.getters.getUserPreference
                this.preference.sortPreference = this.storedPreferences.preference.sortPreference
                this.preference.viewPreference = this.storedPreferences.preference.viewPreference
                this.devices.forEach(device=>{
                    const isVisible = this.storedPreferences.preference.visibilityPreference[device.device_id]
                    this.preference.visibilityPreference[device.device_id] = isVisible;
                })
            }
        }
    }
}
</script>
<style>
.visibility-container{
    position: relative; 
    display: flex;
    flex-wrap: wrap;
    justify-content: flex-start;

}
</style>