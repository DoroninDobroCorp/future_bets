import axios from "axios";

const $calculator_api = axios.create({
    withCredentials: false,
    baseURL: process.env.NEXT_PUBLIC_CALCULATOR_API,
})

export default $calculator_api;