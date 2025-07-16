import axios from "axios";

const $manager_api = axios.create({
    withCredentials: false,
    baseURL: process.env.NEXT_PUBLIC_MANAGER_API,
})

export default $manager_api;