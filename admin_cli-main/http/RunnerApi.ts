import axios from "axios";

const $runner_api = axios.create({
    withCredentials: false,
    baseURL: process.env.NEXT_PUBLIC_RUNNER_API,
})

export default $runner_api;