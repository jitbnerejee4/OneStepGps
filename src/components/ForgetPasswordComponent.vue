<template>
  <div class="main-container">
    <div class="form-container" style="height: 100vh;">
      <div class="row justify-content-center">
        <div class="col-md-6">
          <div class="card mt-5">
            <div class="card-header text-center">
              <h3>Forgot Password</h3>
            </div>
            <div class="card-body">
              <form @submit.prevent="submitForm">
                <div class="mb-3" v-if="!emailSent">
                  <label for="email" class="form-label">Email address</label>
                  <input type="email" class="form-control" id="email" v-model="email" required aria-describedby="emailHelp"/>
                  <div id="emailHelp" class="form-text">
                    Enter the email associated with your account.
                  </div>
                </div>
                <div class="mb-3" v-else>
                    <p> {{ message }}</p>
                </div>
                <button type="submit" class="btn btn-dark" v-if="!emailSent">Submit</button>
              </form>
            </div>
            <div class="card-footer text-center">
              Remember your password? <router-link to="/">Log in</router-link>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
  </template>
  
  <script>
  import AuthService from '@/services/AuthService';
  export default {
    name: "ForgotPassword",
    data() {
      return {
        email: "",
        emailSent: false,
        message: "Password reset link has been sent to your email.",
        success: true
      };
    },
    methods: {
      async submitForm() {
        const response = await AuthService.forgotPassword(this.email)
        if(response.status != 200){
            this.message = "Given email address is not found!"
            this.success = false
        }else{
            this.message = "Password reset link has been sent to your email."
        }
        this.email = ''
        this.emailSent = true

      },
    },
  };
  </script>
  
<style scoped>
.main-container {
  background: rgb(107,139,246);
  background: linear-gradient(146deg, rgba(107,139,246,1) 0%, rgba(251,251,251,1) 100%);
  overflow-x: hidden;
  overflow-y: hidden;
}
.form-container{
  position: relative;
  top: 275px;
}
</style>
  