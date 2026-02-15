import axios from "axios";
import { serverEnv } from "../env/server";

export const axiosInstance = axios.create({
  baseURL: serverEnv.BACKEND_ENDPOINT,
});
