import { createApp } from 'vue'
import App from './App.vue'
import { createStore } from 'vuex';
import createPersistedState from 'vuex-persistedstate';
import router from './router/router';
import 'bootstrap/dist/css/bootstrap.css';
import 'bootstrap/dist/js/bootstrap.js';


function initialState() {
    return {
        userPreferences: null,
        devices: [],
        deviceIcons: {}
    };
}

const store = createStore({
    state: initialState,
    actions: {
        logout({ commit }) {
            commit('resetState');
        }
    },
    mutations:{
        setUserPreference(state, preference){
            state.userPreferences = preference
        },
        setDevices(state, devices){
            state.devices = devices
        },
        resetState(state) {
            Object.assign(state, initialState());
        }
    },
    getters:{
        getUserPreference(state){
            return state.userPreferences
        },
        getDevices(state){
            return state.devices
        }
    },
    plugins: [createPersistedState()] 
})
const app = createApp(App);
app.use(router);
app.use(store)
app.mount('#app');

export default store;
