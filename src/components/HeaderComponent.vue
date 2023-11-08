<template>
    <nav class="navbar navbar-expand-lg" style="background-color: #2BBFFF; color: #003366; font-size: 18px;">
        <div class="container-fluid">
            <transition name="slide" appear>
                <!-- <a class="navbar-brand" href="/"><img src="https://onestepgps.com/images/onestepgps_color_logotype.svg" alt="logo" style="width: 200px; height: 30px;"></a> -->
                <router-link class="navbar-brand" to="/">
                    <img src="https://onestepgps.com/images/onestepgps_color_logotype.svg" alt="logo" style="width: 200px; height: 30px;">
                </router-link>
            </transition>
            <button class="navbar-toggler" type="button" data-bs-toggle="collapse" data-bs-target="#navbarSupportedContent" aria-controls="navbarSupportedContent" aria-expanded="false" aria-label="Toggle navigation">
                <span class="navbar-toggler-icon"></span>
            </button>
            <div class="collapse navbar-collapse" id="navbarSupportedContent">
                <ul class="navbar-nav ms-auto mb-2 mb-lg-0">
                    <transition name="fade" appear>
                        <li class="nav-item">
                            <router-link class="nav-link" to="/home">
                                <i class="fa-solid fa-house"></i>
                            </router-link>
                        </li>
                    </transition>
                    <transition name="fade" appear>
                        <li class="nav-item">
                            <a class="nav-link" aria-current="page" href="#" @click="handleLogout"><i class="fa-solid fa-right-from-bracket"></i></a>
                        </li>
                    </transition>
                </ul>
            </div>
        </div>
    </nav>
</template>
<script>
import AuthService from '@/services/AuthService';
export default{
    name: 'HeaderComponent',
    methods:{
        handleLogout(){
            AuthService.logout()
            this.$router.push('/');
        }
    }
}
</script>
<style>
.nav-link{
    letter-spacing: 2px;
}
.slide-enter-from {
  transform: translateX(-100%);
}

.slide-enter-active {
  transition: transform 0.5s ease-in-out;
}

.slide-enter-to {
  transform: translateX(0);
}
.fade-enter-from {
    opacity: 0;
}


.fade-enter-active {
    transition: opacity 1s ease-in-out;
}

.fade-enter-to {
    opacity: 1;
}
</style>