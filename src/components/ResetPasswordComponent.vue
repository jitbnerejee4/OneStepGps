<template>
    <div class=" reset-container d-flex justify-content-center align-items-center gradient-background" style="height: 100vh;">
        <div v-if="loading" class="loading-overlay">
            <img style="width: 200px; height: 200px; display: flex; justify-content: center;"  src="https://media0.giphy.com/media/QODU6spbkmhzg14hLx/giphy.gif?cid=ecf05e4742053fkspfuyavy4tabyqjg33op18hul4n9hzfxy&ep=v1_stickers_search&rid=giphy.gif&ct=s" alt="loading">
        </div>
        <div class="reset-password-form p-4 shadow col-xs-12 col-sm-10 col-md-8 col-lg-6 col-xl-4">
        <h2 class="text-center mb-4">Reset Password</h2>
        <form @submit.prevent="resetPassword">
            <div class="form-group mb-3">
                <label for="newPassword">New Password</label>
                <input :type="showPassword ? 'text' : 'password'" class="form-control" id="newPassword" v-model="newPassword" required>
                <span class="field-icon" v-if="!showPassword" @click="togglePasswordVisibility"><i class="fa-solid fa-eye-slash"></i></span>
                <span class="field-icon" v-if="showPassword" @click="togglePasswordVisibility"><i class="fa-solid fa-eye"></i></span>
                <p class="text-dark" v-if="!passwordValid">Password must be at least 6 characters long and include at least 1 number and 1 special character.</p>

            </div>
            <div class="form-group mb-4">
                <label for="confirmPassword">Confirm Password</label>
                <input type="password" class="form-control" id="confirmPassword" v-model="confirmPassword" required>
            </div>
            <button type="submit" class="btn btn-primary w-100">Reset Password</button>
        </form>
        </div>
  </div>
</template>
<script>
import AuthService from '@/services/AuthService';
import { toast } from 'vue3-toastify';
import 'vue3-toastify/dist/index.css';

    export default{
        mounted(){
            this.handleReset()
        },
        data() {
            return {
                newPassword: '',
                confirmPassword: '',
                showPassword: false,
                token: '',
                passwordValid: false,
                loading: false
                
            }
        },
        watch: {
            'newPassword': function() {
                this.validatePassword();
            }
        },
        methods:{
            handleReset(){
                this.token = this.$route.query.token;
            },
            validatePassword(){
                const regex = /^(?=.*[0-9])(?=.*[!@#$%^&*])[a-zA-Z0-9!@#$%^&*]{6,}$/;
                this.passwordValid = regex.test(this.newPassword);
            },
            async resetPassword() {
                if (this.newPassword === this.confirmPassword && this.passwordValid == true) {
                    try{
                        const response = await AuthService.resetPassword(this.token, this.newPassword)
                        if(response.status == 200){
                            this.loading = true
                            toast.success("Password reset successful" + " !", {
                                position: toast.POSITION.TOP_RIGHT,
                            });
                            setTimeout(() => {
                                this.loading = false
                                this.$router.push('/'); 
                            }, 3000); 
                        }else if(response.status == 401){
                            toast.error("Token expired!", {
                                position: toast.POSITION.TOP_RIGHT,
                            });
                        }else{
                            toast.error("Server error! Please try again later.", {
                                position: toast.POSITION.TOP_RIGHT,
                            });
                        }
                        
                    }catch(err){
                        toast.error("There is an error!", {
                            position: toast.POSITION.TOP_RIGHT,
                        });
                    }
                }   else if(this.newPassword != this.confirmPassword && this.passwordValid == true){
                        toast.error("Passwords do not match!", {
                            position: toast.POSITION.TOP_RIGHT,
                        });
                    }else if(this.passwordValid == false){
                        toast.error("Password must be at least 6 characters long and include at least 1 number and 1 special character!", {
                            position: toast.POSITION.TOP_RIGHT,
                        });
                    }
                },
            togglePasswordVisibility() {
                this.showPassword = !this.showPassword;
            }
        }
    }
</script>
<style>
.reset-password-form {
    background: white;
    border-radius: 5px;
    box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
    transform: translateY(-50%);
    position: relative;
    top: 10%;
}
.loading-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  display: flex;
  justify-content: center;
  align-items: center;
  background: rgba(255, 255, 255, 0.7);
  z-index: 1000;
}

.gradient-background {
    background: rgb(107,139,246);
    background: linear-gradient(146deg, rgba(107,139,246,1) 0%, rgba(251,251,251,1) 100%);
}

.field-icon{
    float: right;
    margin-left: -25px;
    margin-top: -30px;
    margin-right: 20px;
    position: relative;
    z-index: 2;
}
</style>