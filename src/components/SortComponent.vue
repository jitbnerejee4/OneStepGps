<template>
    <div class="dropdown" style="margin-right: 20px;">
    <a class="btn btn-secondary dropdown-toggle" href="#" role="button" data-bs-toggle="dropdown" aria-expanded="false">Sort</a>

    <ul class="dropdown-menu">
            <li class="dropdown-item" @click="selectedSort('unsorted')">
                <div class="row align-items-center">
                    <div class="col col-md-2">
                        <span v-if="selectedsort == 'unsorted'"><i class="fa-solid fa-check"></i></span>
                    </div>
                    <div class="col col-md-10">Default</div>
                </div>
            </li>
            <li class="dropdown-item" @click="selectedSort('name')">
                <div class="row align-items-center">
                    <div class="col col-md-2">
                        <span v-if="selectedsort == 'name'"><i class="fa-solid fa-check"></i></span>
                    </div>
                    <div class="col col-md-10">By name</div>
                </div>
            </li>
            <li class="dropdown-item" @click="selectedSort('speed')">
                <div class="row align-items-center">
                    <div class="col col-md-2">
                        <span v-if="selectedsort == 'speed'"><i class="fa-solid fa-check"></i></span>
                    </div>
                    <div class="col col-md-10">By speed</div>
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
            selectedsort: 'unsorted'
        }
    },
    mounted(){
        if(this.$store.getters.getUserPreference){
            this.storedPreferences = this.$store.getters.getUserPreference
            this.selectedsort = this.storedPreferences.preference.sortPreference
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
                this.selectedsort = newPreferences.preference.sortPreference;
            }
        }
        }
    },
    methods:{
        selectedSort(sortOption){
            this.selectedsort = sortOption
            this.$emit('selected-sort', this.selectedsort)
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