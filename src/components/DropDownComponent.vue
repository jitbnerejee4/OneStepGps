<template>
    <div class="dropdown">
        <a class="btn btn-secondary dropdown-toggle" href="#" role="button" data-bs-toggle="dropdown" aria-expanded="false">View</a>

        <ul class="dropdown-menu">
                <li @click="selectedView('list')" class="dropdown-item">
                    <div class="row align-items-center">
                        <div class="col col-md-2">
                            <span v-if="currentView == 'list'"><i class="fa-solid fa-check"></i></span>
                        </div>
                        <div class="col col-md-10">List View</div>
                    </div>
                </li>
                <li @click="selectedView('card')" class="dropdown-item">
                    <div class="row align-items-center">
                        <div class="col col-md-2">
                            <span v-if="currentView=='card'"><i class="fa-solid fa-check"></i></span>
                        </div>
                        <div class="col col-md-10">Card View</div>
                    </div>
                </li>
        </ul>
    </div>
</template>

<script>
export default{
    name: 'DropDownComponent',
    data(){
        return{
            currentView: 'list',
            storedPreferences: null
        }
    },
    mounted(){
        if(this.$store.getters.getUserPreference){
            this.storedPreferences = this.$store.getters.getUserPreference
            this.currentView = this.storedPreferences.preference.viewPreference
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
          this.currentView = newPreferences.preference.viewPreference;
        }
      }
    }
  },
    methods:{
        selectedView(view){
            this.currentView = view
            this.$emit('selected-view', view)
        }
    }
}

</script>
<style>

.dropdown-item:hover {
    background-color: #f5f5f5; 
    cursor: pointer;
}
</style>