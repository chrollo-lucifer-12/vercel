import "server-only";
import axios from "axios";
import { serverEnv } from "./env/server";

export const axiosInstance = axios.create({
  baseURL: serverEnv.BACKEND_ENDPOINT,
});

axiosInstance.interceptors.request.use((config) => {
  const token = localStorage.getItem("accessToken");
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});
