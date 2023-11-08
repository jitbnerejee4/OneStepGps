<template>
  <div style="background-color: #FAFAFA;">
    <HeaderComponent></HeaderComponent>
    <div class="row">
      <div class="col col-md-12 col-sm-12" style="position: relative; display: flex; flex-direction: row-reverse; right: 10%; margin-top: 2%;">
        <DropDownComponent @selected-view="handleView"></DropDownComponent>
        <SortComponent @selected-sort="handleSort"></SortComponent>
        <PreferenceComponent></PreferenceComponent>
      </div>
    </div>
    <div class="row">
      <div class="col col-md-12 col-sm-12">
        <CardComponent v-if="selectedView == 'list'" @devices-from-card = "receivedDevices" :selectedSort = "selectedSort"></CardComponent>
        <FlipperComponent v-if="selectedView == 'card'" :selectedSort = "selectedSort"></FlipperComponent>
      </div>
    </div>
  </div>
</template>

<script>
import CardComponent from './CardComponent.vue'
import HeaderComponent from './HeaderComponent.vue'
import FlipperComponent from './FlipperComponent.vue'
import DropDownComponent from './DropDownComponent.vue'
import SortComponent from './SortComponent.vue'
import PreferenceComponent from './PreferenceComponent.vue'
export default{
  name: 'HomeComponent',
  components:{
    CardComponent, 
    HeaderComponent,
    FlipperComponent,
    DropDownComponent,
    SortComponent,
    PreferenceComponent
  },
  
  mounted() {
    if(this.$store.getters.getUserPreference){
      this.storedPreferences = this.$store.getters.getUserPreference
      this.selectedView = this.storedPreferences.preference.viewPreference
      this.selectedSort = this.storedPreferences.preference.sortPreference
    }
  },

  data(){
    return{
      selectedView: 'list',
      selectedSort: 'unsorted',
      devices: [],
      storedPreferences: null
    }
  },
  computed: {
    userPreferences() {
      return this.$store.getters.getUserPreference;
    }
  },

  watch: {
    userPreferences: {
      deep: true,
      handler(newPreferences) {
        if (newPreferences && newPreferences.preference) {
          this.selectedView = newPreferences.preference.viewPreference;
          this.selectedSort = newPreferences.preference.sortPreference;
        }
      }
    }
  },
  methods:{
    handleView(view){
      this.selectedView = view
    },
    handleSort(sortOption){
      this.selectedSort = sortOption
    },
    receivedDevices(data){
      this.devices = data
    }
  }
}

</script>

<style>

</style>