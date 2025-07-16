import axios from "axios";

const $auto_matcher_api = axios.create({
    withCredentials: false,
    baseURL: process.env.NEXT_PUBLIC_AUTO_MATCHER_API,
})

export default $auto_matcher_api;