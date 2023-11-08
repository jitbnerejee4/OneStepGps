import axios from "axios";
import { toast } from 'vue3-toastify';
import 'vue3-toastify/dist/index.css';
import store from "@/main";

const API_URL = 'http://localhost:8000';


class AuthService {
    async register(user) {
        try{
            const response = await axios.post(`${API_URL}/register`, user);
            localStorage.setItem('token', response.data.token);
            localStorage.setItem('email', response.data.email)
            return response
        }catch(err){
            toast.error(err.response.data.error+"! Please login", {
                position: toast.POSITION.TOP_LEFT,
            });
            return { error: "Registration failed", details: err };
        }
    }
    async login(user){
        try{
            const response = await axios.post(`${API_URL}/login`, user)
            localStorage.setItem('token', response.data.token)
            localStorage.setItem('email', response.data.email)
            return response
        }catch(err){
            return { error: "Registration failed", details: err };
        }
    }
    async savePreference(preference){
        try{
            const user = localStorage.getItem("email")
            const response = await axios.post(`${API_URL}/preferences/${user}`, preference)
            toast.success(response.data.message + " !", {
                position: toast.POSITION.TOP_RIGHT,
            });
            return response
        }catch(err){
            toast.error(err.response.data.error, {
                position: toast.POSITION.TOP_RIGHT,
            });
            return { error: "Failed to save preference!", details: err };
        }

    }
    async uploadImage(formData) {
        // Assuming you use axios for HTTP requests
        try{
            const user = localStorage.getItem("email")
            const response = await axios.post(`${API_URL}/upload/${user}`, formData, {
                headers: {
                    'Content-Type': 'multipart/form-data',
                    'Authorization': `Bearer ${localStorage.getItem('token')}`
                }
            });
            return response
        }catch(err){
            return { error: "Failed!", details: err };
        }
    }
    async forgotPassword(email){
        try{
            const response = await axios.get(`${API_URL}/forgot-password/${email}`)
            return response
        }catch(err){
            return {error: "Failed to save email!", details: err}
        }
    }
    async resetPassword(token, password){
        try{
            const data = {
                token: token,
                password: password
            };
            const response = await axios.post(`${API_URL}/reset-password`, data)            
            return response
        }catch(err){
            return {error: "Password update failed!", details: err}
        }
    }
    logout() {
        localStorage.removeItem('token');
        localStorage.removeItem('email')
        store.dispatch('logout'); 
    }
}

export default new AuthService();